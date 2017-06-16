package main

import (
	"encoding/json"
	"strings"
	"os"
	"net"
	"io/ioutil"
	"net/http"
	"context"
	"log"
	"flag"
	"regexp"
	"time"
	"strconv"
	"os/exec"
)

var network string
var timeout int
var failure string
var segment string
var outputFile string
var isManualSourceEnv bool = false
var file *os.File

//syscall.Exec(os.Getenv("SHELL"), []string{os.Getenv("SHELL")}, syscall.Environ())

func main() {

	log.SetPrefix("DCE-APP-ENTRY-POINT ")

	setFlag()
	getEnv()
	showParam()

	if outputFile != "" {
		isManualSourceEnv = true
		f, err := os.Create(outputFile)
		if err != nil {
			fatalLog("create file %s error", outputFile)
		}
		file = f
	}

	switch network {
	case "port":
		setEnvInPortMapping()
	case "mac":
		setEnvInMacVlan()
	default:
		fatalLog("network must is port | mac")
	}

	runCommand()
}

func setAndWriteEnv(key string, val string) {
	if isManualSourceEnv {
		file.WriteString("export " + key + "=" + val + "\n")
	} else {
		os.Setenv(key, val)
	}
}

func runCommand() {

	if isManualSourceEnv {
		return
	}

	argsWithProg := flag.Args()
	command := argsWithProg[0]
	args := argsWithProg[1:]
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	log.Printf("command [ %s ], args %s , \n    Environ %s", command, args, os.Environ())

	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(0)
}

func setFlag() {
	flag.StringVar(&network, "network", "port", "network model:  port | mac , if use env key is DAE_NETWORK ")
	flag.IntVar(&timeout, "timeout", 20, "if in MACVLAN network is timeout, if use env key is DAE_TIMEOUT")
	flag.StringVar(&failure, "failure", "exit", "if set env failure, exit | continue , if use env key is DAE_FAILURE")
	flag.StringVar(&segment, "segment", "", "MACVLAN network segment regexp pattern, if use env key is DAE_SEGMENT")
	flag.StringVar(&outputFile, "output", "", "output file, if set this value, please source output.file, if use env key is DAE_OUTPUT ")
	flag.Parse()
}

func getEnv() {
	if os.Getenv("DAE_NETWORK") != "" {
		network = os.Getenv("DAE_NETWORK")
	}
	if os.Getenv("DAE_TIMEOUT") != "" {
		timeout, _ = strconv.Atoi(os.Getenv("DAE_TIMEOUT"))
	}
	if os.Getenv("DAE_FAILURE") != "" {
		network = os.Getenv("DAE_FAILURE")
	}
	if os.Getenv("DAE_SEGMENT") != "" {
		segment = os.Getenv("DAE_SEGMENT")
	}
	if os.Getenv("DAE_OUTPUT") != "" {
		outputFile = os.Getenv("DAE_OUTPUT")
	}

	//Default value
	if segment == "" {
		segment = "^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$";
	}
}

func showParam() {
	log.Printf("network: [ %s ], timeout: [ %d ], failure: [ %s ], segment: [ %s ]", network, timeout, failure, segment)
}

func setEnvInMacVlan() {
	log.Println("try set env in MACVLAN network")

	timeoutAt := time.Now().Add(time.Second * time.Duration(timeout))

	var matched bool = false

LOOP:
	for time.Now().Before(timeoutAt) && !matched {
		ifaces, err := net.Interfaces()
		if err != nil {
			fatalLog("can't get net Interfaces")
		}
		for _, i := range ifaces {
			addrs, err := i.Addrs()
			if err != nil {
				fatalLog("can't get interface ip")
			}

			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}

				if !isPublicIP(ip) {
					continue
				}

				ipString := ip.String()
				log.Printf("find ip [ %s ]", ipString)

				matched, err = regexp.MatchString(segment, ipString)
				if err != nil {
					fatalLog("can't MatchString %s with %s", ipString, segment)
				}

				if matched {
					setAndWriteEnv("DCE_ADVERTISE_IP", ipString)
					log.Printf("set DCE_ADVERTISE_IP to [ %s ]", ipString)
					break LOOP
				}

			}
		}

		time.Sleep(time.Second * time.Duration(5))
	}

	if !matched {
		fatalLog("timeout can't get macvlan ip...")
	}
}

func setEnvInPortMapping() {
	log.Println("try set env in Portmapping network")

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/dce-metadata/dce-metadata.sock")
			},
		},
	}

	hostname := os.Getenv("HOSTNAME")
	log.Printf("HOSTNAME is %s", hostname)

	hostInfoUrl := "http://unix/containers/" + hostname + "/json"
	log.Printf("hostInfoUrl is %s", hostInfoUrl)

	resp, err := client.Get(hostInfoUrl)
	if err != nil {
		fatalLog("can't get host info from " + hostInfoUrl)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fatalLog("read data from resp err")
	}

	var portInfo = &PortInfo{}

	err = json.Unmarshal(data, portInfo)
	if err != nil {
		fatalLog("Unmarshal json data error: " + string(data))
	}
	var isOnly bool = true
	for key, vale := range portInfo.NetworkSettings.Ports {
		keys := strings.Split(key, "/")
		innerPort := keys[0]
		innerProtocol := keys[1]
		hostPort := vale[0].HostPort

		log.Printf("innerPort [%s], innerProtocol [%s], hostPort [%s]", innerPort, innerProtocol, hostPort)

		if isOnly {
			setAndWriteEnv("DCE_ADVERTISE_PORT", hostPort)
			isOnly = false
		}

		setAndWriteEnv("DCE_ADVERTISE_PORT_"+innerPort, hostPort)
	}

	hostInfoUrl = "http://unix/info"
	resp, err = client.Get(hostInfoUrl)
	if err != nil {
		fatalLog("can't get host info from " + hostInfoUrl)
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fatalLog("read data from resp err")
	}

	var ipInfo = &IpInfo{}
	err = json.Unmarshal(data, ipInfo)
	if err != nil {
		fatalLog("Unmarshal json data error: " + string(data))
	}

	log.Printf("ip address [%s]", ipInfo.Swarm.NodeAddr)
	setAndWriteEnv("DCE_ADVERTISE_IP", ipInfo.Swarm.NodeAddr)
}

func isPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

func fatalLog(format string, v ...interface{}) {
	log.Printf(format, v...)
	switch failure {
	case "continue":
		runCommand()
	default:
		os.Exit(1)
	}
}

//Json def

type PortInfo struct {
	NetworkSettings NetworkSettings `json:"NetworkSettings"`
}

type NetworkSettings struct {
	Ports Ports `json:"Ports"`
}

type Ports map[string][]Port

type Port struct {
	HostIp   string `json:"HostIp"`
	HostPort string `json:"HostPort"`
}

// Host IP INFO
type IpInfo struct {
	Swarm Swarm `json:"Swarm"`
}

type Swarm struct {
	NodeAddr string `json:"NodeAddr"`
}

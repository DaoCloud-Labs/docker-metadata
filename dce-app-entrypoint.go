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

//syscall.Exec(os.Getenv("SHELL"), []string{os.Getenv("SHELL")}, syscall.Environ())

func main() {

	log.SetPrefix("DCE-APP-ENTRY-POINT ")

	setFlag()
	getEnv()
	showParam()

	switch network {
	case "port":
		setEnvInPortMapping()
	case "mac":
		setEnvInMacVlan()
	default:
		fatalLog("network must is port | mac")
	}

	argsWithProg := flag.Args()
	command := argsWithProg[0]
	args := argsWithProg[1:]
	log.Printf("command [ %s ], args %s", command, args)

	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func setFlag() {
	flag.StringVar(&network, "network", "port", "network model:  port | mac , if use env key is DAE_NETWORK ")
	flag.IntVar(&timeout, "timeout", 20, "if in MACVLAN network is timeout, if use env key is DAE_TIMEOUT")
	flag.StringVar(&failure, "failure", "exit", "if set env failure, exit | continue , if use env key is DAE_FAILURE")
	flag.StringVar(&segment, "segment", "", "MACVLAN network segment regexp pattern, if use env key is DAE_SEGMENT")
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
}

func showParam() {
	log.Printf("network: [ %s ], timeout: [ %d ], failure: [ %s ], segment: [ %s ]", network, timeout, failure, segment)
}

func setEnvInMacVlan() {
	log.Println("try set env in MACVLAN network")

	timeoutAt := time.Now().Add(time.Second * time.Duration(timeout))

	var matched bool = false

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

				ipString := ip.String()
				log.Printf("find ip [ %s ]", ipString)

				matched, err = regexp.MatchString(segment, ipString)

				if err != nil {
					fatalLog("can't MatchString %s with %s", ipString, segment)
				}

				if matched {
					os.Setenv("DCE_ADVERTISE_IP", ipString)
					log.Printf("set DCE_ADVERTISE_IP to [ %s ]", ipString)
					break
				}

			}
		}
		if matched {
			break
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
	log.Printf("HOSTNAME is %s \n", hostname)

	hostInfoUrl := "http://unix/containers/" + hostname + "/json"
	log.Printf("hostInfoUrl is %s \n", hostInfoUrl)

	resp, err := client.Get(hostInfoUrl)
	if err != nil {
		fatalLog("can't get host info from " + hostInfoUrl)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fatalLog("read data from resp err")
	}

	var portInfo *PortInfo = &PortInfo{}

	err = json.Unmarshal(data, portInfo)
	if err != nil {
		fatalLog("Unmarshal json data error: " + string(data[:]))
	}
	var isOnly bool = true
	for key, vale := range portInfo.NetworkSettings.Ports {
		keys := strings.Split(key, "/")
		innerPort := keys[0]
		innerProtocol := keys[1]
		hostPort := vale[0].HostPort

		log.Printf("innerPort [%s], innerProtocol [%s], hostPort [%s] \n", innerPort, innerProtocol, hostPort)

		if isOnly {
			os.Setenv("DCE_ADVERTISE_PORT", hostPort)
			isOnly = false
		}

		os.Setenv("DCE_ADVERTISE_PORT_"+innerPort, hostPort)
	}

	hostInfoUrl = "http://unix/info"
	resp, err = client.Get(hostInfoUrl)
	if err != nil {
		fatalLog("can't get host info from " + hostInfoUrl)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fatalLog("read data from resp err")
	}

	var ipInfo *IpInfo = &IpInfo{}
	err = json.Unmarshal(data, ipInfo)
	if err != nil {
		fatalLog("Unmarshal json data error: " + string(data[:]))
	}

	log.Printf("ip address [%s] \n", ipInfo.Swarm.NodeAddr)
	os.Setenv("DCE_ADVERTISE_IP", ipInfo.Swarm.NodeAddr)
}

func fatalLog(v ...interface{}) {
	log.Println(v)
	switch failure {
	case "continue":
		os.Exit(0)
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
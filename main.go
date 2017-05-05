package main

import (
	"fmt"
	"encoding/json"
	"strings"
	"os"
	"net"
	"io/ioutil"
	"net/http"
	"context"
)

func main() {

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/dce-metadata/dce-metadata.sock")
			},
		},
	}

	hostname := os.Getenv("HOSTNAME")
	resp, err := client.Get("http://unix/containers/" + hostname + "/json")

	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	//data := []byte(`{"Id":"c9f053ffc7359f05ced33b8a72143557cdaf27ae9cfce1744cfdc0d4e565653a","Created":"2017-05-04T16:22:43.641224059Z","Path":"/bin/bash","Args":[],"State":{"Status":"running","Running":true,"Paused":false,"Restarting":false,"OOMKilled":false,"Dead":false,"Pid":15303,"ExitCode":0,"Error":"","StartedAt":"2017-05-04T16:22:43.903611376Z","FinishedAt":"0001-01-01T00:00:00Z"},"Image":"sha256:a8493f5f50ffda70c2eeb2d09090debf7d39c8ffcd63b43ff81b111ece6f28bf","ResolvConfPath":"/var/lib/docker/containers/c9f053ffc7359f05ced33b8a72143557cdaf27ae9cfce1744cfdc0d4e565653a/resolv.conf","HostnamePath":"/var/lib/docker/containers/c9f053ffc7359f05ced33b8a72143557cdaf27ae9cfce1744cfdc0d4e565653a/hostname","HostsPath":"/var/lib/docker/containers/c9f053ffc7359f05ced33b8a72143557cdaf27ae9cfce1744cfdc0d4e565653a/hosts","LogPath":"/var/lib/docker/containers/c9f053ffc7359f05ced33b8a72143557cdaf27ae9cfce1744cfdc0d4e565653a/c9f053ffc7359f05ced33b8a72143557cdaf27ae9cfce1744cfdc0d4e565653a-json.log","Name":"/agitated_lamarr","RestartCount":0,"Driver":"aufs","MountLabel":"","ProcessLabel":"","AppArmorProfile":"docker-default","ExecIDs":null,"HostConfig":{"Binds":["/var/run/dce-metadata:/var/run/dce-metadata"],"ContainerIDFile":"","LogConfig":{"Type":"json-file","Config":{}},"NetworkMode":"default","PortBindings":{"9090/tcp":[{"HostIp":"","HostPort":"899"}]},"RestartPolicy":{"Name":"no","MaximumRetryCount":0},"AutoRemove":false,"VolumeDriver":"","VolumesFrom":null,"CapAdd":null,"CapDrop":null,"Dns":[],"DnsOptions":[],"DnsSearch":[],"ExtraHosts":null,"GroupAdd":null,"IpcMode":"","Cgroup":"","Links":null,"OomScoreAdj":0,"PidMode":"","Privileged":false,"PublishAllPorts":false,"ReadonlyRootfs":false,"SecurityOpt":null,"UTSMode":"","UsernsMode":"","ShmSize":67108864,"Runtime":"runc","ConsoleSize":[0,0],"Isolation":"","CpuShares":0,"Memory":0,"NanoCpus":0,"CgroupParent":"","BlkioWeight":0,"BlkioWeightDevice":null,"BlkioDeviceReadBps":null,"BlkioDeviceWriteBps":null,"BlkioDeviceReadIOps":null,"BlkioDeviceWriteIOps":null,"CpuPeriod":0,"CpuQuota":0,"CpuRealtimePeriod":0,"CpuRealtimeRuntime":0,"CpusetCpus":"","CpusetMems":"","Devices":[],"DeviceCgroupRules":null,"DiskQuota":0,"KernelMemory":0,"MemoryReservation":0,"MemorySwap":0,"MemorySwappiness":-1,"OomKillDisable":false,"PidsLimit":0,"Ulimits":null,"CpuCount":0,"CpuPercent":0,"IOMaximumIOps":0,"IOMaximumBandwidth":0},"GraphDriver":{"Data":null,"Name":"aufs"},"Mounts":[{"Type":"bind","Source":"/var/run/dce-metadata","Destination":"/var/run/dce-metadata","Mode":"","RW":true,"Propagation":""}],"Config":{"Hostname":"c9f053ffc735","Domainname":"","User":"","AttachStdin":true,"AttachStdout":true,"AttachStderr":true,"ExposedPorts":{"9090/tcp":{}},"Tty":true,"OpenStdin":true,"StdinOnce":true,"Env":["PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"],"Cmd":["/bin/bash"],"Image":"centos","Volumes":null,"WorkingDir":"","Entrypoint":null,"OnBuild":null,"Labels":{"build-date":"20170406","license":"GPLv2","name":"CentOS Base Image","vendor":"CentOS"}},"NetworkSettings":{"Bridge":"","SandboxID":"1679aaa65bebbf123d739724dc1b5a99e14b9183d7c06e70c3a857e770e6b11b","HairpinMode":false,"LinkLocalIPv6Address":"","LinkLocalIPv6PrefixLen":0,"Ports":{"9090/tcp":[{"HostIp":"0.0.0.0","HostPort":"899"}]},"SandboxKey":"/var/run/docker/netns/1679aaa65beb","SecondaryIPAddresses":null,"SecondaryIPv6Addresses":null,"EndpointID":"f5686ffb986d289c6ce925d77a38e9119e7bc37f90106a865b08fd91f7661019","Gateway":"172.17.0.1","GlobalIPv6Address":"","GlobalIPv6PrefixLen":0,"IPAddress":"172.17.0.6","IPPrefixLen":16,"IPv6Gateway":"","MacAddress":"02:42:ac:11:00:06","Networks":{"bridge":{"IPAMConfig":null,"Links":null,"Aliases":null,"NetworkID":"d873a3f962ee8f8e3e990c779d5b745aee6a012399d148e6433ee2f816329299","EndpointID":"f5686ffb986d289c6ce925d77a38e9119e7bc37f90106a865b08fd91f7661019","Gateway":"172.17.0.1","IPAddress":"172.17.0.6","IPPrefixLen":16,"IPv6Gateway":"","GlobalIPv6Address":"","GlobalIPv6PrefixLen":0,"MacAddress":"02:42:ac:11:00:06"}}}}`)

	var portInfo *PortInfo = &PortInfo{}

	err = json.Unmarshal(data, portInfo)
	if err != nil {
		panic(err)
	}

	w, _ := os.Create("set_env.sh")

	for key, vale := range portInfo.NetworkSettings.Ports {
		keys := strings.Split(key, "/")
		innerPort := keys[0]
		innerProtocol := keys[1]
		hostPort := vale[0].HostPort

		fmt.Printf("innerPort [%s], innerProtocol [%s], hostPort [%s] \n", innerPort, innerProtocol, hostPort)

		w.WriteString("export D_HOST_PORT_" + strings.ToUpper(innerProtocol) + "_" + innerPort + "=" + hostPort + "\n")
	}

	resp, err = client.Get("http://unix/info")
	if err != nil {
		panic(err)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	//data = []byte(`{"ID":"SDAF:PAFN:OBD7:T5MB:IPTI:VRZZ:MKXP:O33Z:R3VV:MAR3:N3VK:FFQ2","Containers":120,"ContainersRunning":10,"ContainersPaused":0,"ContainersStopped":110,"Images":59,"Driver":"aufs","DriverStatus":[["Root Dir","/var/lib/docker/aufs"],["Backing Filesystem","extfs"],["Dirs","645"],["Dirperm1 Supported","false"]],"SystemStatus":null,"Plugins":{"Volume":["local"],"Network":["bridge","host","macvlan","null","overlay"],"Authorization":[]},"MemoryLimit":true,"SwapLimit":false,"KernelMemory":true,"CpuCfsPeriod":true,"CpuCfsQuota":true,"CPUShares":true,"CPUSet":true,"IPv4Forwarding":true,"BridgeNfIptables":true,"BridgeNfIp6tables":true,"Debug":false,"NFd":378,"OomKillDisable":true,"NGoroutines":618,"SystemTime":"2017-05-05T02:48:06.318720524Z","LoggingDriver":"json-file","CgroupDriver":"cgroupfs","NEventsListener":9,"KernelVersion":"3.13.0-48-generic","OperatingSystem":"Ubuntu 14.04.2 LTS","OSType":"linux","Architecture":"x86_64","IndexServerAddress":"https://index.docker.io/v1/","RegistryConfig":{"InsecureRegistryCIDRs":["0.0.0.0/0","127.0.0.0/8"],"IndexConfigs":{"docker.io":{"Name":"docker.io","Mirrors":["http://537c24fb.m.daocloud.io/"],"Secure":true,"Official":true}},"Mirrors":["http://537c24fb.m.daocloud.io/"]},"NCPU":2,"MemTotal":4143964160,"DockerRootDir":"/var/lib/docker","HttpProxy":"","HttpsProxy":"","NoProxy":"","Name":"ip-172-31-60-44","Labels":null,"ExperimentalBuild":false,"ServerVersion":"17.04.0-ce","ClusterStore":"","ClusterAdvertise":"","Runtimes":{"runc":{"path":"docker-runc"}},"DefaultRuntime":"runc","Swarm":{"NodeID":"swkxw9lfo8ln78ct9krihw8xp","NodeAddr":"172.31.60.44","LocalNodeState":"active","ControlAvailable":false,"Error":"","RemoteManagers":[{"NodeID":"rt7noge0po1ld9qoqll873tm3","Addr":"172.31.60.236:2377"}],"Nodes":0,"Managers":0,"Cluster":{"ID":"","Version":{},"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","Spec":{"Labels":null,"Orchestration":{},"Raft":{"ElectionTick":0,"HeartbeatTick":0},"Dispatcher":{},"CAConfig":{},"TaskDefaults":{},"EncryptionConfig":{"AutoLockManagers":false}}}},"LiveRestoreEnabled":false,"Isolation":"","InitBinary":"","ContainerdCommit":{"ID":"422e31ce907fd9c3833a38d7b8fdd023e5a76e73","Expected":"422e31ce907fd9c3833a38d7b8fdd023e5a76e73"},"RuncCommit":{"ID":"9c2d8d184e5da67c95d601382adf14862e4f2228","Expected":"9c2d8d184e5da67c95d601382adf14862e4f2228"},"InitCommit":{"ID":"949e6fa","Expected":"949e6fa"},"SecurityOptions":["name=apparmor"]}`)

	var ipInfo *IpInfo = &IpInfo{}
	err = json.Unmarshal(data, ipInfo)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ip address [%s] \n", ipInfo.Swarm.NodeAddr)
	w.WriteString("export D_HOST_IP=" + ipInfo.Swarm.NodeAddr)

	w.Close()
}

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

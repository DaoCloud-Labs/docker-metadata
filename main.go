package main

import (
	"fmt"
	//"net"
	//"io/ioutil"
	"encoding/json"
)

func main() {
	fmt.Println("loading info from ")

	//c, err := net.Dial("/var/run/dce-metadata/dce-metadata.sock", "http:/containers/$(hostname)/json")
	//
	//if err != nil {
	//	panic(err)
	//}

	//data, err := ioutil.ReadAll(c)
	//if err != nil {
	//	panic(err)
	//}

	data := []byte(`{"Id":"c9f053ffc7359f05ced33b8a72143557cdaf27ae9cfce1744cfdc0d4e565653a","Created":"2017-05-04T16:22:43.641224059Z","Path":"/bin/bash","Args":[],"State":{"Status":"running","Running":true,"Paused":false,"Restarting":false,"OOMKilled":false,"Dead":false,"Pid":15303,"ExitCode":0,"Error":"","StartedAt":"2017-05-04T16:22:43.903611376Z","FinishedAt":"0001-01-01T00:00:00Z"},"Image":"sha256:a8493f5f50ffda70c2eeb2d09090debf7d39c8ffcd63b43ff81b111ece6f28bf","ResolvConfPath":"/var/lib/docker/containers/c9f053ffc7359f05ced33b8a72143557cdaf27ae9cfce1744cfdc0d4e565653a/resolv.conf","HostnamePath":"/var/lib/docker/containers/c9f053ffc7359f05ced33b8a72143557cdaf27ae9cfce1744cfdc0d4e565653a/hostname","HostsPath":"/var/lib/docker/containers/c9f053ffc7359f05ced33b8a72143557cdaf27ae9cfce1744cfdc0d4e565653a/hosts","LogPath":"/var/lib/docker/containers/c9f053ffc7359f05ced33b8a72143557cdaf27ae9cfce1744cfdc0d4e565653a/c9f053ffc7359f05ced33b8a72143557cdaf27ae9cfce1744cfdc0d4e565653a-json.log","Name":"/agitated_lamarr","RestartCount":0,"Driver":"aufs","MountLabel":"","ProcessLabel":"","AppArmorProfile":"docker-default","ExecIDs":null,"HostConfig":{"Binds":["/var/run/dce-metadata:/var/run/dce-metadata"],"ContainerIDFile":"","LogConfig":{"Type":"json-file","Config":{}},"NetworkMode":"default","PortBindings":{"9090/tcp":[{"HostIp":"","HostPort":"899"}]},"RestartPolicy":{"Name":"no","MaximumRetryCount":0},"AutoRemove":false,"VolumeDriver":"","VolumesFrom":null,"CapAdd":null,"CapDrop":null,"Dns":[],"DnsOptions":[],"DnsSearch":[],"ExtraHosts":null,"GroupAdd":null,"IpcMode":"","Cgroup":"","Links":null,"OomScoreAdj":0,"PidMode":"","Privileged":false,"PublishAllPorts":false,"ReadonlyRootfs":false,"SecurityOpt":null,"UTSMode":"","UsernsMode":"","ShmSize":67108864,"Runtime":"runc","ConsoleSize":[0,0],"Isolation":"","CpuShares":0,"Memory":0,"NanoCpus":0,"CgroupParent":"","BlkioWeight":0,"BlkioWeightDevice":null,"BlkioDeviceReadBps":null,"BlkioDeviceWriteBps":null,"BlkioDeviceReadIOps":null,"BlkioDeviceWriteIOps":null,"CpuPeriod":0,"CpuQuota":0,"CpuRealtimePeriod":0,"CpuRealtimeRuntime":0,"CpusetCpus":"","CpusetMems":"","Devices":[],"DeviceCgroupRules":null,"DiskQuota":0,"KernelMemory":0,"MemoryReservation":0,"MemorySwap":0,"MemorySwappiness":-1,"OomKillDisable":false,"PidsLimit":0,"Ulimits":null,"CpuCount":0,"CpuPercent":0,"IOMaximumIOps":0,"IOMaximumBandwidth":0},"GraphDriver":{"Data":null,"Name":"aufs"},"Mounts":[{"Type":"bind","Source":"/var/run/dce-metadata","Destination":"/var/run/dce-metadata","Mode":"","RW":true,"Propagation":""}],"Config":{"Hostname":"c9f053ffc735","Domainname":"","User":"","AttachStdin":true,"AttachStdout":true,"AttachStderr":true,"ExposedPorts":{"9090/tcp":{}},"Tty":true,"OpenStdin":true,"StdinOnce":true,"Env":["PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"],"Cmd":["/bin/bash"],"Image":"centos","Volumes":null,"WorkingDir":"","Entrypoint":null,"OnBuild":null,"Labels":{"build-date":"20170406","license":"GPLv2","name":"CentOS Base Image","vendor":"CentOS"}},"NetworkSettings":{"Bridge":"","SandboxID":"1679aaa65bebbf123d739724dc1b5a99e14b9183d7c06e70c3a857e770e6b11b","HairpinMode":false,"LinkLocalIPv6Address":"","LinkLocalIPv6PrefixLen":0,"Ports":{"9090/tcp":[{"HostIp":"0.0.0.0","HostPort":"899"}]},"SandboxKey":"/var/run/docker/netns/1679aaa65beb","SecondaryIPAddresses":null,"SecondaryIPv6Addresses":null,"EndpointID":"f5686ffb986d289c6ce925d77a38e9119e7bc37f90106a865b08fd91f7661019","Gateway":"172.17.0.1","GlobalIPv6Address":"","GlobalIPv6PrefixLen":0,"IPAddress":"172.17.0.6","IPPrefixLen":16,"IPv6Gateway":"","MacAddress":"02:42:ac:11:00:06","Networks":{"bridge":{"IPAMConfig":null,"Links":null,"Aliases":null,"NetworkID":"d873a3f962ee8f8e3e990c779d5b745aee6a012399d148e6433ee2f816329299","EndpointID":"f5686ffb986d289c6ce925d77a38e9119e7bc37f90106a865b08fd91f7661019","Gateway":"172.17.0.1","IPAddress":"172.17.0.6","IPPrefixLen":16,"IPv6Gateway":"","GlobalIPv6Address":"","GlobalIPv6PrefixLen":0,"MacAddress":"02:42:ac:11:00:06"}}}}`)
	var portInfo *PortInfo = &PortInfo{}
	err := json.Unmarshal(data, portInfo)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(portInfo)
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
//
//type SingPort struct {
//	HostIp   string `json:"HostIp"`
//	HostPort string `json:"HostPort"`
//}

# Docker - metadata

从 dce-metadata-agent 服务中获取 端口映射信息和IP地址

使用Go请求获得Host的IP和端口

## PRE

创建全局 dce-metadata-agent 服务
```bash
    docker service create --name dce-metadata-agent --mode global   --mount type=bind,src=/var/run,dst=/var/run -l io.daocloud.dce.system=build-in daocloud.io/daocloud/dce-metadata-agent 
```

## Build

推荐使用 docker

```bash
    docker run -it  -v $PWD:/usr/code -w /usr/code  golang go build docker-metadata.go
```

或者使用 

```bash
    env GOOS=linux GOARCH=amd64 go build .
```

---

## Usage

```bash
    ./docker-metadata && sh set_env.sh
```

---

## Sample
假设我们使用以下命令创建一个容器 

```bash
    docker run -it -v /var/run/dce-metadata:/var/run/dce-metadata -v $PWD:/usr/code -w /usr/code  -p 8001:8888 -p 8002:8889 centos
```

我们可以在容器中执行

```bash
    ./docker-metadata && source set_env.sh
```

使用 env 我们可以看见如下信息
```bash
HOSTNAME=e6fb4eb9618a
...
D_HOST_IP=172.31.60.44
D_HOST_PORT_TCP_8889=8002
D_HOST_PORT_TCP_8888=8001
```

---

## Eureka Client Settings

add following properties in application properties file
```
eureka.instance.prefer-ip-address=true
```
## Dockerfile

Please refer to sample-Dockerfile

### PS
仓库中的是AMD64 LINUX 编译版本

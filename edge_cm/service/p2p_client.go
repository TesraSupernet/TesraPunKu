package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/DOSNetwork/core/edge_cm/conf"
	"github.com/DOSNetwork/core/edge_cm/tools"
	_ "github.com/alecthomas/template"
	docker "github.com/fsouza/go-dockerclient"
	"os"
)

//const dockerFileTemplate = `
//#基础镜像
//FROM ubuntu
//#作者
//MAINTAINER wangXin <929592332@qq.com>
//#定义工作目录
//ENV WORK_PATH /usr/local/p2p-client
//ENV P2P_VERSION %s
//EXPOSE 9097/udp
//RUN  apt-get update \
//  && apt-get install -y wget \
//  && rm -rf /var/lib/apt/lists/*
//RUN mkdir -p /usr/local/p2p-client/ \
//  && wget -q -O /usr/local/p2p-client/main "%s" \
//  && chmod u+x /usr/local/p2p-client/main
//ENTRYPOINT ["./usr/local/p2p-client/main","%s"]
//`
const dockerFileTemplate = `
#基础镜像
FROM ubuntu:p2p-client
#作者
MAINTAINER wangXin <929592332@qq.com>
#定义工作目录
ENV WORK_PATH /usr/local/p2p-client
ENV P2P_VERSION %s
EXPOSE 9097/udp
RUN mkdir -p /usr/local/p2p-client/ \
  && wget -q -O /usr/local/p2p-client/main "%s" \
  && chmod u+x /usr/local/p2p-client/main
ENTRYPOINT ["./usr/local/p2p-client/main","%s"]
`

type Mount struct {
	Type        string
	Source      string
	Destination string
	Mode        string
	RW          bool
	Propagation string
}

// start p2p client
func StartP2pClient() {
	if conf.Conf.P2pClient.Enabled {
		filter := []string{"serverName=p2pClient"}
		containers, err := Containers(map[string][]string{"label": filter})
		if err != nil {
			tools.Logger.Error(fmt.Sprintf("[start-p2p-client]:get container list error: %s", err.Error()))
			return
		}
		isStartContainer := true
		for _, container := range containers {
			labels := container.Labels
			version := labels["version"]
			if conf.Conf.P2pClient.DefaultVersion == version {
				isStartContainer = false
				continue
			}
			err := Client.StopContainer(container.ID, conf.Conf.Interval.StopContainerTimeout)
			if err != nil {
				tools.Logger.Error(fmt.Sprintf("[start-p2p-client]:stop container(%s) error: %s", container.Names, err.Error()))
				// TODO send mq message
			}
		}
		if isStartContainer {
			err := startContainer()
			if err != nil {
				tools.Logger.Error(fmt.Sprintf("[start-p2p-client]:start p2p-client container error: %s", err.Error()))
				// TODO send mq message
			}
		}
		//
	}
}

// build p2pClient image
func buildImage() (imageName string, err error) {
	version := conf.Conf.P2pClient.DefaultVersion
	//tempDownloadPath := ztools.GetTempDownloadPath(bucketName, fmt.Sprintf("%s/main", version))
	tempDownloadPath := fmt.Sprintf("%s/p2p-client/%s/main", conf.Conf.Minio.Endpoint, version)
	imageName = fmt.Sprintf("p2p-client:%s", version)
	mac := tools.GetMacByName(conf.Conf.MMLServer.MAC_ADDR)
	tools.Logger.Info(fmt.Sprintf("[p2p-client images building]:tempDownloadPath is %s", tempDownloadPath))
	dockerFileStr := fmt.Sprintf(dockerFileTemplate, version, tempDownloadPath, mac)
	dockerfile, err := BuildDockerfile(dockerFileStr)
	if err != nil {
		return
	}
	opts := docker.BuildImageOptions{
		Context:             context.Background(),
		Name:                imageName,
		InputStream:         dockerfile,
		OutputStream:        bytes.NewBuffer(nil),
		RmTmpContainer:      true,
		ForceRmTmpContainer: true,
		Labels:              map[string]string{"serverName": "p2pClient", "version": version}, // set meta data
	}
	err = Client.BuildImage(opts)
	return
}

func startContainer() (err error) {
	version := conf.Conf.P2pClient.DefaultVersion
	containerName := fmt.Sprintf("p2p-client-container-%s", version)
	imageName, err := buildImage()
	if err != nil {
		tools.Logger.Error(fmt.Sprintf("[start-p2p-client]:build p2p-client image error: %s", err.Error()))
		// TODO send mq message
		return
	}
	config := &docker.Config{
		//Cmd:   append(strings.Fields(job.DsmCommand), string(buf)), // add start command
		Image:  imageName,                                                        // image id
		Labels: map[string]string{"serverName": "p2pClient", "version": version}, // set meta data
		//Mounts: []docker.Mount{
		//	{
		//		RW:          true,
		//		Source:      "/home/data/",
		//		Destination: "/mnt/",
		//	},
		//},
	}
	hostConfig := docker.HostConfig{
		//AutoRemove: true,
		// 挂载卷在Mounts上面配置无效，好像是个无效配置；挂载请用以下配置，详情请看：https://github.com/fsouza/go-dockerclient/issues/633
		Binds:        []string{"/home/data/:/home/data/"},
		PortBindings: map[docker.Port][]docker.PortBinding{"9097/udp": {{HostIP: "0.0.0.0", HostPort: "9097"}}},
	}
	createContainerOptions := docker.CreateContainerOptions{
		Name:             containerName,
		Config:           config,
		HostConfig:       &hostConfig,
		NetworkingConfig: nil,
		Context:          context.Background(),
	}
	container, err := Client.CreateContainer(createContainerOptions)
	if err != nil {
		tools.Logger.Error(fmt.Sprintf("[start-p2p-client]:create p2p-client container error: %s", err.Error()))
		// TODO send mq message
		return
	}
	err = Client.StartContainer(container.ID, nil)
	return
}

func DeleteNodeData(jobID string) {
	fileName := fmt.Sprintf("/home/data/%s", jobID)
	err := os.Remove(fileName)
	if err != nil {
		tools.Logger.Error(fmt.Sprintf("[Delete Node Data]:file %s delete failed!", fileName))
	}
}

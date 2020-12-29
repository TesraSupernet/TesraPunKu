package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	get         = kingpin.Command("get", "get dockerfile.")
	getFlag     = get.Flag("taskid", "get the dockerfile corresponding to taskid.").Required().String()
	pull        = kingpin.Command("pull", "pull image.")
	pullFlag    = pull.Flag("taskid", "pull image by taskid.").Required().String()
	del         = kingpin.Command("del", "delete image.")
	delFlag     = del.Flag("taskid", "delete image by taskid.").Required().String()
	run         = kingpin.Command("run", "run container.")
	runFlag     = run.Flag("taskid", "run container by taskid.").Required().String()
	stop        = kingpin.Command("stop", "stop container.")
	stopFlag    = stop.Flag("taskid", "stop container by taskid.").Required().String()
	restart     = kingpin.Command("restart", "restart container.")
	restartFlag = restart.Flag("taskid", "restart container by taskid.").Required().String()
)

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("deng hong")
	kingpin.CommandLine.Help = "container management subsystem command line assistants."
	switch kingpin.Parse() {
	case "get":
		Request(*getFlag, GET_DOCKERFILE_URL)
	case "pull":
		Request(*pullFlag, PULL_IMAGE_URL)
	case "del":
		Request(*delFlag, DEL_IMAGE_URL)
	case "run":
		Request(*runFlag, RUN_CONTAINER_URL)
	case "stop":
		Request(*stopFlag, STOP_CONTAINER_URL)
	case "restart":
		Request(*restartFlag, RESTART_CONTAINER_URL)
	}
}

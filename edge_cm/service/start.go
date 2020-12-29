package service

func Run() {
	SystemCheck()

	go PushJobStatusChanged()
	go PushHardwareInfo()
	go PushRunningJob(true)

	go PullJobImageBuildRun()
	go PullJobStatus()
	go PullJobCommand()

	go CheckFailedMessageAndResend()
	go CheckAndAutoRmImage()
	go ClearDanglingImage()
	go CheckContainerStatus()
	go CheckPendingJob()
	go StartP2pClient()
	//Serve()

}

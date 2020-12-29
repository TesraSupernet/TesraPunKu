package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DOSNetwork/core/edge_cm/conf"
	"github.com/DOSNetwork/core/edge_cm/tools"
	docker "github.com/fsouza/go-dockerclient"
	"strings"
	"sync"
	"time"
)

var JobDispatchers sync.Map

// Start the container and generate the state based on the job information
func StartContainer(job *tools.EdgeCmJob, jobCh chan<- *tools.JobStatus) {
	status := &tools.JobStatus{
		TaskId:    job.TaskId,
		JobId:     job.JobId,
		Timestamp: int64(time.Now().Unix()),
	}
	buf, err := json.Marshal(job)
	if err != nil {
		msg := err.Error()
		status.Code = conf.StartContainerFailedCode
		status.Status = conf.StartContainerFailure
		status.Msg = msg
		// send start container failed status
		jobCh <- status
		tools.Logger.Error(msg)
	}
	config := &docker.Config{
		Tty:   true,                                                // Distribution terminal
		Cmd:   append(strings.Fields(job.DsmCommand), string(buf)), // add start command
		Image: job.JobId,                                           // image id
		//Labels: map[string]string{"job": "true", "taskid": job.TaskId, "jobid": job.JobId}, // set meta data
	}
	var hostConfig docker.HostConfig
	// auto remove container
	hostConfig.AutoRemove = conf.Conf.Container.AutoRm

	switch job.Type {
	// AI job
	case conf.AI, conf.VASP:
		// ai job need nvidia-docker e.g: docker run -ti --runtime=nvidia tesra:v1.0
		hostConfig.Runtime = "nvidia"
	case conf.P2P:
		hostConfig.Binds = []string{"/home/data/:/tmp/ec.job.p2p/"}
	}

	createContainerOptions := docker.CreateContainerOptions{
		Name:             job.JobId, // container name
		Config:           config,
		HostConfig:       &hostConfig,
		NetworkingConfig: nil,
		Context:          context.Background(),
	}
	// create container
	container, err := Client.CreateContainer(createContainerOptions)
	if err != nil {
		msg := err.Error()
		status.Msg = msg
		status.Code = conf.StartContainerFailedCode
		status.Status = conf.StartContainerFailure
		// send start container failed status
		jobCh <- status
		tools.Logger.Error(msg)
		return
	}
	// start container
	err = Client.StartContainer(container.ID, nil)
	if err != nil {
		msg := err.Error()
		status.Msg = msg
		status.Code = conf.StartContainerFailedCode
		status.Status = conf.StartContainerFailure
		// send start container failed status
		jobCh <- status
		tools.Logger.Error(msg)
		return
	}
	status.Msg = "start container success."
	status.Code = conf.StartContainerSuccessCode
	status.Status = conf.StartContainerSuccess
	// send start container success status
	jobCh <- status
	tools.Logger.Info("start container success.")
	// store cwa to sqlite
	dispatcher := StoreCWA(job.JobId)
	// save to JobDispatchers
	JobDispatchers.Store(job.JobId, dispatcher)
	// send running job
	//PushRunningJob(false)
}

// Stop the container and generate the state based on the job information
func StopContainer(taskId, jobId string, jobCh chan<- *tools.JobStatus) {
	status := &tools.JobStatus{
		TaskId:    taskId,
		JobId:     jobId,
		Timestamp: int64(time.Now().Unix()),
	}
	// stop container
	err := Client.StopContainer(jobId, conf.Conf.Interval.StopContainerTimeout)
	if err != nil {
		msg := err.Error()
		status.Code = conf.StopContainerFailedCode
		status.Status = conf.StopFailure
		status.Msg = msg
		// send stop container failed status
		jobCh <- status
		tools.Logger.Error(msg)
		return
	}
	status.Code = conf.StopContainerSuccessCode
	status.Status = conf.StopSuccess
	status.Msg = "container  stop success."
	// send stop container success status
	jobCh <- status
	tools.Logger.Info("container  stop success.")
}

// Restart the container and generate the state based on the job information
func RestartContainer(taskId, jobId string, jobCh chan<- *tools.JobStatus) {
	status := &tools.JobStatus{
		TaskId:    taskId,
		JobId:     jobId,
		Timestamp: int64(time.Now().Unix()),
	}

	// Check if the container exists
	if ContainerExistByName(jobId) {
		// restart container
		err := Client.RestartContainer(jobId, conf.Conf.Interval.RestartContainerTimeout)
		if err != nil {
			msg := err.Error()
			status.Code = conf.RestartContainerFailedCode
			status.Status = conf.RestartFailure
			status.Msg = msg
			// send restart container failed status
			jobCh <- status
			tools.Logger.Error(msg)
		} else {
			status.Code = conf.RestartContainerSuccessCode
			status.Status = conf.RestartSuccess
			// send restart container failed status
			jobCh <- status
		}
	} else {
		var job tools.EdgeCmJob
		db := tools.First(&job, "jobid = ?", jobId)
		if db.Error != nil {
			msg := db.Error.Error()
			status.Code = conf.RestartContainerFailedCode
			status.Status = conf.RestartFailure
			status.Msg = msg
			// send restart container failed status
			jobCh <- status
			tools.Logger.Error(msg)
			return
		}
		if CheckImageExistByName(jobId) {
			go StartContainer(&job, jobCh)
		} else {
			go BuildImage(true, &job, jobCh)
		}
		status.Code = conf.RestartContainerSuccessCode
		status.Status = conf.RestartSuccess
		// send restart container failed status
		jobCh <- status
	}
}

// get containers by filters
func Containers(filters map[string][]string) ([]docker.APIContainers, error) {
	opt := docker.ListContainersOptions{Filters: filters}
	return Client.ListContainers(opt)
}

// check container exist by name
func ContainerExistByName(name string) bool {
	cond := fmt.Sprintf("jobid=%s", name)
	filters := map[string][]string{"label": []string{cond}}
	containers, err := Containers(filters)
	if err != nil {
		tools.Logger.Error(fmt.Sprintf("%s", err.Error()))
		return false
	}
	if len(containers) > 0 {
		return true
	}
	return false
}

// get current running jobs
func RunningJob() ([]map[string]string, error) {
	var jobs []map[string]string

	// search status=running metadata label job=true jobs
	filters := map[string][]string{"status": []string{"running"}, "label": []string{"job=true"}}

	containers, err := Containers(filters)
	if err != nil {
		tools.Logger.Error(fmt.Sprintf("%s", err.Error()))
		return nil, err
	}
	for _, c := range containers {
		labels := c.Labels
		// delete key job
		jobMap := make(map[string]string)
		jobMap["jobid"] = labels["jobid"]
		jobMap["taskid"] = labels["taskid"]
		//delete(labels, "job")
		jobs = append(jobs, jobMap)
	}
	return jobs, nil
}

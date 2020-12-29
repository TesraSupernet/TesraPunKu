package service

import (
	"fmt"
	"github.com/DOSNetwork/core/edge_cm/conf"
	"github.com/DOSNetwork/core/edge_cm/tools"
	"github.com/DOSNetwork/core/edge_cm/tools/mq"
	"time"
)

// Check to send the failed message and resend it
func CheckFailedMessageAndResend() {
	for {
		time.Sleep(time.Duration(conf.Conf.Interval.ResendMessageInterval) * time.Second)
		var failedMessages []*tools.FailedMessage
		db := tools.Find(&failedMessages)
		if db.Error != nil {
			tools.Logger.Error(db.Error.Error())
			continue
		}
		for _, msg := range failedMessages {
			_, err := mq.Send(false, msg.Topic, msg.Message)
			if err == nil {
				db := tools.Delete(msg)
				if db.Error != nil {
					tools.Logger.Error(db.Error.Error())
				}
			}
		}
	}
}

// check containers
func CheckContainerStatus() {
	for {
		time.Sleep(time.Second * time.Duration(conf.Conf.Container.CheckLive))
		var jobs []*tools.EdgeCmJob
		db := tools.Find(&jobs, "job_status in (?)", []int{conf.StartContainerSuccess, conf.Decomposing,
			conf.Decomposed, conf.Calculating, conf.Calculated, conf.Uploading, conf.DataSetDownloading,
			conf.DataSetDownloaded, conf.DataSetDecompressSuccess, conf.DataSetDecompression})
		if db.Error != nil {
			tools.Logger.Error(fmt.Sprintf("database find error: %s", db.Error.Error()))
			continue
		}
		tools.Logger.Debug(fmt.Sprintf("all running jobs: %+v", jobs))
		time.Sleep(time.Second * time.Duration(conf.Conf.Container.CheckLive)) // TODO 等待数据更新数据库
		for _, job := range jobs {
			filter := []string{fmt.Sprintf("jobid=%s", job.JobId)}
			containers, err := Containers(map[string][]string{"label": filter})
			if err != nil {
				tools.Logger.Error(fmt.Sprintf("get container list error: %s", err.Error()))
				continue
			}
			if len(containers) == 0 {
				var edgeJob tools.EdgeCmJob
				db := tools.First(&edgeJob, "jobid = ?", job.JobId)
				if db.Error != nil {
					tools.Logger.Error(db.Error.Error())
					continue
				}
				if edgeJob.JobStatus == job.JobStatus {
					status := &tools.JobStatus{
						TaskId:    job.TaskId,
						JobId:     job.JobId,
						Timestamp: int64(time.Now().Unix()),
						Msg:       "container abnormal exit.",
						Code:      conf.ContainerExitedCode,
						Status:    conf.ContainerExited,
					}
					statusCh <- status
					edgeJob.JobStatus = conf.ContainerExited
					tools.Logger.Debug(fmt.Sprintf("container exited, jobid=%s;status=%d", edgeJob.JobId, edgeJob.JobStatus))
					tools.Update(edgeJob)
				}
			}
		}
	}
}

// Check the image build or the task that the image build successfully hangs
func CheckPendingJob() {
	for {
		time.Sleep(time.Second * time.Duration(conf.Conf.Container.CheckLive))
		var jobs []*tools.EdgeCmJob
		db := tools.Find(&jobs, "job_status in (?)", []int{conf.ImageBuilding, conf.ImageBuilt})
		if db.Error != nil {
			tools.Logger.Error(db.Error.Error())
			continue
		}
		now := int64(time.Now().Unix())
		for _, job := range jobs {
			if now-job.ReceiveTimestamp > conf.Conf.Image.MaxDealTimeout {
				tools.Logger.Info(fmt.Sprintf("job %s timeout, status=%d", job.JobId, job.JobStatus))
				status := &tools.JobStatus{
					TaskId:    job.TaskId,
					JobId:     job.JobId,
					Timestamp: int64(time.Now().Unix()),
				}
				switch job.JobStatus {
				case conf.ImageBuilding:
					status.Code = conf.ImageBuildFailedCode
					status.Status = conf.ImageBuildFailure
				case conf.ImageBuilt:
					status.Code = conf.StartContainerFailedCode
					status.Status = conf.StartContainerFailure
				}
				statusCh <- status
			}
		}
	}
}

// system check
func SystemCheck() {
	tools.Logger.Info("system checking...")
	jobs, err := RunningJob()
	var tmp = make(map[string]struct{})
	if err != nil {
		tools.Logger.Error(fmt.Sprintf("get running jobs error: %s", err.Error()))
	}

	for _, job := range jobs {
		tmp[job["jobid"]] = struct{}{}
	}

	for jobid, _ := range tmp {
		tools.Logger.Debug(fmt.Sprintf("laod job %s success.", jobid))
		// store cwa to sqlite
		dispatcher := StoreCWA(jobid)
		// save to JobDispatchers
		JobDispatchers.Store(jobid, dispatcher)
	}
}

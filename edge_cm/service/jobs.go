package service

import (
	"encoding/json"
	"fmt"
	"github.com/DOSNetwork/core/edge_cm/conf"
	"github.com/DOSNetwork/core/edge_cm/tools"
	"github.com/DOSNetwork/core/edge_cm/tools/mq"
	"time"
)

var statusCh = make(chan *tools.JobStatus)
var IsRunningJob = make(chan bool)

// pull task information
func PullJobImageBuildRun() {
	mac := tools.GetMacByName(conf.Conf.MMLServer.MAC_ADDR)
	mc := mq.Receive(conf.JobInfoPullTopic, mac, "shared")
	for msg := range mc.C {
		var edgeCmJob tools.EdgeCmJob
		payload := msg.Payload()
		err := json.Unmarshal(payload, &edgeCmJob)
		if err != nil {
			tools.Logger.Error(fmt.Sprintf("%s", err.Error()))
		}
		if edgeCmJob.Mac == mac {
			e := &tools.EdgeCmJob{
				Mac:              edgeCmJob.Mac,
				RDS:              edgeCmJob.RDS,
				DDM:              edgeCmJob.DDM,
				DDS:              edgeCmJob.DDS,
				NCP:              edgeCmJob.NCP,
				Type:             edgeCmJob.Type,
				Index:            edgeCmJob.Index,
				Count:            edgeCmJob.Count,
				JobId:            edgeCmJob.JobId,
				TaskId:           edgeCmJob.TaskId,
				NcpArgs:          edgeCmJob.NcpArgs,
				JobStatus:        edgeCmJob.JobStatus,
				Dockerfile:       edgeCmJob.Dockerfile,
				NcpCommand:       edgeCmJob.NcpCommand,
				DsmCommand:       edgeCmJob.DsmCommand,
				ReceiveTimestamp: int64(time.Now().Unix()),
			}
			db := tools.Insert(e)
			if db.Error != nil {
				// send failed status
				jobStatus := &tools.JobStatus{
					Msg:    db.Error.Error(),
					Code:   conf.ImageBuildFailedCode,
					Status: conf.ImageBuildFailure,
					TaskId: edgeCmJob.TaskId,
					JobId:  edgeCmJob.JobId,
				}
				buf, err := json.Marshal(jobStatus)
				if err != nil {
					tools.Logger.Error(fmt.Sprintf("marshal job status error: %s", err.Error()))
				}
				_, err = mq.Send(false, conf.JobStatusPushTopic, buf)
				if err != nil {
					tools.Logger.Info(fmt.Sprintf("send status failed: %s.", err.Error()))
					db := tools.Insert(&tools.FailedMessage{Message: buf, Topic: conf.JobStatusPushTopic})
					if db.Error != nil {
						tools.Logger.Error(fmt.Sprintf("insert send failed message error: %s", db.Error.Error()))
					}
				}
			} else {
				// build image
				go BuildImage(true, &edgeCmJob, statusCh)
			}
		}
	}
}

// pull job command
func PullJobCommand() {
	mac := tools.GetMacByName(conf.Conf.MMLServer.MAC_ADDR)
	mc := mq.Receive(conf.JobCommandPullTopic, mac, "shared")
	for msg := range mc.C {
		var jobCommand JobCommand
		payload := msg.Payload()
		err := json.Unmarshal(payload, &jobCommand)
		if err != nil {
			tools.Logger.Error(fmt.Sprintf("%s", err.Error()))
		}
		if jobCommand.Mac == mac {
			switch jobCommand.Command {
			case conf.START:
				// TODO start命令不暂不处理
			case conf.STOP:
				// stop container
				go StopContainer(jobCommand.TaskId, jobCommand.JobId, statusCh)
			case conf.RESTART:
				// restart container
				go RestartContainer(jobCommand.TaskId, jobCommand.JobId, statusCh)
			case conf.DELETE:
				go DeleteNodeData(jobCommand.JobId)
			}
		}
	}
}

// pull DSM job status
func PullJobStatus() {
	mac := tools.GetMacByName(conf.Conf.MMLServer.MAC_ADDR)
	mc := mq.Receive(conf.JobStatusPullTopic, mac, "shared")
	for msg := range mc.C {
		var status tools.JobStatus
		payload := msg.Payload()
		err := json.Unmarshal(payload, &status)
		if err != nil {
			tools.Logger.Error(fmt.Sprintf("%s", err.Error()))
		}
		if status.Mac == mac {
			statusCh <- &status
		}
	}
}

// push hardware information every 30 seconds
func PushHardwareInfo() {
	h := tools.HardwareGenerator(time.Second * time.Duration(conf.Conf.Interval.PushHardwareInterval))
	for hw := range h.C {
		buf, _ := json.Marshal(hw)
		status, err := mq.Send(false, conf.HardwarePushTopic, buf)
		if err != nil {
			info := fmt.Sprintf("push hardware:\n information: %s\n status:%t\n error: %v", string(buf), status, err)
			tools.Logger.Debug(info)
		}
	}
}

// push running jobs
func PushRunningJob(forever bool) {
	jobObj := make(map[string]interface{})
	mac := tools.GetMacByName(conf.Conf.MMLServer.MAC_ADDR)
	jobObj["mac"] = mac
	for {
		// get current running jobs
		jobs, _ := RunningJob()
		if jobs != nil {
			IsRunningJob <- true
			jobObj["jobs"] = jobs
		} else {
			IsRunningJob <- false
			jobObj["jobs"] = nil
		}
		buf, _ := json.Marshal(jobObj)
		tools.Logger.Debug(string(buf))
		_, e := mq.Send(false, conf.RunningJobPushTopic, buf)
		if e != nil {
			tools.Logger.Error("send running job failed.")
		}
		if !forever {
			break
		}
		time.Sleep(time.Second * time.Duration(conf.Conf.Interval.PushRunningJobInterval))
	}
}

// push job status
func PushJobStatusChanged() {
	for {
		// listen statusCh send job status
		select {
		case status := <-statusCh:
			switch status.Status {
			// The final state will send stop peak statistics and send
			case conf.DecompositionFailure, conf.DataSetDownloadFailure, conf.DataSetDecompressFailure,
				conf.CalculateFailure, conf.UploadFailure, conf.Uploaded, conf.StopSuccess, conf.SystemError,
				conf.ContainerExited:
				if v, ok := JobDispatchers.Load(status.JobId); !ok {
					tools.Logger.Error(fmt.Sprintf("key %s not found!", status.JobId))
				} else {
					if dispatcher, ok := v.(*Dispatcher); !ok {
						tools.Logger.Error("type error.")
					} else {
						dispatcher.Stop()
					}
				}
				peakInfo, err := Peak(status.JobId)
				if err != nil {
					tools.Logger.Error(err.Error())
				} else {
					status.PeakInfo = peakInfo
				}
				// delete this job dispatcher
				JobDispatchers.Delete(status.JobId)
				// delete this job cwa
				DeleteCWA(status.JobId)
			}

			// update database job_status
			db := tools.Updates(&tools.EdgeCmJob{}, &tools.EdgeCmJob{JobId: status.JobId}, &tools.EdgeCmJob{JobStatus: status.Status})
			if db.Error != nil {
				tools.Logger.Error(fmt.Sprintf("%s", db.Error.Error()))
			}
			buf, _ := json.Marshal(status)
			tools.Logger.Debug(string(buf))
			_, e := mq.Send(false, conf.JobStatusPushTopic, buf)
			if e != nil {
				tools.Logger.Error("send status failed.")
				db := tools.Insert(&tools.FailedMessage{Message: buf, Topic: conf.JobStatusPushTopic})
				if db.Error != nil {
					tools.Logger.Error(fmt.Sprintf("insert send failed message error: %s", db.Error.Error()))
				}
			}
		}
	}
}

package service

import (
	"encoding/json"
	"fmt"
	"github.com/DOSNetwork/core/edge_cm/conf"
	"github.com/DOSNetwork/core/edge_cm/tools"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

// request body
type ReqParam struct {
	TaskId string
}

// response body
type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// api response
type ResponseInfo struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Code    int    `json:"code"`
	Data    tools.EdgeCmJob
}

// build http response body
func Response(c *gin.Context, statusCode int, resp interface{}) {
	c.JSON(statusCode, resp)
}

// handle request data parse struct
func HandleRequestData(c *gin.Context, out interface{}) bool {
	buf, err := c.GetRawData()
	if err != nil {
		resp := Resp{
			Code: conf.GetDockerfileFailed,
			Msg:  err.Error(),
		}
		Response(c, http.StatusBadRequest, resp)
		return false
	}
	err = json.Unmarshal(buf, &out)
	if err != nil {
		resp := Resp{
			Code: conf.GetDockerfileFailed,
			Msg:  err.Error(),
		}
		Response(c, http.StatusBadRequest, resp)
		return false
	}
	return true
}

// call api and store job information result
func GetAndStoreJobInfo(url, taskId string, r *Resp) {
	url = fmt.Sprintf("%s%s", url, taskId)
	resp, err := http.Get(url)
	if err != nil {
		tools.Logger.Error(fmt.Sprintf("containermt error: %s.\n", err.Error()))
		return
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	var responseInfo ResponseInfo
	err = json.Unmarshal(body, &responseInfo)
	if err != nil {
		return
	}
	if resp.StatusCode == 200 && responseInfo.Code == 200 {
		responseInfo.Data.JobId = tools.UUID()
		responseInfo.Data.Mac = tools.GetMacByName(conf.Conf.MMLServer.MAC_ADDR)
		db := tools.Insert(&responseInfo.Data)
		if db.Error != nil {
			r.Code = conf.GetDockerfileFailed
			r.Msg = db.Error.Error()
		} else {
			r.Code = conf.GetDockerfileSuccess
			r.Msg = "Dockerfile downloaded."
		}
	} else {
		r.Code = conf.GetDockerfileFailed
		r.Msg = "task does't not exists."
	}

}

// get dockerfile by taskid
func GetDockerfile(c *gin.Context) {
	var reqParam ReqParam
	if !HandleRequestData(c, &reqParam) {
		return
	}
	var job tools.EdgeCmJob
	db := tools.First(&job, "taskid = ?", reqParam.TaskId)
	if db.Error != nil {
		var r Resp
		GetAndStoreJobInfo(conf.Conf.Interface.JOB_INFO_URL, reqParam.TaskId, &r)
		Response(c, http.StatusOK, r)
	} else {
		resp := Resp{
			Code: conf.GetDockerfileSuccess,
			Msg:  "Dockerfile downloaded.",
		}
		Response(c, http.StatusOK, resp)
	}
}

// pull image by taskid
func PullImage(c *gin.Context) {
	var reqParam ReqParam
	if !HandleRequestData(c, &reqParam) {
		return
	}
	var job tools.EdgeCmJob
	db := tools.First(&job, "taskid = ?", reqParam.TaskId)
	if db.Error != nil {
		resp := Resp{
			Code: conf.ImageBuildFailedCode,
			Msg:  db.Error.Error(),
		}
		Response(c, http.StatusBadRequest, resp)
		return
	}
	BuildImage(false, &job, statusCh)
	resp := Resp{
		Code: conf.ImageBuildingCode,
		Msg:  "Image building...",
	}
	Response(c, http.StatusOK, resp)
}

// delete image by taskid
func DelImage(c *gin.Context) {
	var reqParam ReqParam
	if !HandleRequestData(c, &reqParam) {
		return
	}
	var job tools.EdgeCmJob
	db := tools.First(&job, "taskid = ?", reqParam.TaskId)
	if db.Error != nil {
		resp := Resp{
			Code: conf.RemoveImageFailedCode,
			Msg:  db.Error.Error(),
		}
		Response(c, http.StatusBadRequest, resp)
		return
	}
	err := RemoveImage(job.JobId)
	if err != nil {
		resp := Resp{
			Code: conf.RemoveImageFailedCode,
			Msg:  err.Error(),
		}
		Response(c, http.StatusBadRequest, resp)
		return
	}
	resp := Resp{
		Code: conf.RemoveImageSuccessCode,
		Msg:  "Image remove success.",
	}
	Response(c, http.StatusOK, resp)
}

// run container by taskid
func RunContainer(c *gin.Context) {
	var reqParam ReqParam
	if !HandleRequestData(c, &reqParam) {
		return
	}
	var job tools.EdgeCmJob
	db := tools.First(&job, "taskid = ?", reqParam.TaskId)
	if db.Error != nil {
		resp := Resp{
			Code: conf.StartContainerFailedCode,
			Msg:  db.Error.Error(),
		}
		Response(c, http.StatusBadRequest, resp)
		return
	}
	StartContainer(&job, statusCh)
	resp := Resp{
		Code: conf.StartContainerSuccessCode,
		Msg:  "The container is being starting.",
	}
	Response(c, http.StatusOK, resp)
}

// stop container by taskid
func ContainerStop(c *gin.Context) {
	var reqParam ReqParam
	if !HandleRequestData(c, &reqParam) {
		return
	}
	var job tools.EdgeCmJob
	db := tools.First(&job, "taskid = ?", reqParam.TaskId)
	if db.Error != nil {
		resp := Resp{
			Code: conf.StartContainerFailedCode,
			Msg:  db.Error.Error(),
		}
		Response(c, http.StatusBadRequest, resp)
		return
	}
	StopContainer(job.TaskId, job.JobId, statusCh)
	resp := Resp{
		Code: conf.StopContainerSuccessCode,
		Msg:  "The container stopped.",
	}
	Response(c, http.StatusOK, resp)
}

// restart container by taskid
func ContainerRestart(c *gin.Context) {
	var reqParam ReqParam
	if !HandleRequestData(c, &reqParam) {
		return
	}
	var job tools.EdgeCmJob
	db := tools.First(&job, "taskid = ?", reqParam.TaskId)
	if db.Error != nil {
		resp := Resp{
			Code: conf.RestartContainerFailedCode,
			Msg:  db.Error.Error(),
		}
		Response(c, http.StatusBadRequest, resp)
		return
	}
	RestartContainer(job.TaskId, job.JobId, statusCh)
	resp := Resp{
		Code: conf.RestartContainerSuccessCode,
		Msg:  "The container restarted.",
	}
	Response(c, http.StatusOK, resp)
}

func Serve() {
	//router := gin.Default()
	router := gin.New()
	//router.Use()
	router.POST("/get-dockerfile", GetDockerfile)
	router.POST("/pull-image", PullImage)
	router.POST("/del-image", DelImage)
	router.POST("/run-container", RunContainer)
	router.POST("/stop-container", ContainerStop)
	router.POST("/restart-container", ContainerRestart)
	router.Run(conf.Conf.MMLServer.SERVER_ADDR)
}

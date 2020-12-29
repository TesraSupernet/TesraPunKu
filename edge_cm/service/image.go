package service

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"github.com/DOSNetwork/core/edge_cm/conf"
	"github.com/DOSNetwork/core/edge_cm/tools"
	docker "github.com/fsouza/go-dockerclient"
	"time"
)

var Client *docker.Client

func init() {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		tools.Logger.Fatal(fmt.Sprintf("%s", err.Error()))
	}
	Client = client
}

func BuildDockerfile(dockerFile string) (inputBuf *bytes.Buffer, err error) {
	inputBuf = bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputBuf)
	defer func() {
		err = tr.Close()
		if err != nil {
			return
		}
	}()
	currentTime := time.Now()
	err = tr.WriteHeader(&tar.Header{
		Name:       "Dockerfile",
		Size:       int64(len(dockerFile)),
		ModTime:    currentTime,
		AccessTime: currentTime,
		ChangeTime: currentTime,
	})
	if err != nil {
		return
	}
	_, err = tr.Write([]byte(dockerFile))

	if err != nil {
		return
	}
	return
}

// Generate a image and deduce the state based on the job information
func BuildImage(auto bool, job *tools.EdgeCmJob, jobCh chan<- *tools.JobStatus) {
	status := &tools.JobStatus{
		TaskId:    job.TaskId,
		JobId:     job.JobId,
		Timestamp: int64(time.Now().Unix()),
	}

	t := time.Now()
	inputBuf, outputBuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputBuf)
	err := tr.WriteHeader(&tar.Header{
		Name:       "Dockerfile",
		Size:       int64(len(job.Dockerfile)),
		ModTime:    t,
		AccessTime: t,
		ChangeTime: t,
	})
	if err != nil {
		tools.Logger.Error(fmt.Sprintf("%s", err.Error()))
		status.Msg = err.Error()
		status.Code = conf.ImageBuildFailedCode
		status.Status = conf.ImageBuildFailure
		// send build images failed status
		jobCh <- status
		return
	}

	_, err = tr.Write([]byte(job.Dockerfile))

	if err != nil {
		tools.Logger.Error(fmt.Sprintf("%s", err.Error()))
		status.Msg = err.Error()
		status.Code = conf.ImageBuildFailedCode
		status.Status = conf.ImageBuildFailure
		// send build image failed status
		jobCh <- status
		return
	}
	err = tr.Close()
	if err != nil {
		tools.Logger.Error(fmt.Sprintf("%s", err.Error()))
		status.Msg = err.Error()
		status.Code = conf.ImageBuildFailedCode
		status.Status = conf.ImageBuildFailure
		jobCh <- status
		return
	}
	// Generate image options
	opts := docker.BuildImageOptions{
		Context:      context.Background(),
		Name:         job.JobId,
		InputStream:  inputBuf,
		OutputStream: outputBuf,
		//NoCache:             true,
		RmTmpContainer:      true,
		ForceRmTmpContainer: true,
		Labels:              map[string]string{"job": "true", "taskid": job.TaskId, "jobid": job.JobId}, // set meta data
	}
	status.Msg = "image building..."
	status.Code = conf.ImageBuildingCode
	status.Status = conf.ImageBuilding
	// send image building status
	jobCh <- status
	tools.Logger.Info("start build image...")

	// send running job
	PushRunningJob(false)

	// Generate image based on dockerfile
	if err := Client.BuildImage(opts); err != nil {
		tools.Logger.Error(fmt.Sprintf("%s", err.Error()))
		status.Msg = err.Error()
		status.Code = conf.ImageBuildFailedCode
		status.Status = conf.ImageBuildFailure
		// send build images failed status
		jobCh <- status
		return
	}
	tools.Logger.Info("images build success.")
	status.Msg = "image build success."
	status.Code = conf.ImageBuildSuccessCode
	status.Status = conf.ImageBuilt
	// send images built status
	jobCh <- status

	//r := make([]byte, 1024*10)
	//// Read the build process information
	//for {
	//	n, err := outputBuf.Read(r)
	//	if err != nil {
	//		if err == io.EOF {
	//			break
	//		} else {
	//			ztools.Logger.Error(fmt.Sprintf("%s", err.Error()))
	//		}
	//	}
	//	ztools.Logger.Info(fmt.Sprintf("\n%s", string(r[:n])))
	//}
	// auto start container
	if auto {
		// start container
		StartContainer(job, jobCh)
	}

}

// check images and auto remove images
func CheckAndAutoRmImage() {
	for {
		time.Sleep(time.Second * time.Duration(conf.Conf.Image.AutoRmImgInterval))
		if images, err := JobImageList(); err == nil && len(images) >= conf.Conf.Image.MaxImages {
			var jobs []*tools.EdgeCmJob
			db := tools.Find(&jobs, "job_status in (?)", []int{
				conf.StartContainerFailure, conf.Uploaded, conf.StopSuccess, conf.DecompositionFailure,
				conf.CalculateFailure, conf.UploadFailure, conf.DataSetDownloadFailure, conf.SystemError})
			if db.Error != nil {
				tools.Logger.Error(fmt.Sprintf("%s", db.Error.Error()))
				continue
			}
			tools.Logger.Debug(fmt.Sprintf("stoped jobs: %+v", jobs))
			for _, job := range jobs {
				err := RemoveImage(job.JobId)
				if err != nil {
					tools.Logger.Error(fmt.Sprintf("remove image %s failed, error: %s.", job.JobId, err.Error()))
				} else {
					tools.Logger.Info(fmt.Sprintf("remove image %s success.", job.JobId))
				}
				db := tools.Delete(&tools.EdgeCmJob{}, "jobid= ? ", job.JobId)
				if db.Error != nil {
					tools.Logger.Error(fmt.Sprintf("delete record %s error: %s", job.JobId, db.Error.Error()))
				}
			}
		}
	}
}

// clear invalid images
func ClearDanglingImage() {
	for {
		var count int
		db := tools.Count(&tools.EdgeCmJob{}, &count, "job_status = ?", conf.ImageBuilding)
		if db.Error != nil || 0 != count {
			tools.Logger.Error(fmt.Sprintf("get building image count %d, error: %+v", count, db.Error))
			time.Sleep(time.Second * time.Duration(conf.Conf.Image.AutoRmImgInterval))
			continue
		}
		filter := map[string][]string{"dangling": []string{"true"}, "label": []string{"job=true"}}
		//filter := map[string][]string{"dangling": []string{"true"}}
		images, err := Images(filter)
		if err != nil {
			tools.Logger.Error(fmt.Sprintf("%s", err))
			time.Sleep(time.Second * time.Duration(conf.Conf.Image.AutoRmImgInterval))
			continue
		}
		for _, image := range images {
			e := RemoveImage(image.ID)
			if e != nil {
				tools.Logger.Error(fmt.Sprintf("%s", e))
				continue
			}
			tools.Logger.Debug(fmt.Sprintf("remove image %s success.", image.ID))
		}
		time.Sleep(time.Second * time.Duration(conf.Conf.Image.AutoRmImgInterval))
	}
}

// search images by filters
func Images(filters map[string][]string) ([]docker.APIImages, error) {
	return Client.ListImages(docker.ListImagesOptions{Filters: filters})
}

// check images exist by name
func CheckImageExistByName(name string) bool {
	cond := fmt.Sprintf("jobid=%s", name)
	filters := map[string][]string{"label": []string{cond}}
	images, err := Images(filters)
	if err != nil {
		tools.Logger.Error(fmt.Sprintf("%s", err.Error()))
		return false
	}
	if len(images) > 0 {
		return true
	}
	return false
}

// get job images
func JobImageList() ([]docker.APIImages, error) {
	filters := map[string][]string{"label": []string{"job=true"}}
	return Images(filters)
}

// remove images
func RemoveImage(name string) error {
	//return Client.RemoveImage(name)
	return Client.RemoveImageExtended(name, docker.RemoveImageOptions{Force: true})
}

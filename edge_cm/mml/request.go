package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	GET_DOCKERFILE_URL    = "http://localhost:9090/get-dockerfile"
	PULL_IMAGE_URL        = "http://localhost:9090/pull-image"
	DEL_IMAGE_URL         = "http://localhost:9090/del-image"
	RUN_CONTAINER_URL     = "http://localhost:9090/run-container"
	STOP_CONTAINER_URL    = "http://localhost:9090/stop-container"
	RESTART_CONTAINER_URL = "http://localhost:9090/restart-container"
)

// deal with request
func Request(taskId, url string) {
	data := map[string]string{"taskid": taskId}
	buf, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewReader(buf))
	if err != nil {
		fmt.Printf("containermt error: %s.\n", err.Error())
		return
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

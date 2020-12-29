package ztools

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpMethod int

const (
	Get HttpMethod = iota
	Post
	Put
	Delete
)

var httpCli *http.Client

func init() {
	var trans http.Transport
	trans = http.Transport{
		DisableKeepAlives: true,
	}
	httpCli = &http.Client{
		Transport: &trans,
	}
}

//将数据发送出去
// HttpRequest ...
func HttpRequest(url, ContentType string, method HttpMethod, v io.Reader, header map[string]string) ([]byte, error) {
	var (
		err  error
		resp *http.Response
		req  *http.Request
		body []byte
	)

	headerSet := func(req *http.Request) {
		if len(header) != 0 {
			for k, v := range header {
				req.Header.Set(k, v)
			}
		}
		return
	}
	err = nil
	switch method {
	case Get:
		req, err = http.NewRequest("GET", url, v)
		if err != nil {
			break
		}
		headerSet(req)

		break
	case Post:
		req, err = http.NewRequest("POST", url, v)
		if err != nil {
			break
		}
		headerSet(req)
		req.Header.Set("Content-Type", ContentType)

		break
	case Put:
		req, err = http.NewRequest("PUT", url, v)
		if err != nil {
			break
		}
		headerSet(req)
		req.Header.Set("Content-Type", ContentType)
		break
	case Delete:
		req, err = http.NewRequest("PUT", url, v)
		if err != nil {
			break
		}
		headerSet(req)
		req.Header.Set("Content-Type", ContentType)
		break
	default:
		return nil, errors.New("no http type")
	}
	//
	if err != nil {
		return nil, err
	}

	//
	resp, err = httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//defer req.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

package tools

import (
	"github.com/DOSNetwork/core/edge_cm/conf"
	"github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
	"net/url"
	"sync"
	"time"
)

var minioClient *minio.Client
var lock sync.Mutex

func Client() *minio.Client {
	if minioClient == nil {
		lock.Lock()
		endpoint := conf.Conf.Minio.Endpoint
		accessKeyID := conf.Conf.Minio.Access
		secretAccessKey := conf.Conf.Minio.Secret
		useSSL := conf.Conf.Minio.UseSsl
		client, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
		if err != nil {
			log.Panicf("create minio client error: %s", err.Error())
		}
		minioClient = client
		lock.Unlock()
	}
	return minioClient
}

func GetTempDownloadPath(bucketName string, objectName string) (tempDownloadPath *url.URL) {
	tempDownloadPath, err := Client().PresignedGetObject(bucketName, objectName, time.Second*2*60*60, nil)
	if err != nil {
		log.Panicf("get temp download path: %s", err.Error())
	}
	return
}

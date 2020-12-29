package ztools

import (
	"crypto/md5"
	"encoding/base64"
	"io"
	"os"
)

//获取经过base64处理的md5值
func GetFileMd5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return "", err
	}

	md5Sum := md5hash.Sum(nil)
	str := base64.StdEncoding.EncodeToString(md5Sum)

	return str, nil
}

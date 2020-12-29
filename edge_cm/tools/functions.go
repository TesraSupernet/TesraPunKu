package tools

import (
	"bufio"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"net"
	"os"
	"strings"
)

// get mac address
func GetMacByName(name string) string {
	inter, err := net.InterfaceByName(name)
	if err != nil {
		log.Println(err)
		return ""
	}
	return inter.HardwareAddr.String()
}

func UUID() string {
	u := uuid.NewV4()
	return u.String()
}

func CheckSuffix(path string) string {
	sep := string(os.PathSeparator)
	if !strings.HasSuffix(path, sep) {
		path = fmt.Sprintf("%s%s", path, sep)
	}
	return path
}

func MakePath(path string) string {
	f, err := os.Stat(path)
	if err != nil {
		e := os.Mkdir(path, 0700)
		if e != nil {
			log.Panic(e)
		}
		return CheckSuffix(path)
	}
	if !f.IsDir() {
		log.Panicf("%s is not dir.", path)
	}
	return CheckSuffix(path)
}

// write string to file
func WriteFile(prefix, filename, content string) error {
	if prefix == "" {
		prefix = "/home/files/"
	}
	file := fmt.Sprintf("%s%s", prefix, filename)
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		Logger.Error(fmt.Sprintf("%s", err.Error()))
		return err
	}
	defer f.Close()
	write := bufio.NewWriter(f)
	write.WriteString(content)
	write.Flush()
	return nil
}

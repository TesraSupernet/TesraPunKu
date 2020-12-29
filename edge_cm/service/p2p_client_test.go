package service

import (
	"fmt"
	_ "github.com/alecthomas/template"
	"testing"
)

func TestVolumes(t *testing.T) {
	filter := []string{"serverName=p2pClient"}
	containers, err := Containers(map[string][]string{"label": filter})
	if err != nil {
		fmt.Print(err)
	}
	for _, container := range containers {
		mounts := container.Mounts
		for _, mount := range mounts {
			fmt.Print(mount)
		}
	}
}

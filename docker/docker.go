package main

import (
	"github.com/vrgakos/livemigrate/ssh"
	"encoding/json"
	"fmt"
)

type InspectResult struct {
	Name string
	State struct {
		Running bool
	}
	Config struct {
		Image string
	}
	HostConfig struct {
		Binds []string
	}
	NetworkSettings struct {
		Networks map[string] struct{
			IPAMConfig struct{
				IPv4Address string
			}
		}
	}
}

func InspectDocker(client *ssh.SshClient, containerId string) *InspectResult {
	var ins InspectResult

	result, err := client.RunAndWait("docker inspect -f '{{ json . }}' " + containerId)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if err := json.Unmarshal(result, &ins); err != nil {
		panic(err)
	}
	return &ins
}

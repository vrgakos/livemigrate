package node

import (
	"github.com/docker/docker/client"
	"fmt"
)

type Node struct {
	Alias		string
	Host		string

	SshPort		int
	SshUser		string
	SshKey		string

	DockerApiPort	int

	sshClient	*SshClient
	dockerClient	*client.Client
}

func (n *Node) GetSshClient() (*SshClient, error) {
	if n.sshClient == nil {
		var err error
		n.sshClient, err = NewSshClient(n.Host, n.SshPort, n.SshUser, n.SshKey)
		if err != nil {
			return nil, err
		}
	}

	return n.sshClient, nil
}

func (n *Node) GetDockerClient() (*client.Client, error) {
	if n.dockerClient == nil {
		var err error
		headers := map[string]string { "User-Agent": "livemigrate-0.1" }
		n.dockerClient, err = client.NewClient(fmt.Sprintf("tcp://%s:%d", n.Host, n.DockerApiPort), "1.29", nil, headers)
		if err != nil {
			return nil, err
		}
	}

	return n.dockerClient, nil
}
package main

import (
	"github.com/vrgakos/livemigrate/ssh"
	"github.com/google/uuid"
)

type TcpDump struct {
	client		*ssh.SshClient
	iface		string
	fileName	string
	remoteFileName	string

	job		*ssh.BackGroundjob
}

func (t *TcpDump) Run() {
	t.job = t.client.NewBackGroundjob("tcpdump -s 0 -i " + t.iface + " -w " + t.remoteFileName + " 2>/dev/null")
}

func (t *TcpDump) Stop() {
	t.job.Stop()
	t.client.RunAndWriteFile("cat " + t.remoteFileName + " && rm " + t.remoteFileName, t.fileName)
}

func NewTcpdump(client *ssh.SshClient, iface string, fileName string) *TcpDump {
	id, _ := uuid.NewRandom()
	randomId := id.String()

	tcpDump := &TcpDump{
		client:		client,
		iface:		iface,
		fileName:	fileName,
		remoteFileName:	"/tmp/" + randomId + ".pcap",
	}
	tcpDump.Run()
	return tcpDump
}

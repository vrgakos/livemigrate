package ssh

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"fmt"
	"bytes"
	"log"
	"io"
)

func PublicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to read key file: %s", err)
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse key file: %s", err)
	}

	return ssh.PublicKeys(key), nil
}


type SshClient struct {
	Alias		string
	Host		string
	Port 		int
	User		string

	config		*ssh.ClientConfig
	conn		*ssh.Client
}

func NewSshClient(alias string, host string, port int, user string, keyFile string) (*SshClient, error) {
	auth, err := PublicKeyFile(keyFile)
	if err != nil {
		return nil, err
	}

	client := &SshClient{
		Alias:		alias,
		Host:		host,
		Port:		port,
		User:		user,

		config:		&ssh.ClientConfig{
					User: user,
					Auth: []ssh.AuthMethod { auth },
				},
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port), client.config)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial: %s", err)
	}


	client.conn = conn

	return client, nil
}

func (c *SshClient) RunAndWait(command string) ([]byte, error) {
	session, err := c.conn.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %s", err)
	}


	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	var stderrBuf bytes.Buffer
	session.Stderr = &stderrBuf

	err = session.Run(command)
	if err != nil {
		return nil, fmt.Errorf("Failed to run command: %s\n%s", err, stderrBuf.String())
	}

	session.Close()

	return stdoutBuf.Bytes(), nil
}

func (c *SshClient) RunAndWriteFile(command, fileName string) error {
	res, err := c.RunAndWait(command)
	if (err != nil) {
		log.Println("Error at running command:", err)
		return err
	}

	err = ioutil.WriteFile(fileName, res, 777)
	if (err != nil) {
		log.Println("Error at writing output file:", err)
		return err
	}
	return nil
}

type BackGroundjob struct {
	command		string
	client		*SshClient
	session		*ssh.Session

	result		[]byte
	err		error
	done		chan bool

	stdin		io.WriteCloser
}

func (j *BackGroundjob) Run() {
	var err error
	j.session, err = j.client.conn.NewSession()
	if err != nil {
		j.err = err
		j.done <- false
		return
	}

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	j.session.Stdout = &stdoutBuf
	j.session.Stderr = &stderrBuf
	j.stdin, err = j.session.StdinPipe()
	if err != nil {
		j.err = err
		j.done <- false
		return
	}
	//j.session.Stdin = &j.stdinBuf

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	j.session.RequestPty("xterm", 40, 80, modes)

	err = j.session.Run(j.command)
	if err != nil {
		j.err = err
		j.done <- false
		return
	}
	j.result = stdoutBuf.Bytes()
	j.done <- true

	j.session.Close()
}

func (j *BackGroundjob) Result() ([]byte, error) {
	<- j.done
	if (j.err != nil) {
		return []byte{}, j.err
	} else {
		return j.result, nil
	}
}

func (j *BackGroundjob) Stop() ([]byte, error) {
	//j.session.Signal(ssh.SIGTERM)
	//j.stdinBuf.WriteString("\x03")
	j.stdin.Write([]byte{ 0x03 })

	return j.Result()
}

func (j *BackGroundjob) StopAndWriteFile(fileName string) {
	res, _ := j.Stop()
	err := ioutil.WriteFile(fileName, res, 777)
	if (err != nil) {
		log.Println("Error at writing output file:", err)
	}
}



func (c *SshClient) NewBackGroundjob(command string) *BackGroundjob {
	job := &BackGroundjob{
		command:	command,
		client:		c,
		done:		make(chan bool),
	}

	go job.Run()
	return job
}
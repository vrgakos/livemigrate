package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/vrgakos/livemigrate/ssh"
	. "github.com/vrgakos/livemigrate/measure"
	"sync"
	"time"
)

func migrateDockerCentralNfs(from *ssh.SshClient, to *ssh.SshClient, nfs *ssh.SshClient, containerId string) *Measure {
	m := NewMeasure()

	id, _ := uuid.NewRandom()
	snapshotId := id.String()[:8]
	fmt.Printf("Snapshot id: %s\n", snapshotId)

	// START TCPDUMP AND WAIT
	dumpNfs := NewTcpdump(nfs, "ens33", "D:\\nfs.pcap")
	dumpFrom := NewTcpdump(from, "ens33", "D:\\from.pcap")
	dumpTo := NewTcpdump(to, "ens33", "D:\\to.pcap")
	time.Sleep(time.Millisecond * 1000)

	// START TCP-CLIENT AND WAIT



	m.MilestoneStart()
	fmt.Printf("Container id: %s\n", containerId)
	ins := InspectDocker(from, containerId)
	m.AddMilestone("Inspect done")

	if ins.State.Running == false {
		fmt.Printf("Container not running!\n")
		return nil
	}

	// GET CONTAINER SETTINGS
	containerName := ins.Name[1:]
	var netName string
	var ipv4Address string
	for k, v := range ins.NetworkSettings.Networks {
		netName = k
		ipv4Address = v.IPAMConfig.IPv4Address
		break
	}
	fmt.Printf("Container name=%s, network name=%s, ip=%s\n", containerName, netName, ipv4Address)



	binds := ""
	for _, bind := range ins.HostConfig.Binds {
		binds += "-v " + bind + " "
	}

	// CREATE NEW CONTAINER
	createStr := fmt.Sprintf("docker rm -f %s ; docker create --name %s %s --net %s --ip %s %s", containerName, containerName, binds, netName, ipv4Address, ins.Config.Image)
	fmt.Printf("Create: %s\n", createStr)
	_, err := to.RunAndWait(createStr);
	/*newIdStr := strings.Trim(string(newId), "\r\n ")
	fmt.Printf("New id is: %s\n", newIdStr)*/
	if err != nil {
		fmt.Println(err)
		return nil
	}
	m.AddMilestone("Create done")



	// SNAPSHOT CREATE
	cpDir := "/mnt/docker/"
	fromStr := fmt.Sprintf("rm -rf " + cpDir +"%s/ ; docker checkpoint create --checkpoint-dir=%s %s %s", snapshotId, cpDir, containerId, snapshotId)
	fmt.Printf("Stop: %s\n", fromStr)
	if _, err := from.RunAndWait(fromStr); err != nil {
		fmt.Println(err)
		return nil
	}
	m.AddMilestone("Checkpoint done")



	// SNAPSHOT RESTORE
	startStr := fmt.Sprintf("docker start --checkpoint %s --checkpoint-dir %s %s", snapshotId, cpDir, containerName)
	fmt.Printf("Start: %s\n", startStr)
	if _, err := to.RunAndWait(startStr); err != nil {
		fmt.Println(err)
		return nil
	}
	m.AddMilestone("Restore done")



	// WAIT AND STOP TCP-CLIENT
	//time.Sleep(time.Second * 2)

	// WAIT AND STOP TCPDUMP
	time.Sleep(time.Millisecond * 1000)
	var wg sync.WaitGroup
	wg.Add(3)
	go (func() {
		dumpNfs.Stop()
		wg.Done()
	})()
	go (func() {
		dumpFrom.Stop()
		wg.Done()
	})()
	go (func() {
		dumpTo.Stop()
		wg.Done()
	})()
	wg.Wait()


	return m
}
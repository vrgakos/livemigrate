package main

import (
	"fmt"
	"os"
	"os/signal"
	"github.com/google/uuid"
	"encoding/json"
	"time"
	"strings"
	"github.com/vrgakos/livemigrate/ssh"
)


func dumpJSON(v interface{}, kn string) {
	iterMap := func(x map[string]interface{}, root string) {
		var knf string
		if root == "root" {
			knf = "%q:%q"
		} else {
			knf = "%s:%q"
		}
		for k, v := range x {
			dumpJSON(v, fmt.Sprintf(knf, root, k))
		}
	}

	iterSlice := func(x []interface{}, root string) {
		var knf string
		if root == "root" {
			knf = "%q:[%d]"
		} else {
			knf = "%s:[%d]"
		}
		for k, v := range x {
			dumpJSON(v, fmt.Sprintf(knf, root, k))
		}
	}

	switch vv := v.(type) {
	case string:
		fmt.Printf("%s => (string) %q\n", kn, vv)
	case bool:
		fmt.Printf("%s => (bool) %v\n", kn, vv)
	case float64:
		fmt.Printf("%s => (float64) %f\n", kn, vv)
	case map[string]interface{}:
		fmt.Printf("%s => (map[string]interface{}) ...\n", kn)
		iterMap(vv, kn)
	case []interface{}:
		fmt.Printf("%s => ([]interface{}) ...\n", kn)
		iterSlice(vv, kn)
	default:
		fmt.Printf("%s => (unknown?) ...\n", kn)
	}
}

func migrate(from *ssh.SshClient, to *ssh.SshClient, containerId string) {
	start := time.Now()

	fmt.Printf("Container id: %s\n", containerId)
	result, err := from.RunAndWait("docker inspect -f '{{ json . }}' " + containerId)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("-- Inspect: %d ms\n", time.Now().Sub(start).Nanoseconds() / 1000000)
	start = time.Now()

	var ins struct {
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
	if err := json.Unmarshal(result, &ins); err != nil {
		panic(err)
	}

	if ins.State.Running == false {
		fmt.Printf("Container not running!\n")
		return
	}

	name := ins.Name[1:]

	//fmt.Printf("%v", net)
	var netName string
	var ipv4Address string
	for k, v := range ins.NetworkSettings.Networks {
		netName = k
		ipv4Address = v.IPAMConfig.IPv4Address
		break
	}
	fmt.Printf("Network name=%s ip=%s\n", netName, ipv4Address)

	id, _ := uuid.NewRandom()
	idString := id.String()
	fmt.Printf("Snapshot id: %s\n", idString)

	binds := ""
	for _, bind := range ins.HostConfig.Binds {
		binds += "-v " + bind + " "
	}
	createStr := fmt.Sprintf("docker rm -f %s ; docker create --name %s %s --net %s --ip %s %s", name, name, binds,
		netName, ipv4Address, ins.Config.Image)
	fmt.Printf("Create: %s\n", createStr)
	newId, err := to.RunAndWait(createStr);
	newIdStr := strings.Trim(string(newId), "\r\n ")
	fmt.Printf("New id is: %s\n", newIdStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("-- Create new: %d ms\n", time.Now().Sub(start).Nanoseconds() / 1000000)
	start = time.Now()

	fromStr := fmt.Sprintf("rm -rf /mnt/docker/%s/ ; docker checkpoint create --checkpoint-dir=/mnt/docker/ %s %s", idString, containerId, idString)
	fmt.Printf("Stop: %s\n", fromStr)
	if _, err := from.RunAndWait(fromStr); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("-- Snapshot done: %d ms\n", time.Now().Sub(start).Nanoseconds() / 1000000)
	start = time.Now()


	startStr := fmt.Sprintf("docker start --checkpoint %s --checkpoint-dir /mnt/docker %s", idString, name)
	fmt.Printf("Start: %s\n", startStr)
	if _, err := to.RunAndWait(startStr); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("-- Start: %d ms\n", time.Now().Sub(start).Nanoseconds() / 1000000)
	start = time.Now()

	return

	updateArp := fmt.Sprintf("docker exec %s arping -c 1 -A -U %s ; echo OK", name, ipv4Address)
	fmt.Printf("Update ARP: %s\n", updateArp)
	if _, err := to.RunAndWait(updateArp); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("-- Update ARP: %d ms\n", time.Now().Sub(start).Nanoseconds() / 1000000)
	start = time.Now()
}



func main() {

	ubuntu1, _ := ssh.NewSshClient("ubuntu1","192.168.220.197", 22, "root", "d:\\KEY\\vrg-201603.id_rsa")
	ubuntu2, _ := ssh.NewSshClient("ubuntu2","192.168.220.198", 22, "root", "d:\\KEY\\vrg-201603.id_rsa")
	ubuntu3, _ := ssh.NewSshClient("ubuntu3","192.168.220.199", 22, "root", "d:\\KEY\\vrg-201603.id_rsa")
	nfs, _ := ssh.NewSshClient("nfs","192.168.220.200", 22, "root", "d:\\KEY\\vrg-201603.id_rsa")

	ids := make(map[string] *ssh.SshClient)
	ids["ubuntu1"] = ubuntu1
	ids["ubuntu2"] = ubuntu2
	ids["ubuntu3"] = ubuntu3
	ids["nfs"] = nfs

	go (func(){
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		os.Exit(0)
	})()


	for {
		var fromId, toId, contId string
		fmt.Scanf("%s %s %s\r\n", &fromId, &toId, &contId)

		from, ok1 := ids[fromId]
		to, ok2 := ids[toId]

		if ok1 && ok2 {
			migrateDockerDirectNfs(from, to, contId)
			//migrateDockerCentralNfs(from, to, nfs, contId)
			fmt.Printf("\nDone.\n")
		} else {
			fmt.Printf("%s\n", "Unknown targets")
		}
	}



}

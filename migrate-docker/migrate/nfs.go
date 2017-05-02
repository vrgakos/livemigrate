package migrate

import (
	. "github.com/vrgakos/livemigrate/node"
	"fmt"
	"github.com/google/uuid"
	"log"
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types"
)

func Nfs(from *Node, to *Node, containerId string, migrateOpts *DoOpts) *Measure {
	m := NewMeasure(migrateOpts)

	// NEW SNAPSHOT ID
	id, _ := uuid.NewRandom()
	checkPoint := id.String()[:8]
	fmt.Printf("Snapshot id: %s\n", checkPoint)

	m.AddStat(from, "/mnt/docker/stat/" + checkPoint + "/from.xlsx")
	m.AddStat(to, "/mnt/docker/stat/" + checkPoint + "/to.xlsx")
	m.Start()

	m.MilestoneStart()
	fmt.Printf("Provided container id: %s\n", containerId)

	dockerCliFrom, err := from.GetDockerClient()
	if err != nil {
		log.Println(err)
		return nil
	}

	dockerCliTo, err := to.GetDockerClient()
	if err != nil {
		log.Println(err)
		return nil
	}


	// INSPECT
	ins, err := dockerCliFrom.ContainerInspect(context.Background(), containerId)
	if err != nil {
		log.Println(err)
		return nil
	}
	containerName := ins.Name[1:]
	m.AddMilestone("Inspect done")

	if ins.State.Running == false {
		fmt.Printf("Container is not running!\n")
		return nil
	}


	// CREATE CONTAINER DUMMY
	dockerCliTo.ContainerRemove(context.Background(),
		containerName,
		types.ContainerRemoveOptions{
			Force: true,
		},
	)
	containerCreatedBody, err := dockerCliTo.ContainerCreate(context.Background(),
		&container.Config{
			Image: ins.Config.Image,
		},
		&container.HostConfig{
			Binds: ins.HostConfig.Binds,
			DNS: ins.HostConfig.DNS,
			DNSOptions: ins.HostConfig.DNSOptions,
			DNSSearch: ins.HostConfig.DNSSearch,
		},
		&network.NetworkingConfig{
			EndpointsConfig: ins.NetworkSettings.Networks,
		},
		containerName,
	)
	if err != nil {
		log.Println("ContainerCreate error:", err)
		return nil
	}
	m.AddMilestone("ContainerCreate done")


	cpDir := "/mnt/" + to.Alias + "/"
	parentCpDir := ""
	statCalc := NewStatCalc(migrateOpts)

	for statCalc.Resume() {
		// CHECKPOINT CREATE
		cpId := fmt.Sprintf("%s.%d", checkPoint, statCalc.GetIters())
		stats, err := dockerCliFrom.CheckpointCreate(context.Background(),
			ins.ID,
			types.CheckpointCreateOptions{
				CheckpointDir: cpDir,
				CheckpointID: cpId,
				ParentPath: parentCpDir,
				PreDump: true,
				Exit: false,
			},
		)
		if err != nil {
			log.Println("CheckpointCreate (pre-dump) error:", err)
			return nil
		}
		statCalc.Add(&stats)
		m.AddMilestone("CheckpointCreate (pre-dump) done")
		parentCpDir = "../" + cpId
	}


	// FINAL CHECKPOINT CREATE
	_, err = dockerCliFrom.CheckpointCreate(context.Background(),
		ins.ID,
		types.CheckpointCreateOptions{
			CheckpointDir: cpDir,
			CheckpointID: checkPoint,
			ParentPath: parentCpDir,
			PreDump: false,
			Exit: true,
		},
	)
	if err != nil {
		log.Println("CheckpointCreate error:", err)
		return nil
	}
	m.AddMilestone("CheckpointCreate done")


	// CHECKPOINT RESTORE
	err = dockerCliTo.ContainerStart(context.Background(),
		containerCreatedBody.ID,
		types.ContainerStartOptions{
			CheckpointDir: cpDir,
			CheckpointID: checkPoint,
		},
	)
	if err != nil {
		log.Println("ContainerStart error:", err)
		return nil
	}
	m.AddMilestone("ContainerStart done")


	m.Stop()
	return m
}


package node

import (
	"io/ioutil"
	"encoding/json"
	"log"
)

type NodeStore struct {
	file		string
	Nodes		[]*Node
}

func NewNodeStore(file string) *NodeStore {
	store := &NodeStore{
		file:	file,
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)

		store.Nodes = make([]*Node, 0)
		err := store.Save()
		if err != nil {
			// cry
			log.Println(err)
			return nil
		}

		return store
	}

	if err := json.Unmarshal(b, &store.Nodes); err != nil {
		log.Println(err)

		store.Nodes = make([]*Node, 0)
		err := store.Save()
		if err != nil {
			// cry
			log.Println(err)
			return nil
		}

		return store
	}
	return store
}


func (ns *NodeStore) Save() error {
	b, err := json.Marshal(ns.Nodes)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(ns.file, b, 0777)
	if err != nil {
		return err
	}
	return nil
}


func (ns *NodeStore) NewNode(alias string, host string, sshUser string, sshKey string) *Node {
	node := &Node{
		Alias:		alias,
		Host:		host,
		SshUser:	sshUser,
		SshKey:		sshKey,

		SshPort:	22,
		DockerApiPort:	2376,
	}

	ns.Nodes = append(ns.Nodes, node)
	err := ns.Save()
	if err != nil {
		// cry
		log.Println(err)
		return nil
	}

	return node
}

func (ns *NodeStore) GetNode(alias string) *Node {
	for _, node := range ns.Nodes {
		if node.Alias == alias {
			return node
		}
	}

	return nil
}
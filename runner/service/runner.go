package service

import (
	"fmt"
	"sync"

	"github.com/kavos113/seseragi/model"
)

type NodeRunner interface {
	Run(node model.Node) error
}

type NodeInfo struct {
	node      model.Node
	dependsOn []*NodeInfo
	doneCh    chan struct{} // 自分のタスクが完了したことを通知
	runneer   NodeRunner
	err       error
}

type WorkflowRunner struct {
}

func (wr *WorkflowRunner) RunWorkflow(workflow model.Workflow) error {
	nodes := make(map[string]*NodeInfo)
	for _, node := range workflow.Nodes {
		nodes[node.Name] = &NodeInfo{
			node:   node,
			doneCh: make(chan struct{}),
			err:    nil,
		}
	}

	for _, nodeInfo := range nodes {
		for _, dep := range nodeInfo.node.Dependencies {
			depNodeInfo, ok := nodes[dep]
			if !ok {
				return fmt.Errorf("dependency task %s not found", dep)
			}
			nodeInfo.dependsOn = append(nodeInfo.dependsOn, depNodeInfo)
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(nodes))

	for _, nodeInfo := range nodes {
		go func(n *NodeInfo) {
			defer close(n.doneCh)
			defer wg.Done()

			for _, dep := range n.dependsOn {
				<-dep.doneCh

				if dep.err != nil {
					n.err = fmt.Errorf("dependency task %s failed: %w", dep.node.TaskID, dep.err)
					return
				}
			}
			fmt.Printf("all dependencies of task %s are completed, running task...\n", n.node.TaskID)
			n.err = n.runneer.Run(n.node)
		}(nodeInfo)
	}

	wg.Wait()
	for _, nodeInfo := range nodes {
		if nodeInfo.err != nil {
			return fmt.Errorf("task %s failed: %w", nodeInfo.node.TaskID, nodeInfo.err)
		}
	}
	return nil
}

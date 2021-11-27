package install

import (
	"fmt"
	"sync"
)

//BuildExec is
func BuildExec(execNodes []string) {
	if len(execNodes) > 0 {
		execNodesFunc(execNodes)
	}
}

func execNodesFunc(execNodes []string) {
	nodes := execNodes
	i := &SdosInstaller{
		Nodes: nodes,
	}
	i.ExecNodes()
	NodeIPs = append(NodeIPs, execNodes...)
}

func (s *SdosInstaller) ExecNodes() {
	var wg sync.WaitGroup
	for _, node := range s.Nodes {
		wg.Add(1)
		go func(node string) {
			defer wg.Done()
			name := StringMd5(ExecConfig.Command)
			_ = SSHConfig.SaveFileX(node, fmt.Sprintf("/tmp/.hook_%s", name), BaseHookUtils(node))
			_ = SSHConfig.CmdAsync(node, fmt.Sprintf("/tmp/.hook_%s %s", name, Base64Decode(ExecConfig.Param)))
			execDone(node)
		}(node)
	}
	wg.Wait()
	ResultInstall.Range(execWalk)
}

func execDone(node string) {
	ResultInstall.Store(node, "success")
}

func execWalk(_ interface{}, _ interface{}) bool {
	return true
}
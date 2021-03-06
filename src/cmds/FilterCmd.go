package cmds

import "fmt"
import "subsystems"
import "visitors"

type FilterCmd struct {
	enforcer   *subsystems.Enforcer
	agent_ip   string
	agent_port string
}

func NewFilterCmd(agent_ip string, agent_port string) FilterCmd {
	return FilterCmd{
		enforcer:   visitors.EnforcerVisitor.EnforcerInstance,
		agent_ip:   agent_ip,
		agent_port: agent_port,
	}
}

// Call agent.enforcer to begin NFQ
func (filterCmd FilterCmd) ExecCmd() {
	fmt.Println("Executing Filtering Command...")
	filterCmd.enforcer.Start()
}

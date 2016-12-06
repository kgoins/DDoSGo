package cmds

import "fmt"
import "subsystems"
import "visitors"

type FilterCmd struct {
	enforcer           subsystems.Enforcer
	agent_ip           string
	startFilterPackets bool
}

func NewFilterCmd(agent_ip string, startFilteringPackets bool) FilterCmd {
	return FilterCmd{
		enforcer:           *visitors.EnforcerVisitor.EnforcerInstance,
		agent_ip:           agent_ip,
		startFilterPackets: startFilteringPackets,
	}
}

// Call agent.enforcer to begin NFQ
func (filterCmd FilterCmd) ExecCmd() {
	fmt.Println("Executing filtering command")
	filterCmd.enforcer.Start()
}

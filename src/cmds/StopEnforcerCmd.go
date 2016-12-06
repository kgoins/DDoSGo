package cmds

import "fmt"
import "subsystems"
import "visitors"

type StopEnforcerCmd struct {
	enforcer   *subsystems.Enforcer
	agent_ip   string
	agent_port string
}

func NewStopEnforcerCmd(agent_ip string, agent_port string) StopEnforcerCmd {
	return StopEnforcerCmd{
		enforcer:   visitors.EnforcerVisitor.EnforcerInstance,
		agent_ip:   agent_ip,
		agent_port: agent_port,
	}
}

// Call agent.enforcer to begin NFQ
func (stopEnforcerCmd StopEnforcerCmd) ExecCmd() {
	fmt.Println("Executing Stop Filtering Command...")
	stopEnforcerCmd.enforcer.Close()
}

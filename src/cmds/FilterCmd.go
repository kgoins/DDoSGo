package cmds

import "subsystems"
import "visitors"

type FilterCmd struct {
	enforcer           subsystems.Enforcer
	startFilterPackets bool
}

func NewFilterCmd(startFilteringPackets bool) FilterCmd {
	return FilterCmd{
		enforcer:           *visitors.EnforcerVisitor.EnforcerInstance,
		startFilterPackets: startFilteringPackets,
	}
}

func (filterCmd FilterCmd) ExecCmd() {
}

// Call agent.enforcer

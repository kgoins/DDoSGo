package visitors

import "subsystems"

var (
	alertVisitor *AlertVisitor
	agentRegVisitor *AgentRegistryVisitor
)

type AlertVisitor struct {
	alertingSys *subsystems.AlertSystem
}

type AgentRegistryVisitor struct {
	agentReg *subsystems.AgentRegistry
}

func NewAlertSystemVisitor(alertingSys *subsystems.AlertSystem) *AlertVisitor {
	return &AlertVisitor{alertingSys: alertingSys}
}

func (alertVisitor *AlertVisitor) HeldInstance() *subsystems.AlertSystem {
	return alertVisitor.alertingSys
}

func NewAgentRegistryVisitor(agentReg *subsystems.AgentRegistry) *AgentRegistryVisitor {
	return &AgentRegistryVisitor{agentReg: agentReg}
}



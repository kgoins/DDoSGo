package visitors

import "subsystems"

var (
	AlertingVisitor *AlertVisitor
	AgentRegVisitor *AgentRegistryVisitor
)

type AlertVisitor struct {
	AlertingSys *subsystems.AlertSystem
}

type AgentRegistryVisitor struct {
	AgentReg *subsystems.AgentRegistry
}

func NewAlertSystemVisitor(alertingSys *subsystems.AlertSystem) *AlertVisitor {
	return &AlertVisitor{AlertingSys: alertingSys}
}

func (alertVisitor *AlertVisitor) HeldInstance() *subsystems.AlertSystem {
	return alertVisitor.AlertingSys
}

func (Reg *AgentRegistryVisitor) HeldInstance() *subsystems.AgentRegistry {
	return Reg.AgentReg
}

func NewAgentRegistryVisitor(agentReg *subsystems.AgentRegistry) *AgentRegistryVisitor {
	return &AgentRegistryVisitor{AgentReg: agentReg}
}

func SetupVisitors(alertingSys *subsystems.AlertSystem, agentReg *subsystems.AgentRegistry) {
	AlertingVisitor = NewAlertSystemVisitor(alertingSys)
	AgentRegVisitor = NewAgentRegistryVisitor(agentReg)
}

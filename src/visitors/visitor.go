package visitors

import "subsystems"
import "fmt"

var (
	AlertingVisitor *AlertVisitor
	AgentRegVisitor *AgentRegistryVisitor
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

func (Reg *AgentRegistryVisitor) HeldInstance() *subsystems.AgentRegistry {
	return Reg.agentReg
}

func NewAgentRegistryVisitor(agentReg *subsystems.AgentRegistry) *AgentRegistryVisitor {
	return &AgentRegistryVisitor{agentReg: agentReg}
}

func SetupVisitors(alertingSys *subsystems.AlertSystem, agentReg *subsystems.AgentRegistry) {
	AlertingVisitor = NewAlertSystemVisitor(alertingSys)
	AgentRegVisitor = NewAgentRegistryVisitor(agentReg)
}

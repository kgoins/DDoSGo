package visitors

import "subsystems"

var (
	AlertingVisitor *AlertVisitor
	AgentRegVisitor *AgentRegistryVisitor
	EnforcerVisitor *EnforcementVisitor
)

type AlertVisitor struct {
	AlertingSys *subsystems.AlertSystem
}

type AgentRegistryVisitor struct {
	AgentReg *subsystems.AgentRegistry
}

type EnforcementVisitor struct {
	EnforcerInstance *subsystems.Enforcer
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

func NewEnforcementVisitor(enforcer *subsystems.Enforcer) *EnforcementVisitor {
	return &EnforcementVisitor{EnforcerInstance: enforcer}
}

func (EnforcerVisitor EnforcementVisitor) HeldInstance() *subsystems.Enforcer {
	return EnforcerVisitor.EnforcerInstance
}

func SetupHandlerVisitors(alertingSys *subsystems.AlertSystem, agentReg *subsystems.AgentRegistry) {
	AlertingVisitor = NewAlertSystemVisitor(alertingSys)
	AgentRegVisitor = NewAgentRegistryVisitor(agentReg)
}

func SetupAgentVisitors(enforcer *subsystems.Enforcer) {
	EnforcerVisitor = NewEnforcementVisitor(enforcer)
}

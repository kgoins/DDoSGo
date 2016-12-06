package cmds

import (
	// "fmt"
	"subsystems"
	"visitors"
)

type ProcDataStreamCmd struct {
	alertingSystem *subsystems.AlertSystem
	agentReg       *subsystems.AgentRegistry
	agent_ip       string
	cpu            int
	mem            int
	bytesRecvd     int
	bytesSent      int
}

func NewProcDataStreamCmd(agent_ip string, cpu int, mem int, bytesRecvd int, bytesSent int) ProcDataStreamCmd {
	return ProcDataStreamCmd{
		alertingSystem: visitors.AlertingVisitor.AlertingSys,
		agentReg:       visitors.AgentRegVisitor.AgentReg,
		agent_ip:       agent_ip,
		cpu:            cpu,
		mem:            mem,
		bytesRecvd:     bytesRecvd,
		bytesSent:      bytesSent,
	}
}

func (procStreamCmd ProcDataStreamCmd) ExecCmd() {

	// Update record status to show it has been in contact recently
	procStreamCmd.agentReg.UpdateRecordStatus(procStreamCmd.agent_ip)

	// Check if record is currently filtering
	isFiltering := procStreamCmd.agentReg.IsAgentFiltering(procStreamCmd.agent_ip)

	// Page alert system to process data for anomalies
	procStreamCmd.alertingSystem.ProcessDataStream(procStreamCmd.agent_ip, procStreamCmd.cpu,
		procStreamCmd.mem, procStreamCmd.bytesRecvd, procStreamCmd.bytesSent, isFiltering)

}

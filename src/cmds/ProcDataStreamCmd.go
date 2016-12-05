package cmds

import (
	// "fmt"
	"subsystems"
	"visitors"
)

type ProcDataStreamCmd struct {
	alertingSystem subsystems.AlertSystem
	agentReg		subsystems.AgentRegistry
	agent_ip	string
	cpu        int
	mem        int
	bytesRecvd int
	bytesSent  int
}

func NewProcDataStreamCmd(agent_ip string, cpu int, mem int, bytesRecvd int, bytesSent int) ProcDataStreamCmd {
	return ProcDataStreamCmd{
		alertingSystem: *visitors.AlertingVisitor.AlertingSys,
		agentReg:		*visitors.AgentRegVisitor.AgentReg,
		agent_ip:	agent_ip,
		cpu: 		cpu,
		mem:        mem,
		bytesRecvd: bytesRecvd,
		bytesSent:  bytesSent,
	}
}

func (procStreamCmd ProcDataStreamCmd) ExecCmd() {
	// fmt.Println("From procStream cmd: ", procStreamCmd)
	procStreamCmd.agentReg.UpdateRecordStatus(procStreamCmd.agent_ip)
	procStreamCmd.alertingSystem.ProcessDataStream(procStreamCmd.cpu, procStreamCmd.mem, procStreamCmd.bytesRecvd, procStreamCmd.bytesSent)
}

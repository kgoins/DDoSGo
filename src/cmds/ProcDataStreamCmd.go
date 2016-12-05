package cmds

import (
	"fmt"
	"subsystems"
	"visitors"
)

type ProcDataStreamCmd struct {
	alertingSystem subsystems.AlertSystem
	agentReg		subsystems.AgentRegistry
	cpu        int
	mem        int
	bytesRecvd int
	bytesSent  int
}

func NewProcDataStreamCmd(cpu int, mem int, bytesRecvd int, bytesSent int) ProcDataStreamCmd {
	return ProcDataStreamCmd{
		alertingSystem: *visitors.AlertingVisitor.AlertingSys,
		agentReg:		*visitors.AgentRegVisitor.AgentReg,
		cpu: 		cpu,
		mem:        mem,
		bytesRecvd: bytesRecvd,
		bytesSent:  bytesSent,
	}
}

func (procStreamCmd ProcDataStreamCmd) ExecCmd() {
	fmt.Println("From procStream cmd: ", procStreamCmd)
	procStreamCmd.alertingSystem.ProcessDataStream(procStreamCmd.cpu, procStreamCmd.mem, procStreamCmd.bytesRecvd, procStreamCmd.bytesSent)
}

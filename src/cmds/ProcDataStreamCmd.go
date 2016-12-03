package cmds

import "fmt"

type ProcDataStream struct {
	cpu        int
	mem        int
	bytesRecvd int
	bytesSent  int
}

func NewProcDataStream(cpu int, mem int, bytesRecvd int, bytesSent int) ProcDataStream {
	return ProcDataStream{cpu: cpu,
		mem:        mem,
		bytesRecvd: bytesRecvd,
		bytesSent:  bytesSent}
}

func (procStream ProcDataStream) ExecCmd() {
	fmt.Println("From procStream cmd: ", procStream)
}

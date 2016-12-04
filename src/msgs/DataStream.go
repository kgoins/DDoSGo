package msgs

import "fmt"
import "cmds"

type DataStream struct {
	Cpu        int
	Mem        int
	BytesRecvd int
	BytesSent  int
}

func NewDataStream(cpu int, mem int, bytesRecvd int, bytesSent int) DataStream {
	return DataStream{
		Cpu:        cpu,
		Mem:        mem,
		BytesRecvd: bytesRecvd,
		BytesSent:  bytesSent}
}

func (stream DataStream) GetType() string {
	return "DataStream"
}

func (stream DataStream) String() string {
	data := fmt.Sprintf("cpu: %d, mem: %d, ntwk: %d %d", stream.Cpu, stream.Mem,
		stream.BytesRecvd, stream.BytesSent)
	return fmt.Sprintf("type: %s, payload: %s", stream.GetType(), data)
}

func (stream DataStream) BuildCmd() cmds.Cmd {
	return cmds.NewProcDataStreamCmd(stream.Cpu, stream.Mem, stream.BytesRecvd, stream.BytesSent)
}

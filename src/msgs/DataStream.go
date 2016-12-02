package msgs

import "fmt"

type DataStream struct {
	Cpu        int
	Mem        int
	BytesRecvd int
	BytesSent  int
}

func NewDataStream(cpu, mem, bytesRecvd, bytesSent int) DataStream {
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

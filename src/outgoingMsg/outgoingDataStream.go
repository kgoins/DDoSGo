package outgoingMsg

import "fmt"

type OutgoingDataStream struct {
	Cpu        int
	Mem        int
	BytesRecvd int
	BytesSent  int
}

func NewOutgoingDataStream(cpu int, mem int, bytesRecvd int, bytesSent int) OutgoingDataStream {
	return OutgoingDataStream{
		Cpu:        cpu,
		Mem:        mem,
		BytesRecvd: bytesRecvd,
		BytesSent:  bytesSent}
}

func (stream OutgoingDataStream) GetType() string {
	return "DataStream"
}

func (stream OutgoingDataStream) String() string {
	data := fmt.Sprintf("cpu: %d, mem: %d, ntwk: %d %d", stream.Cpu, stream.Mem,
		stream.BytesRecvd, stream.BytesSent)
	return fmt.Sprintf("type: %s, payload: %s", stream.GetType(), data)
}

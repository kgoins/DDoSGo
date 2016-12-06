package outgoingMsg

import "fmt"

type OutgoingDataStream struct {
	Agent_ip     string
	Agent_port   string
	Handler_ip   string
	Handler_port string
	Cpu          int
	Mem          int
	BytesRecvd   int
	BytesSent    int
}

func NewOutgoingDataStream(agent_ip string, aPort string, handler_ip string, handler_port string, cpu int, mem int, bytesRecvd int, bytesSent int) OutgoingDataStream {
	return OutgoingDataStream{
		Agent_ip:     agent_ip,
		Agent_port:   aPort,
		Handler_ip:   handler_ip,
		Handler_port: handler_port,
		Cpu:          cpu,
		Mem:          mem,
		BytesRecvd:   bytesRecvd,
		BytesSent:    bytesSent}
}

func (stream OutgoingDataStream) GetType() string {
	return "DataStream"
}

func (stream OutgoingDataStream) String() string {
	data := fmt.Sprintf("cpu: %d, mem: %d, ntwk: %d %d", stream.Cpu, stream.Mem,
		stream.BytesRecvd, stream.BytesSent)
	return fmt.Sprintf("type: %s, payload: %s", stream.GetType(), data)
}

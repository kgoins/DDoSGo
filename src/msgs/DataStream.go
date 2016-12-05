package msgs

import "fmt"
import "cmds"

type DataStream struct {
	Agent_ip 		string
	Handler_ip 		string
	Handler_port 	string 
	Cpu        int
	Mem        int
	BytesRecvd int
	BytesSent  int
}

func NewDataStream(agent_ip string, handler_ip string, handler_port string, cpu int, mem int, bytesRecvd int, bytesSent int) DataStream {
	return DataStream{
		Agent_ip: 		agent_ip,
		Handler_ip: 	handler_ip,
		Handler_port: 	handler_port,
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
	return cmds.NewProcDataStreamCmd(stream.Agent_ip, stream.Cpu, stream.Mem, stream.BytesRecvd, stream.BytesSent)
}

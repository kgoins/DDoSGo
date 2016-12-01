package msgs

import "fmt"

type DataStream struct {
	Cpu  int
	Mem  int
	Ntwk int
}

func NewDataStream(cpu, mem, ntwk int) DataStream {
	return DataStream{
		Cpu:  cpu,
		Mem:  mem,
		Ntwk: ntwk}
}

func (stream DataStream) GetType() string {
	return "DataStream"
}

func (stream DataStream) String() string {
	data := fmt.Sprintf("cpu: %d, mem: %d, ntwk: %d", stream.Cpu, stream.Mem, stream.Ntwk)
	return fmt.Sprintf("type: %s, payload: %s", stream.GetType(), data)
}

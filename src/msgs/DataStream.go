package msgs

import "fmt"

type DataStream struct {
	Cpu  int
	Ram  int
	Ntwk int
}

// Construct new data stream -- take in # workers & function they should call
func NewDataStream(cpu, ram, ntwk int) DataStream {
	return DataStream{
		Cpu:  cpu,
		Ram:  ram,
		Ntwk: ntwk}
}

func (stream DataStream) GetType() string {
	return "DataStream"
}

func (stream DataStream) String() string {
	data := fmt.Sprintf("cpu: %d, ram: %d, ntwk: %d", stream.Cpu, stream.Ram, stream.Ntwk)
	return fmt.Sprintf("type: %s, payload: %s", stream.GetType(), data)
}

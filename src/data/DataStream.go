package data

import "fmt"

type DataStream struct {
	cpu  int
	ram  int
	ntwk int
}

// Construct new data stream -- take in # workers & function they should call
func NewDataStream(cpu, ram, ntwk int) DataStream {
	return DataStream{
		cpu:  cpu,
		ram:  ram,
		ntwk: ntwk}
}

func (stream DataStream) GetType() string {
	return "DataStream"
}

func (stream DataStream) String() string {
	data := fmt.Sprintf("cpu: %s, ram: %s, ntwk: %s", stream.cpu, stream.ram, stream.ntwk)
	return fmt.Sprintf("type: %s, payload: %s", GetType(), data)
}

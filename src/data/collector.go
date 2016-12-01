package data

import "msgs"

type DataCollector struct {
	msgChan chan Msg
	intVal  int
}

func NewDataCollector(msgChan chan Msg, intval int) DataCollector {
	return DataCollector{msgChan: msgChan, intVal: intval}
}

func (collector DataCollector) Start() {
	go func() {
		for {
			cpu := cpuUtil()
			ram := ramUtil()
			ntwk := ntwkUtil()

			dataStream := NewDataStream(cpu, ram, ntwk)
			collector.msgChan <- dataStream

			time.Sleep(collector.intVal * time.Millisecond)
		}
	}()
}

func cpuUtil() int {
	return 12
}

func ramUtil() int {
	return 60
}

func ntwkUtil() int {
	return 256
}

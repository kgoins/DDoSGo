package data

import "msgs"
import "time"

type DataCollector struct {
	msgChan chan msgs.Msg
	intVal  int
}

func NewDataCollector(msgChan chan msgs.Msg, intval int) DataCollector {
	return DataCollector{msgChan: msgChan, intVal: intval}
}

func (collector DataCollector) Start() {
	go func() {
		for {
			cpu := cpuUtil()
			ram := ramUtil()
			ntwk := ntwkUtil()

			dataStream := msgs.NewDataStream(cpu, ram, ntwk)
			collector.msgChan <- dataStream

			time.Sleep(time.Second * time.Duration(collector.intVal))
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

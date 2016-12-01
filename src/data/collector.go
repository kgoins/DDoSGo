package data

import "msgs"
import "time"

type DataCollector struct {
	msgChan  chan msgs.Msg
	intVal   int
	shutdown chan bool
}

func NewDataCollector(msgChan chan msgs.Msg, intval int) DataCollector {
	shutdown := make(chan bool)
	return DataCollector{
		msgChan:  msgChan,
		intVal:   intval,
		shutdown: shutdown}
}

func (collector DataCollector) Start() {
	go func() {
		for {
			select {
			case <-collector.shutdown:
				return

			default:
				cpu := cpuUtil()
				ram := ramUtil()
				ntwk := ntwkUtil()

				dataStream := msgs.NewDataStream(cpu, ram, ntwk)
				collector.msgChan <- dataStream

				time.Sleep(time.Second * time.Duration(collector.intVal))
			}
		}
	}()
}

func (collector DataCollector) Close() {
	collector.shutdown <- true
	close(collector.shutdown)
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

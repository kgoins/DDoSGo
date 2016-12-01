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
				data := collectData()

				dataStream := buildDataStream(data)
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

// Data Collection
type Data struct {
	cpu      int
	ram      int
	ntwkUtil int
}

func collectData() Data {
	cpu := cpuUtil()
	ram := ramUtil()
	ntwkUtil := ntwkUtil()

	return Data{cpu: cpu, ram: ram, ntwkUtil: ntwkUtil}
}

func buildDataStream(data Data) msgs.DataStream {
	return msgs.NewDataStream(data.cpu, data.ram, data.ntwkUtil)
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

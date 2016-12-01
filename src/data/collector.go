package main

// package data

import "time"
import "os"

import "bufio"
import "fmt"
import "strings"
import "strconv"

import "msgs"

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
	mem      int
	ntwkUtil int
}

func collectData() Data {
	cpu := cpuUtil()
	mem := memUtil()
	ntwkUtil := ntwkUtil()

	return Data{cpu: cpu, mem: mem, ntwkUtil: ntwkUtil}
}

func buildDataStream(data Data) msgs.DataStream {
	return msgs.NewDataStream(data.cpu, data.mem, data.ntwkUtil)
}

func cpuUtil() int {
	return 12
}

func memUtil() int {
	fileHandle, _ := os.Open("/proc/meminfo")
	defer fileHandle.Close()

	fileScanner := bufio.NewScanner(fileHandle)

	// memTotal = memStats[0]
	// memFree = memStats[1]
	var memStats [2]int

	for i := 0; i < 2; i++ {
		fileScanner.Scan()

		line := strings.Fields(fileScanner.Text())
		memStats[i], _ = strconv.Atoi(line[1])
	}

	memUtil := float64(memStats[1]) / float64(memStats[0])
	fmt.Println(memUtil)

	return int(memUtil * 100)
}

func ntwkUtil() int {
	return 256
}

// ** Test of Collector ** //
func main() {
	fmt.Println(memUtil())
}

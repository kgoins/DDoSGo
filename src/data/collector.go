package main

// package data

import "time"
import "os"

import "bufio"
import "io/ioutil"
import "fmt"
import "strings"
import "strconv"

import "msgs"

type DataCollector struct {
	msgChan       chan msgs.Msg
	collectIntval int
	sendIntval    int
	shutdown      chan bool
}

func NewDataCollector(msgChan chan msgs.Msg, sendIntval int, collectIntval int) DataCollector {
	shutdown := make(chan bool)
	return DataCollector{
		msgChan:       msgChan,
		sendIntval:    sendIntval,
		collectIntval: collectIntval,
		shutdown:      shutdown}
}

func (collector DataCollector) Start() {
	go func() {
		for {
			select {
			case <-collector.shutdown:
				return

			default:
				data := collectData(collector.collectIntval)

				dataStream := buildDataStream(data)
				collector.msgChan <- dataStream

				time.Sleep(time.Second * time.Duration(collector.sendIntval))
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

func collectData(intVal int) Data {
	cpu := cpuUtil(intVal)
	mem := memUtil()
	bytesRecv, _ := ntwkUtil(intVal)

	return Data{cpu: cpu, mem: mem, ntwkUtil: bytesRecv}
}

func buildDataStream(data Data) msgs.DataStream {
	return msgs.NewDataStream(data.cpu, data.mem, data.ntwkUtil)
}

// Calculates the current percentage CPU utilization
func cpuUtil(intVal int) int {
	prevIdle, prevTotal := currCpuStats()
	time.Sleep(time.Second * time.Duration(intVal))
	currIdle, currTotal := currCpuStats()

	diffIdle := currIdle - prevIdle
	diffTotal := currTotal - prevTotal

	return ((1000 * (diffTotal - diffIdle) / diffTotal) + 5) / 10

}

func currCpuStats() (int, int) {
	fileHandle, _ := os.Open("/proc/stat")
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)

	fileScanner.Scan()
	cpuRow := strings.Fields(fileScanner.Text())
	fmt.Println(cpuRow)
	cpuRow = cpuRow[1:]

	cpuIdle := 0
	cpuTotal := 0
	for i := 0; i < len(cpuRow); i++ {
		cpuVal, _ := strconv.Atoi(cpuRow[i])
		cpuTotal += cpuVal

		if i == 3 {
			cpuIdle = cpuVal
		}
	}

	return cpuIdle, cpuTotal
}

// Calculates the percentage of memory currently in use
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
	return int(memUtil * 100)
}

// Calculates byte throughput (sent and received) across all interfaces
func ntwkUtil(intVal int) (int, int) {
	prevBytesRecv, prevBytesSent := getNtwkThroughput()
	time.Sleep(time.Second * time.Duration(intVal))
	currBytesRecv, currBytesSent := getNtwkThroughput()

	bytesRecv := (currBytesRecv - prevBytesRecv) / intVal
	bytesSent := (currBytesSent - prevBytesSent) / intVal

	return bytesRecv, bytesSent
}

func getNtwkThroughput() (int, int) {
	netDevFile, _ := ioutil.ReadFile("/proc/net/dev")
	netDevLines := strings.Split(string(netDevFile), "\n")

	// remove header rows
	netDevLines = netDevLines[2:]

	var ifaceValues [][]int

	for _, line := range netDevLines {
		row := strings.Fields(line)
		byteVals := []int{0, 0}

		if len(row) == 0 || strings.Contains(row[0], "lo:") {
			continue
		} else {
			byteVals[0], _ = strconv.Atoi(row[1])
			byteVals[1], _ = strconv.Atoi(row[9])
		}

		ifaceValues = append(ifaceValues, byteVals)
	}

	bytesRecv := sumCols(ifaceValues, 0)
	bytesSent := sumCols(ifaceValues, 1)

	return bytesRecv, bytesSent
}

// Unility functions
func sumCols(array [][]int, colNum int) int {
	sum := 0
	for i := 0; i < len(array); i++ {
		sum += array[i][colNum]
	}

	return sum
}

func getMax(num1, num2 int) int {
	if num1 > num2 {
		return num1
	} else {
		return num2
	}
}

// ** Test of Collector ** //
func main() {
	cpuIdle, cpuTotal := currCpuStats()
	bytesRecv, bytesSent := getNtwkThroughput()

	fmt.Println("cpu times: ", cpuIdle, cpuTotal)
	fmt.Println("mem util: ", memUtil())
	fmt.Println("ntwk bytes: ", bytesRecv, bytesSent)
}

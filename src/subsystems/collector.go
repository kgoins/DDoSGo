package subsystems

// package data

import "time"
import "os"

import "bufio"
import "io/ioutil"
import "strings"
import "strconv"

import "outgoingMsg"
import "fmt"

type DataCollector struct {
	agent_ip		string
	handler_ip		string
	handler_port	string
	msgChan       chan outgoingMsg.OutgoingMsg
	collectIntval int
	sendIntval    int
	shutdown      chan bool
}

func NewDataCollector(agent_ip string, handler_ip string, handler_port string, msgChan chan outgoingMsg.OutgoingMsg, sendIntval int, collectIntval int) DataCollector {
	shutdown := make(chan bool)
	return DataCollector{
		agent_ip: agent_ip,
		handler_ip: handler_ip,
		handler_port: handler_port,
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
				dataStream := collectOutgoingData(collector.agent_ip, collector.handler_ip, collector.handler_port, collector.collectIntval)
				collector.msgChan <- dataStream

				time.Sleep(time.Second * time.Duration(collector.sendIntval))
			}
		}
	}()
}

func (collector DataCollector) Close() {
	collector.shutdown <- true
	fmt.Println("made it")
	close(collector.shutdown)
}

// Data Collection

/*func collectData(intVal int) msgs.DataStream {
	cpu := cpuUtil(intVal)
	mem := memUtil()
	bytesRecv, bytesSent := ntwkUtil(intVal)

	return msgs.NewDataStream(cpu, mem, bytesRecv, bytesSent)
}*/

func collectOutgoingData(agent_ip string, handler_ip string, handler_port string, intVal int) outgoingMsg.OutgoingDataStream {
	cpu := cpuUtil(intVal)
	mem := memUtil()
	bytesRecv, bytesSent := ntwkUtil(intVal)

	return outgoingMsg.NewOutgoingDataStream(agent_ip, handler_ip, handler_port, cpu, mem, bytesRecv, bytesSent)
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

	memUsed := memStats[0] - memStats[1]
	memUtil := float64(memUsed) / float64(memStats[0])
	return int(memUtil * 100)
}

// Calculates byte throughput (sent and received) across all interfaces
func ntwkUtil(intVal int) (int, int) {
	prevBytesRecv, prevBytesSent := getNtwkThroughput()
	time.Sleep(time.Second * time.Duration(intVal))
	currBytesRecv, currBytesSent := getNtwkThroughput()

	float_intval := float64(intVal)
	bytesRecv := (currBytesRecv - prevBytesRecv) / float_intval
	bytesSent := (currBytesSent - prevBytesSent) / float_intval

	return int(bytesRecv), int(bytesSent)
}

func getNtwkThroughput() (float64, float64) {
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

	return float64(bytesRecv), float64(bytesSent)
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
// func main() {
// 	recvBPS, sentBPS := ntwkUtil(1)

// 	fmt.Println("cpu utilization:", cpuUtil(1), "%")
// 	fmt.Println("mem used:", memUtil(), "%")
// 	fmt.Println("ntwk bytes:", recvBPS, sentBPS)
// }

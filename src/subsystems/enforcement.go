package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/subgraph/go-nfnetlink/nfqueue"
)

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func callIptables(action string, queueNum int) {
	var arg string

	switch action {
	case "start":
		arg = "-A"
		fmt.Println("starting Iptables")

	case "stop":
		arg = "-D"
		fmt.Println("stopping Iptables")
	}

	cmd := "iptables"
	args := []string{arg, "OUTPUT", "-p", "0", "-j", "NFQUEUE", "--queue-num", strconv.Itoa(queueNum), "--queue-bypass"}
	fmt.Println(args)

	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func startNFQ(queueNum int) {
	nfq := nfqueue.NewNFQueue(uint16(queueNum))
	fmt.Println("running queue")

	nfqPacketStream, err := nfq.Open()
	if err != nil {
		fmt.Printf("Error opening NFQueue: %v\n", err)
		os.Exit(1)
	}
	defer nfq.Close()

	for nfqPacket := range nfqPacketStream {
		filterPacket(nfqPacket)
	}
}

func filterPacket(nfqPacket *nfqueue.NFQPacket) {
	fmt.Println("Processing packet")

	if isPacketBad(nfqPacket.Packet) {
		nfqPacket.Accept()
	} else {
		nfqPacket.Drop()
	}
}

var (
	icmpLayer layers.ICMPv4
	ipLayer   layers.IPv4

	tcpLayer layers.TCP
	udpLayer layers.UDP

	dnsLayer layers.DNS
)

func getPacketLayers(packetData []byte) []gopacket.LayerType {

	parser := gopacket.NewDecodingLayerParser(
		layers.LayerTypeIPv4,
		&ipLayer,
		&icmpLayer,
		&tcpLayer,
		&udpLayer,
	)

	foundLayerTypes := []gopacket.LayerType{}
	err := parser.DecodeLayers(packetData, &foundLayerTypes)
	if err != nil {
		// fmt.Println("Trouble decoding layers: ", err)
	}

	return foundLayerTypes
}

func isPacketBad(packet gopacket.Packet) bool {
	packetLayers := getPacketLayers(packet.Data())

	for _, layer := range packetLayers {
		if layer == layers.LayerTypeIPv4 {
			fmt.Println("IPv4: ", ipLayer.SrcIP, "->", ipLayer.DstIP)
		}

		if layer == layers.LayerTypeICMPv4 {
			fmt.Println("dropping icmp")
			// return false
			return true
		}

		// if layerType == layers.LayerTypeTCP {
		// 	fmt.Println("TCP Port: ", tcpLayer.SrcPort, "->", tcpLayer.DstPort)
		// 	fmt.Println("TCP SYN:", tcpLayer.SYN, " | ACK:", tcpLayer.ACK)
		// }
	}

	return true
}

func main() {
	queueNum := 12

	callIptables("start", queueNum)

	startNFQ(queueNum)

	callIptables("stop", queueNum)
}

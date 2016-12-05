package subsystems

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

// static vars needed for DecodeLayers to run
var (
	icmpLayer layers.ICMPv4
	ipLayer   layers.IPv4

	tcpLayer layers.TCP
	udpLayer layers.UDP

	dnsLayer layers.DNS
)

type Enforcer struct {
	killsig  chan bool
	nfq      *nfqueue.NFQueue
	queueNum int
}

func NewEnforcer() *Enforcer {
	queueNum := 12 // From agent's conf??

	killsig := make(chan bool)
	nfq := nfqueue.NewNFQueue(uint16(queueNum))

	return &Enforcer{
		killsig:  killsig,
		queueNum: queueNum,
		nfq:      nfq}
}

func (enforcer *Enforcer) Close() {
	fmt.Println("closing enforcer")

	enforcer.Stop()
	close(enforcer.killsig)
}

func (enforcer *Enforcer) Start() {

	iptables("start", enforcer.queueNum)

	fmt.Println("made it")
	go enforcer.startNFQ()
}

func (enforcer *Enforcer) Stop() {
	enforcer.killsig <- true
	enforcer.nfq.Close()

	iptables("stop", enforcer.queueNum)
}

func (enforcer *Enforcer) startNFQ() {
	fmt.Println("starting packet filter")

	nfqPacketChan, err := enforcer.nfq.Open()
	if err != nil {
		fmt.Printf("Error opening NFQueue: %v\n", err)
		os.Exit(1)
	}

	for nfqPacket := range nfqPacketChan {
		if <-enforcer.killsig {
			return
		}

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
	}

	return true
}

// Utility functions
func iptables(action string, queueNum int) {
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

// func main() {
// 	queueNum := 12

// 	iptables("start", queueNum)

// 	startNFQ(queueNum)

// 	iptables("stop", queueNum)
// }

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
	"config"
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
	nfq      *nfqueue.NFQueue
	queueNum string
	running  bool
	offendingIPs []string
}

func NewEnforcer(queueNum string) *Enforcer {
	queueNum_int, _ := strconv.Atoi(queueNum)
	fmt.Println(queueNum_int)
	nfq := nfqueue.NewNFQueue(uint16(queueNum_int))

	ips,_ := config.ReadIpConf()


	return &Enforcer{
		queueNum: queueNum,
		nfq:      nfq,
		running:  false,
		offendingIPs: ips.IPs}
}

//Update the offending ips list
func (enforcer *Enforcer) UpdateOffendingIps(newIPs []string){
     enforcer.offendingIPs = newIPs

		 fmt.Println("Received ips: ", newIPs)

     for _, ip := range enforcer.offendingIPs{
     	 fmt.Println("Blocking new ip: ", ip)
     }
}

func (enforcer *Enforcer) Close() {

	if enforcer.running == true {

		fmt.Println("Closing Enforcer...")

		enforcer.Stop()

		enforcer.running = false
		fmt.Println("Enforcer Closed")

	} else {
		fmt.Println("Enforcer Not Currently Active, Ignoring Cmd")
	}

}

func (enforcer *Enforcer) Start() {
	if enforcer.running == false {
		enforcer.running = true
		iptables("start", enforcer.queueNum)
		go enforcer.startNFQ()
	} else {
		fmt.Println("Enforcer Already Running, Ignoring Filter Cmd")
	}

}

func (enforcer *Enforcer) Stop() {
	enforcer.nfq.Close()
	// fmt.Println("Closing nfq")

	iptables("stop", enforcer.queueNum)
}

func (enforcer *Enforcer) startNFQ() {
	fmt.Println("Starting Packet Filter...")

	nfqPacketChan, err := enforcer.nfq.Open()
	if err != nil {
		fmt.Printf("Error Opening NFQueue: %v\n", err)
		return
	}

	fmt.Println("Filtering Packets...")
	for nfqPacket := range nfqPacketChan {
		filterPacket(nfqPacket, enforcer.offendingIPs)
	}

	fmt.Println("exiting NFQ")
}

func filterPacket(nfqPacket *nfqueue.NFQPacket, offendingIPs []string) {
	// fmt.Println("Processing packet")

	if isPacketBad(nfqPacket.Packet, offendingIPs) {
		nfqPacket.Accept()
	} else {
		nfqPacket.Drop()
	}
}

func isPacketBad(packet gopacket.Packet, offendingIPs []string) bool {
	packetLayers := getPacketLayers(packet.Data())

	for _, layer := range packetLayers {

		if layer == layers.LayerTypeIPv4 {

			// dstIp := ipLayer.SrcIP.String()
			sourceIp := ipLayer.DstIP.String()

			// fmt.Println("ips: ", sourceIp, "  ", dstIp)

			if sourceIp == "127.0.0.1" || sourceIp == "192.168.56.103" {
				// fmt.Println("verdict: keeping")
				return true
			}

			for _, ip := range offendingIPs{

				if sourceIp == ip {
				   // fmt.Println("verdict: dropping")
				   return false
			        }
			}
		}

		if layer == layers.LayerTypeICMPv4 {
			// fmt.Println("verdict: dropping ping")
			return false
		}
	}

	// fmt.Println("unknown packet type")
	return true
}

// Utility functions
func iptables(action string, queueNum string) {
	var arg string

	switch action {
	case "start":
		arg = "-A"
		// fmt.Println("starting Iptables")

	case "stop":
		arg = "-D"
		// fmt.Println("stopping Iptables")
	}

	cmd := "iptables"
	args := []string{arg, "OUTPUT", "-p", "0", "-j", "NFQUEUE", "--queue-num", queueNum, "--queue-bypass"}
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
		// // fmt.Println("Trouble decoding layers: ", err)
	}

	return foundLayerTypes
}

// func main() {
// 	queueNum := 12

// 	iptables("start", queueNum)

// 	startNFQ(queueNum)

// 	iptables("stop", queueNum)
// }

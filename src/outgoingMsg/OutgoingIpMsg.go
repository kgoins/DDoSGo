package outgoingMsg

type OutgoingIPsMsg struct{
     Agent_ip string
     Agent_port string
     OffendingIPs []string
}

func NewOutgoingIPsMsg(aIP string, port string, ips []string) OutgoingIPsMsg{
     return OutgoingIPsMsg{
     	    Agent_ip: aIP,
	    Agent_port: port,
	    OffendingIPs: ips}
}

func (ipMsg OutgoingIPsMsg) String() string{
     return "type: IPs" + "; payload: agent= " +ipMsg.Agent_ip+ " " + ipMsg.Agent_port
}

func (ipMsg OutgoingIPsMsg) GetType() string{
     return "IPs"
}


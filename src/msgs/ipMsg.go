package msgs

import "cmds"

type IPsMsg struct{
     Agent_ip string
     Agent_port string
     OffendingIPs []string
}

func NewIPsMsg(aIP string, port string, ips []string) IPMsg{
     return IPsMsg{
     	    Agent_ip: aIP,
	    Agent_port: port,
	    OffendingIPs: ips}
}

func (ipMsg IPsMsg) String() string{
     return "type: IPs" + "; payload: agent= " + ipMsg.Agent_ip + " " + ipMsg.Agent_port
}

func (ipMsg IPsMsg) GetType() string{
     return "IPs"
}

func (ipMsg IPsMsg) BuildCmd() cmds.Cmd{
     return cmds.NewIPsCmd(ipMsg.Agent_ip, ipMsg.Agent_port, ipMsg.OffendingIPs)
}

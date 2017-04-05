package msgs

import "cmds"

// IP msg, send to update enforcer ips
type IPsMsg struct{
     Agent_ip string
     Agent_port string
     OffendingIPs []string
}

// Create new IP msg
func NewIPsMsg(aIP string, port string, ips []string) IPsMsg{
     return IPsMsg {
 	    Agent_ip: aIP,
	    Agent_port: port,
	    OffendingIPs: ips,
    }
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

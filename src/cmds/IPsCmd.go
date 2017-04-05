package cmds

import "subsytems"
import "visitors"

type IPsCmd struct{
     enforcer *subsystems.Enforcer
     agent_ip string
     agent_port string
     OffendingIPs []string
}

// Create new IPsCmd
func NewIPsCmd(aIP string, aPort string, ips []string) IPsCmd{
     return IPsCmd{
 	    enforcer: visitors.EnforcerVisitor.EnforcerInstance,
	    agent_ip: aIP,
	    agent_port: aPort,
	    OffendingIPs: ips}


// Execute cmd to update ips
func (IPsCmd ipcmd) ExecCmd() {
  print("Updating enforcer instance ips...")
  ipcmd.enforcer.updateOffendingIps(ipcmd.OffendingIPs)
}

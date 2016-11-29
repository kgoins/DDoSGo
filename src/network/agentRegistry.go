package network

/*AGENT REGISTRY- handler side
Initialized with this.start function

Maintains a list of all of the agents a handler is in charge of
Removes agents when they shuts down or request to be removed
*/

/*MODIFYING AGENTS CLASS
Send a "register" message to handler with critical information about self
	Critical information: Every router between itself and interwebz
*/

/*NEW REGISTER MSG
Utilize interface to create a new "register" message for agent to send to handler
*/

/*LOOK/FEEL OF MSG
Critical info about agent

*/

type AgentRecord struct {
	//What makes up an agent record
	agent_hostname   string
	handler_hostname string
	handler_port     int
	traceroute       []string //An empty list, use append to add to it
}

/**
 *The constructor for the AgentRecord struct
 *@param aName the agent hostname
 *@param hName The handler hostname
 *@param port The port number that the handler-agent connection is on
 *@param list A slice which contains all of the routers between the agent and the interwebz
 *@return A reference to the agent record itself
 */
func NewAgentRecord(aName string, hName string, port int, list []string) *AgentRecord {

	return &AgentRecord{
		agent_hostname:   aName,
		handler_hostname: hName,
		handler_port:     port,
		traceroute:       list}
}

func (rec *AgentRecord) GetAgHostname() string {
	return rec.agent_hostname
}

/************************************************************************************************/

type AgentRegistry struct {
	//Hash map of agents
	//The registry is a hashmap with hostnames as the key which returns the Agent record
	registry map[string]AgentRecord
}

/**
 *The constructor for the AgentRegistry
 *Initializes a hash map to store the future AgentRecords
 */
func Start() *AgentRegistry {
	//Initialize a hash map for agents
	//Return a pointer to an agent registry
	//Is the constructor
	reg := make(map[string]AgentRecord)

	return &AgentRegistry{
		registry: reg}
}

/**
 *Removes an AgentRecord from the AgentRegistry
 *If the hostname is not found in the registry nothing is done
 *@param hostname A string representing the host name of the agent to be removed
 */
func (reg *AgentRegistry) RemoveAgent(hostname string) {
	//Removes an agent from the list of registered agents

	delete(reg.registry, hostname)
	//if hosname does not exist, delete does nothing
}

/**
 *Adds a new AgentRecord to the AgentRegistry
 *@param agent The AgentRecord to be added to the AgentRegistry
 *@see newAgentRecord
 */
func (reg *AgentRegistry) AddAgent(agent AgentRecord) {
	//Adds and agent to the registry (hash map)
	reg.registry[agent.handler_hostname] = agent
}

/**
 *Gives the list of all routers between the agent and the interwebz
 *@param hostname The hostname of the agent to retrieve list for
 *@return A slice containing all of the devices between the agent and the web
 */
func (reg *AgentRegistry) ReturnTrace(hostname string) []string {
	//Returns all routers between an agent and the internet
	//Returns whatever the config file has in it for a particular agent
	//Return arrayList of host names?

	agent, exists := reg.registry[hostname]

	if exists {
		return agent.traceroute
	} else {
		return nil
	}
}

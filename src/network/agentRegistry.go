package network

/*AGENT REGISTRY- handler side
Initialized with this.start function
Accepts an incoming "register" message and adds that agent to the registry
Maintains a list of all of the agents a handler is in charge of
Removes agents when it shuts down or requests to be removed
Note: Each host should only have one agent
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
Add or Remove command
*/

/* EXPECTED DEF OF AGENT*/
type AgentRecord struct {
	//What makes up an agent record
	agent_hostname   string
	handler_hostname string
	handler_port     int
	traceroute       []string //An empty list, use append to add to it
}

type AgentRegistry struct {
	//Hash map of agents
	//The registry is a hashmap with hostnames as the key which returns the Agent record
	registry map[string]AgentRecord
}

func start() *AgentRegistry {
	//Initialize a hash map for agents
	//Return a pointer to an agent registry
	//Is the constructor
	reg := make(map[string]AgentRecord)

	return &AgentRegistry{
		registry: reg}
}

func (reg *AgentRegistry) removeAgent(hostname string) {
	//Removes an agent from the list of registered agents

	delete(reg.registry, hostname)
	//if hosname does not exist, delete does nothing
}

func (reg *AgentRegistry) addAgent(agent AgentRecord) {
	//Adds and agent to the registry (hash map)
	reg.registry[agent.handler_hostname] = agent
}

func (reg *AgentRegistry) returnTrace(hostname string) []string {
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

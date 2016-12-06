package subsystems

import (
	"fmt"
)

/*AGENT REGISTRY- handler side
Initialized with this.start function

Maintains a list of all of the agents a handler is in charge of
Removes agents when they shut down or request to be removed (request should be made before shutdown)
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
	agent_ip     string
	handler_ip   string
	handler_port string
	agent_port   string
	traceroute   []string //An empty list, use append to add to it
	updated      bool
	isFiltering  bool
}

/**
 *The constructor for the AgentRecord struct
 *@param aName the agent hostname
 *@param hName The handler hostname
 *@param port The port number that the handler-agent connection is on
 *@param list A slice which contains all of the routers between the agent and the interwebz
 *@return A reference to the agent record itself
 */
func NewAgentRecord(aIP string, hIP string, port string, aPort string, list []string) *AgentRecord {

	return &AgentRecord{
		agent_ip:     aIP,
		handler_ip:   hIP,
		handler_port: port,
		agent_port:   aPort,
		traceroute:   list,
		updated:      true,
		isFiltering:  false}
}

func (rec *AgentRecord) GetAgHostname() string {
	return rec.agent_ip
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
func NewAgentRegistry() *AgentRegistry {
	//Initialize a hash map for agents
	//Return a pointer to an agent registry
	//Is the constructor
	reg := make(map[string]AgentRecord)
	fmt.Println("Setting Up New Agent Registry...")

	return &AgentRegistry{
		registry: reg}
}

/**
 *Removes an AgentRecord from the AgentRegistry
 *If the hostname is not found in the registry nothing is done
 *@param hostname A string representing the host name of the agent to be removed
 */
func (reg *AgentRegistry) RemoveAgent(agentIP string, agentPort string) {
	//Removes an agent from the list of registered agents
	fmt.Printf("Deleting Agent %s%s From Registry\n", agentIP, agentPort)
	key := agentIP + agentPort
	delete(reg.registry, key)
	//if hosname does not exist, delete does nothing
}

/**
 *Adds a new AgentRecord to the AgentRegistry
 *@param agent The AgentRecord to be added to the AgentRegistry
 *@see newAgentRecord
 */
func (reg *AgentRegistry) AddAgent(agent AgentRecord) {
	//Adds and agent to the registry (hash map)
	key := agent.agent_ip + agent.agent_port
	reg.registry[key] = agent
}

/**
 *Gives the list of all routers between the agent and the interwebz
 *@param hostname The hostname of the agent to retrieve list for
 *@return A slice containing all of the devices between the agent and the web
 */
func (reg *AgentRegistry) ReturnTrace(aIP string, aPort string) []string {
	//Returns all routers between an agent and the internet
	//Returns whatever the config file has in it for a particular agent
	//Return arrayList of host names?
	key := aIP + aPort
	agent, exists := reg.registry[key]

	if exists {
		return agent.traceroute
	} else {
		return nil
	}
}

// Update given hander record in registry
func (reg *AgentRegistry) UpdateRecordStatus(agent_ip string, aPort string) {

	// Check registry for if agent exists, and then update its status to updated
	key := agent_ip + aPort
	agent, exists := reg.registry[key]
	if exists {
		agent.updated = true
		reg.registry[key] = agent
		// fmt.Println("Updating Record of Agent ", agent.agent_ip)
	}

}

// Function to check registry for unresponsive agent records and alert system to problem if found
func (reg *AgentRegistry) CheckRecords() (bool, []AgentRecord) {
	// fmt.Println("Checking Registry Records...")

	var clean = true                      // Test for unresponsive records
	var unresponsiveRecords []AgentRecord // The unresponsive records

	// Loop through registry and check records
	for _, record := range reg.registry {
		if record.updated == true {
			record.updated = false // If record was updated in last interval, mark rest to false
			key := record.agent_ip + record.agent_port
			reg.registry[key] = record // Reset record into registry
		} else {
			clean = false // Record was found unresponsive, we have to signal alerts
			record.isFiltering = true
			key := record.agent_ip + record.agent_port
			reg.registry[key] = record
			unresponsiveRecords = append(unresponsiveRecords, record) // Append record to collection for return to the alert
		}

	}
	return clean, unresponsiveRecords
}

func (reg *AgentRegistry) IsAgentFiltering(agent_ip string, agent_port string) bool {
	return reg.registry[agent_ip+agent_port].isFiltering
}

// Set a current agent in the registry as filtering
func (reg *AgentRegistry) SetAgentAsFiltering(agent_ip string, agent_port string) {
	key := agent_ip + agent_port
	record := reg.registry[key]
	record.isFiltering = true
	reg.registry[key] = record
}

// Turn off filtering for an agent in the registry
func (reg *AgentRegistry) ClearAgentFilteringStatus(agent_ip string, agent_port string) {
	key := agent_ip + agent_port
	record := reg.registry[key]
	record.isFiltering = false
	reg.registry[key] = record
}

package network


import fmt

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

struct AgentRegistry{
	//Hash map of agents
}

struct Agent{
	//What makes up an agent record
	//host name, ip addr, port
	//ArrayList trace
}

func start(){
	//Initialize a hash map for agents	
	//Create a channel for other objects to communicate with it
	//Return a pointer to an agent registry
	//Is the constructor
}

func removeAgent(){
	//Removes an agent from the list of registered agents
	//Throws error if it does not exist
}

func addAgent(){
	//Adds and agent to the registry (hash map)
}

func returnTrace(){
	//Returns all routers between an agent and the internet
	//Returns whatever the config file has in it for a particular agent
	//Return arrayList of host names?
}




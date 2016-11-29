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


func start(){
	//Initialize a hash map for agents	
	//Create a channel for other objects to communicate with it
}

func removeAgent(){
	//Removes an agent from the list of registered agents
	//Throws error if it does not exist
}

func checkAgentCount(){
	//Checks how many agents are on a given host
	//Each host should only have one agent
	//REDACTED??
}

func getChannel(){
	//Returns a reference to the channel for the registry
}

func addAgent(){
	//Adds and agent to the registry (hash map)
}

func returnTrace(){
	//Returns all routers between an agent and the internet
	//Returns whatever the config file has in it for a particular agent
	//Return arrayList of host names?
}




package config

import "encoding/json"
import "io/ioutil"
import "fmt"

type AgentConf struct {
	HandlerAddr string
	HandlerPort int
}

type HandlerConf struct{}

func ReadAgentConf() (AgentConf, error) {
	confPath := "config/agent_config.json"
	var conf AgentConf

	rawFile, fileErr := ioutil.ReadFile(confPath)

	err := json.Unmarshal(rawFile, &conf)

	if err != nil {
		fmt.Println("error reading agent config")
		fmt.Println(err)
		return conf, err
	}

	return conf, nil
}

func ReadHandlerConf() (AgentConf, error) {
	confPath := "./handler_config.json"
	var conf AgentConf

	rawFile, _ := ioutil.ReadFile(confPath)

	err := json.Unmarshal(rawFile, &conf)

	if err != nil {
		fmt.Println("error reading handler config")
		fmt.Println(err)
		return conf, err
	}

	return conf, nil
}

func main() {
	conf, err := ReadAgentConf()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(conf)
}

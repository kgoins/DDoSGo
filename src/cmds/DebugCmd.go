package cmds

import "fmt"

type Debug struct {
	msgText string
	id      int
}

func NewDebug(text string, id int) Debug {
	return Debug{msgText: text, id: id}
}

func (debug Debug) ExecCmd() {
	fmt.Println("From debug cmd: ", debug.msgText, debug.id)
}

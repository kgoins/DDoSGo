package dispatcher

import "net"
import "fmt"
import "io/ioutil"

import "msgs"

// implements dispatcher.Dispatchable
type MsgDispatchable struct {
	conn net.Conn
}

func NewMsgDispatchable(conn net.Conn) MsgDispatchable {
	return MsgDispatchable{conn: conn}
}

func (msgDisp MsgDispatchable) DispatchableExec() {
	defer msgDisp.conn.Close()
	fmt.Println("handling conn from: " + msgDisp.conn.RemoteAddr().String())

	msgBytes, _ := ioutil.ReadAll(msgDisp.conn)
	msg := msgs.BuildMsg(msgBytes)
	fmt.Println("received msg: " + msg.String())
}

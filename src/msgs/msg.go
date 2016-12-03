package msgs

// import "Cmd"
import "encoding/json"
import "fmt"

type Msg interface {
	// BuildCommand()
	String() string
	GetType() string
}

type MsgBuilder struct {
	MsgType string
	Payload interface{}
}

func EncodeMsg(msg Msg) []byte {
	// msgData, _ := json.Marshal(msg)
	builder := MsgBuilder{MsgType: msg.GetType(), Payload: msg}

	builderData, _ := json.Marshal(builder)

	return builderData
}

func BuildMsg(msgBytes []byte) Msg {
	var rawMsg json.RawMessage
	builder := MsgBuilder{
		Payload: &rawMsg,
	}

	if err := json.Unmarshal(msgBytes, &builder); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Builder Msg Type: ", builder.MsgType)
	switch builder.MsgType {
	case "Debug":
		msg := Debug{}
		json.Unmarshal(rawMsg, &msg)
		return msg

	case "DataStream":
		msg := DataStream{}
		json.Unmarshal(rawMsg, &msg)
		return msg

	default:
		fmt.Println("unknown message type: %s", builder.MsgType)
		return nil
	}
}

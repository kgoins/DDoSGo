package msgs

// import "Cmd"

type Msg interface {
	// BuildCommand() Cmd
	String() string
}

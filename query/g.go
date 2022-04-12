package query

import (
	"fmt"
	"github.com/c-bata/go-prompt"
)

const (
	Empty = ""
	Space = " "

	// desc
	K8sShortcutDesc = "k8s resource shortcut query."

	// all namespace
	AllNamespace = "all"

	// k8s resource kind
	KindDaemonSet   = "DaemonSet"
	KindDeployment  = "Deployment"
	KindReplicaSet  = "ReplicaSet"
	KindStatefulSet = "StatefulSet"
	KindPod         = "Pod"
	KindNode        = "Node"
	KindService     = "Service"
	KindJob         = "Job"
)

var (
	// config
	Debug bool

	ResourcesList       []CmdResourceInterface
	ConsoleStdoutWriter prompt.ConsoleWriter
	// ConsoleStderrWriter prompt.ConsoleWriter

	// use namespace
	GlobalNamespace = "default"
)

type CompareAct string

var (
	Equal        CompareAct = "=="
	GreaterEqual CompareAct = ">="
	LessEqual    CompareAct = "<="
	Greater      CompareAct = ">"
	Less         CompareAct = "<"
)

func MsgExpectLineWord(cnt int) string {
	return fmt.Sprintf("Expect line word == %d", cnt)
}

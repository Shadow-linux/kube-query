package query

import (
	"context"
	"fmt"
	"golang.org/x/term"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/remotecommand"
	"os"
	"time"
)

func init() {
	// wrap exec error.
	runtime.ErrorHandlers = []func(error){
		func(err error) {
			Logger("Exec handle error: ", err.Error())
		},
	}
}

type termSizeQueue chan remotecommand.TerminalSize

func (this termSizeQueue) Next() *remotecommand.TerminalSize {
	size, ok := <-this
	if !ok {
		return nil
	}
	return &size
}

func PodExec(pod *v1.Pod, container, cmd string) string {

	fmt.Printf("Connect to container: %s.%s\n", pod.Name, container)
	fmt.Printf("Commannd: %s \n", cmd)
	req := ClientSet.CoreV1().RESTClient().Post().Resource("pods").
		Namespace(pod.GetNamespace()).Name(pod.GetName()).
		SubResource("exec").
		Param("color", "false").
		VersionedParams(&v1.PodExecOptions{
			Container: container,
			Command:   []string{cmd},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)
	Logger("Request: ", req.URL())

	exec, err := remotecommand.NewSPDYExecutor(RestCliConfig, "POST", req.URL())
	if err != nil {
		WrapError(fmt.Errorf("unable to execute remote command, err: %s", err.Error()))
		return "unable to execute remote command. err: " + err.Error()
	}
	fd := int(os.Stdin.Fd())
	//// Put the terminal into raw mode to prevent it echoing characters twice.
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		WrapError(fmt.Errorf("unable to init terminal, err: %s", err.Error()))
		return "unable to init terminal. err: " + err.Error()
	}

	// init termSize
	s := make(termSizeQueue, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			termWidth, termHeight, err := term.GetSize(fd)
			WrapError(err)
			termSize := remotecommand.TerminalSize{Width: uint16(termWidth), Height: uint16(termHeight)}
			s <- termSize

			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(2 * time.Second)
			}
		}

	}()

	defer func() {
		err := term.Restore(fd, oldState)
		if err != nil {
			WrapError(err)
		}
	}()

	// Connect this process' std{in,out,err} to the remote shell process.
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:             os.Stdin,
		Stdout:            os.Stdout,
		Stderr:            os.Stderr,
		Tty:               true,
		TerminalSizeQueue: s,
	})

	if err != nil {
		if err != io.EOF {
			WrapError(fmt.Errorf("unable to stream shell process, err: %s", err.Error()))
			return "unable to stream shell process. err: " + err.Error()
		}

	}
	return "Connection closed."
}

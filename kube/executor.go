package kube

import (
	"bytes"
	"github.com/Shadow-linux/kube-query/query"
	"os"
	"os/exec"
	"strings"
)

func ExecuteAndGetResult(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		query.Logger("you need to pass the something arguments")
		return ""
	}

	out := &bytes.Buffer{}
	cmd := exec.Command("/bin/sh", "-c", "kubectl "+s)
	cmd.Stdin = os.Stdin
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		query.Logger(err.Error())
		return ""
	}
	r := string(out.Bytes())
	return r
}

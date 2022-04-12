package query

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const (
	CmdHelpInfo = `
Kube-query is a kubectl plugin what easier to query and operate Kubernetes cluster.
  Find more information at: https://github.com/Shadow-linux/kube-query

Basic Commands (beginner):
  clear			Clear console.
  use			Set namespace context. like: use default , use all, use kube-system.
  @				Start run shell command. like: '@ls /tmp' or '@ ls /tmp'
  help			Print help information.
  exit			Exit console.

Resource Commands (quick query):
  pods			Get resource information, relationship and login into container.
  deploy		Get resource information, relationship.
  svc			Get resource information, relationship.
  daemonset		Get resource information, relationship.
  configmap		Get resource information, relationship.
  job			Get resource information, relationship.

---

Native kubectl controls the Kubernetes cluster manager.

 Find more information at: https://kubernetes.io/docs/reference/kubectl/overview/

Basic Commands (Beginner):
  create        Create a resource from a file or from stdin.
  expose        Take a replication controller, service, deployment or pod and expose it as a new Kubernetes Service
  run           Run a particular image on the cluster
  set           Set specific features on objects

Basic Commands (Intermediate):
  explain       Documentation of resources
  get           Display one or many resources
  edit          Edit a resource on the server
  delete        Delete resources by filenames, stdin, resources and names, or by resources and label selector

Deploy Commands:
  rollout       Manage the rollout of a resource
  scale         Set a new size for a Deployment, ReplicaSet or Replication Controller
  autoscale     Auto-scale a Deployment, ReplicaSet, or ReplicationController

Cluster Management Commands:
  certificate   Modify certificate resources.
  cluster-info  Display cluster info
  top           Display Resource (CPU/Memory/Storage) usage.
  cordon        Mark node as unschedulable
  uncordon      Mark node as schedulable
  drain         Drain node in preparation for maintenance
  taint         Update the taints on one or more nodes

Troubleshooting and Debugging Commands:
  describe      Show details of a specific resource or group of resources
  logs          Print the logs for a container in a pod
  attach        Attach to a running container
  exec          Execute a command in a container
  port-forward  Forward one or more local ports to a pod
  proxy         Run a proxy to the Kubernetes API server
  cp            Copy files and directories to and from containers.
  auth          Inspect authorization
  debug         Create debugging sessions for troubleshooting workloads and nodes

Advanced Commands:
  diff          Diff live version against would-be applied version
  apply         Apply a configuration to a resource by filename or stdin
  patch         Update field(s) of a resource
  replace       Replace a resource by filename or stdin
  wait          Experimental: Wait for a specific condition on one or many resources.
  kustomize     Build a kustomization target from a directory or a remote url.

Settings Commands:
  label         Update the labels on a resource
  annotate      Update the annotations on a resource
  completion    Output shell completion code for the specified shell (bash or zsh)

Other Commands:
  api-resources Print the supported API resources on the server
  api-versions  Print the supported API versions on the server, in the form of "group/version"
  config        Modify kubeconfig files
  plugin        Provides utilities for interacting with plugins.
  version       Print the client and server version information
	
`
)

func PureCmdRun(line string, autoOutput bool) string {
	Logger("PureCmdRun: ", "/bin/sh", "-c", line)
	cmd := exec.Command("/bin/sh", "-c", line)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	if autoOutput {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			Logger("Got error: ", err.Error())
		}
	} else {
		var res string
		out, err := cmd.CombinedOutput()
		res = string(out)
		if err != nil {
			res = res + "\n" + fmt.Sprintf("Run cmd error: %s", err.Error())
		}
		return strings.TrimSpace(res)
	}

	return ""
}

func CmdRunWithFile(cmd string, autoOutput bool) string {
	const (
		tmpDir = "/tmp"
	)
	content := fmt.Sprintf(`#!/bin/bash
source /etc/profile;
%s`, cmd)
	file, err := ioutil.TempFile(tmpDir, "*.kube-query")
	WrapError(err)
	defer func() {
		_ = file.Close()
		os.Remove(file.Name())
	}()
	_, err = file.WriteString(content)
	WrapError(err)

	execCMD := fmt.Sprintf("sh %s", file.Name())
	Logger("CmdRunWithFile: ", execCMD)
	return PureCmdRun(execCMD, autoOutput)
}

func fetchArgNs() string {
	var argsNs string
	if GlobalNamespace == "all" {
		argsNs = "-A"
	} else {
		argsNs = "-n " + GlobalNamespace
	}
	return argsNs
}

func K8sCmdRun(line string, autoOutput bool) string {
	var cmdLine string
	if strings.Contains(line, "|") {
		lines := strings.Split(line, "|")
		firstColumn := lines[0] + " " + fetchArgNs()
		cmdLine = fmt.Sprintf("kubectl %s | %s", firstColumn, strings.Join(lines[1:], " "))
	} else {
		cmdLine = fmt.Sprintf("kubectl %s %s", line, fetchArgNs())
	}
	Logger("K8sCmdRun: ", cmdLine)
	return PureCmdRun(cmdLine, autoOutput)
}

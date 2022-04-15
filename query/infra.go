package query

import (
	"bytes"
	"context"
	"fmt"
	"github.com/c-bata/go-prompt"
	"io/ioutil"
	appv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"strings"
)

type CmdResRelationshipInterface interface {
	// generate relationship
	generate() []func() string
}

type CmdResourceInterface interface {
	Name() string
	ShortName() string
	CanExecute(ctx *PromptCtx) bool
	Execute(ctx *PromptCtx) string
	DefaultSuggestions(ctx *PromptCtx, name, shortName string) []prompt.Suggest
	Suggestions(ctx *PromptCtx) []prompt.Suggest
}

type CmdResSuggestionInterface interface {
	Helper() []prompt.Suggest
}

type PromptCtx struct {
	context.Context
	Line string
	Word string
	D    prompt.Document
}

func NewPromptCtx(context context.Context, line string, word string, d prompt.Document) *PromptCtx {
	return &PromptCtx{Context: context, Line: line, Word: word, D: d}
}

type cmdResourceRelaTool struct {
}

func newCmdResourceRelaTool() *cmdResourceRelaTool {
	return &cmdResourceRelaTool{}
}

func (this *cmdResourceRelaTool) setTitle(title *string, resourceKind string) {
	if *title == "" {
		*title = resourceKind
	}
}

func (this *cmdResourceRelaTool) setHeader(header *[]string, headerValue []string) {
	if len(*header) == 0 {
		*header = headerValue
	}
}

func (this *cmdResourceRelaTool) FormatView(objects []runtime.Object) string {
	var (
		title  string
		header = make([]string, 0)
		buffer = bytes.NewBuffer(make([]byte, 10240))
	)

	table := NewTable(buffer)
	for _, obj := range objects {
		switch obj.(type) {
		case *v1.Pod:
			o := obj.(*v1.Pod)

			this.setTitle(&title, KindPod)
			this.setHeader(&header, []string{"Name", "READY", "STATUS", "RESTARTS", "IP", "NODE", "NOMINATED NODE"})
			var (
				readyCnt   int
				restartCnt int32
			)
			for _, cs := range o.Status.ContainerStatuses {
				if cs.Ready {
					readyCnt += 1
				}
				restartCnt += cs.RestartCount
			}

			table.Append([]string{
				o.Name,
				fmt.Sprintf("%d/%d", readyCnt, len(o.Status.ContainerStatuses)),
				string(o.Status.Phase),
				fmt.Sprintf("%d", restartCnt),
				o.Status.PodIP,
				o.Spec.NodeName,
				o.Status.NominatedNodeName,
			})
			break

		case *v1.Service:
			this.setTitle(&title, KindService)
			this.setHeader(&header, []string{"Name", "Type", "Cluster-IP", "External-IP", "PORTS"})

			o := obj.(*v1.Service)
			table.Append([]string{
				o.Name,
				string(o.Spec.Type),
				o.Spec.ClusterIP,
				SliceString2String(o.Spec.ExternalIPs),
				ParserServicePorts2String(o.Spec.Ports),
			})
			break
		case *appv1.DaemonSet:
			this.setTitle(&title, KindDaemonSet)
			this.setHeader(&header, []string{"Name", "DESIRED", "CURRENT", "READY", "UP-TO-DATE", "AVAILABLE"})
			o := obj.(*appv1.DaemonSet)
			table.Append([]string{
				o.Name,
				Int2String(o.Status.DesiredNumberScheduled),
				Int2String(o.Status.CurrentNumberScheduled),
				Int2String(o.Status.NumberReady),
				Int2String(o.Status.UpdatedNumberScheduled),
				Int2String(o.Status.NumberAvailable),
				Map2String(o.Spec.Selector.MatchLabels),
			})
			break
		case *appv1.Deployment:
			this.setTitle(&title, KindDeployment)
			this.setHeader(&header, []string{"Name", "READY", "UP-TO-DATE", "AVAILABLE"})
			o := obj.(*appv1.Deployment)
			table.Append([]string{
				o.Name,
				Int2String(o.Status.ReadyReplicas) + "/" + Int2String(o.Status.Replicas),
				Int2String(o.Status.UpdatedReplicas),
				Int2String(o.Status.AvailableReplicas),
			})
			break
		case *appv1.StatefulSet:
			this.setTitle(&title, KindStatefulSet)
			this.setHeader(&header, []string{"Name", "READY", "UP-TO-DATE", "AVAILABLE"})
			o := obj.(*appv1.StatefulSet)
			table.Append([]string{
				o.Name,
				Int2String(o.Status.ReadyReplicas) + "/" + Int2String(o.Status.Replicas),
				Int2String(o.Status.UpdatedReplicas),
				Int2String(o.Status.AvailableReplicas),
			})
			break
		case *appv1.ReplicaSet:
			this.setTitle(&title, KindReplicaSet)
			this.setHeader(&header, []string{"Name", "DESIRED", "CURRENT", "READY"})
			o := obj.(*appv1.ReplicaSet)
			table.Append([]string{
				o.Name,
				Int2String(o.Status.Replicas),
				Int2String(o.Status.AvailableReplicas),
				Int2String(o.Status.ReadyReplicas),
			})
			break
		case *batchv1.Job:
			this.setTitle(&title, KindReplicaSet)
			this.setHeader(&header, []string{"Name", "COMPLETIONS", "DURATION", "CONTAINERS", "IMAGES", "SELECTOR"})
			o := obj.(*batchv1.Job)
			var (
				containerNames []string
				imgs           []string
			)
			for _, c := range o.Spec.Template.Spec.Containers {
				imgs = append(imgs, c.Image)
				containerNames = append(containerNames, c.Name)
			}
			table.Append([]string{
				o.Name,
				Int2String(o.Status.Active) + "/" + Int2String(o.Spec.Completions),
				Int2String(0),
				SliceString2String(containerNames),
				SliceString2String(imgs),
				Map2String(o.Spec.Selector.MatchLabels),
			})
			break
		default:
			return fmt.Sprintf("%s has no relevant resource.\n", obj.GetObjectKind().GroupVersionKind().Kind)
		}
	}
	if title != "" {
		table.SetHeader(header)
		table.Render()
		b, err := ioutil.ReadAll(buffer)
		WrapError(err)
		content := "##### " + title + " #####" + "\n" + string(b)
		return content
	}
	return ""
}

func (this *cmdResourceRelaTool) FetchAll(funcs []func() string) string {
	var (
		res  []string
		body []string
	)

	for _, f := range funcs {
		s := f()
		if strings.TrimSpace(s) != "" {
			body = append(body, s)
		}
	}
	// title
	header := []string{"Relevant relationship:"}
	res = append(res, header...)
	if len(body) == 0 {
		res = append(res, "No relevant resource.")

	} else {
		res = append(res, body...)

	}
	return strings.Join(res, "\n")

}

type cmdResourceTool struct{}

func newCmdResourceTool() *cmdResourceTool {
	return &cmdResourceTool{}
}

func (this *cmdResourceTool) DefaultSuggestions(ctx *PromptCtx, name, shortName string) []prompt.Suggest {
	Logger("Run DefaultSuggestions.")
	return prompt.FilterHasPrefix(
		[]prompt.Suggest{
			{name, K8sShortcutDesc},
			{shortName, K8sShortcutDesc},
		},
		ctx.Word,
		true)
}

func (this *cmdResourceTool) FindResourceName(ctx *PromptCtx) string {
	lines := FormatLineWithSpace(ctx.Line)
	if len(lines) >= 2 {
		return lines[1]
	}
	return ""
}

type cmdResourceExecutorTool struct {
}

func newCmdResourceExecutorTool() *cmdResourceExecutorTool {
	return &cmdResourceExecutorTool{}
}

func (this *cmdResourceExecutorTool) Output(name, ns, mode, kind string, obj runtime.Object) string {
	if mode == "yaml" {
		return RuntimeObject2Yaml(obj)
	}
	if mode == "json" {
		return RuntimeObject2Json(obj)
	}
	return K8sCmdRun(fmt.Sprintf("describe %v %s", kind, name), false, ns)
}

func (this *cmdResourceExecutorTool) Events(uid types.UID) string {
	contents := []string{
		"Events: ",
	}
	events := FetchEvents(uid)
	var (
		buffer = bytes.NewBuffer(make([]byte, 10240))
	)
	tb := NewTable(buffer)
	tb.SetHeader([]string{"TYPE", "REASON", "MESSAGE"})
	for _, e := range events {
		tb.Append([]string{
			e.Type,
			e.Reason,
			e.Message,
		})
	}
	tb.Render()
	b, e := ioutil.ReadAll(buffer)
	WrapError(e)
	contents = append(contents, string(b))
	return strings.Join(contents, "\n")
}

func (this *cmdResourceExecutorTool) Labels(labels map[string]string) string {
	contents := []string{
		"Labels: ",
	}
	labelSlice := Map2Slice(labels)
	contents = append(contents, labelSlice...)
	return strings.Join(contents, "\n")
}

func (this *cmdResourceExecutorTool) Annotations(annos map[string]string) string {
	contents := []string{
		"Annotations: ",
	}
	annoSlice := Map2Slice(annos)
	contents = append(contents, annoSlice...)
	return strings.Join(contents, "\n")
}

func RegisterCmdResources(resource CmdResourceInterface) {
	ResourcesList = append(ResourcesList, resource)
}

package query

import (
	"bytes"
	"fmt"
	"github.com/c-bata/go-prompt"
	"io/ioutil"
	appv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
	"sort"
	"strings"
)

func init() {
	RegisterCmdResources(NewPod())
}

type Pod struct {
	cmdResourceTool
	suggest *PodsSuggestion
	exec    *PodsExecutor
}

func NewPod() *Pod {
	return &Pod{cmdResourceTool: cmdResourceTool{}}
}

func (this *Pod) Name() string {
	return "pods"
}

func (this *Pod) ShortName() string {
	return "po"
}

func (this *Pod) CanExecute(ctx *PromptCtx) bool {
	if RuleHasPrefix(
		ctx.Line,
		this.Name(),
		this.ShortName()) {
		return true
	}
	return false
}

func (this *Pod) Execute(ctx *PromptCtx) string {
	if this.CanExecute(ctx) {
		var pod *v1.Pod
		this.exec = NewPodsExecutor(ctx)
		// pods <PodName>
		if RuleJudgeLineWordCount(ctx.Line, 2, Equal) {
			pod = FetchPodWithName(ParserResourceName(this.FindResourceName(ctx)))
			if pod == nil {
				return ""
			}
			return this.exec.Output(pod, "")
		}

		if RuleJudgeLineWordCount(ctx.Line, 2, Greater) {

			pod = FetchPodWithName(ParserResourceName(this.FindResourceName(ctx)))
			if pod == nil {
				return ""
			}

			// pods <PodName> -o [yaml|desc|json]
			if RuleJudgeLineHasWords(ctx.Line, ArgOutput.Text) {
				return this.exec.Output(pod, ArgOutput.Text)
			}
			// pods <PodName> -r
			if RuleJudgeLineHasWords(ctx.Line, ArgRelationship.Text) {
				return this.exec.Relationship(pod)
			}
			// pods <PodName> -l
			if RuleJudgeLineHasWords(ctx.Line, ArgLabel.Text) {
				return this.exec.Labels(pod)
			}
			// pods <PodName> -i <ContainerName> -s sh
			if RuleJudgeLineHasWords(ctx.Line, ArgShell.Text, ArgInteractive.Text) {
				return this.exec.Interactive(pod, ArgInteractive.Text, ArgShell.Text)
			}
			// pods <PodName> -i <ContainerName>
			if RuleJudgeLineHasWords(ctx.Line, ArgInteractive.Text) {
				return this.exec.Interactive(pod, ArgInteractive.Text, "")
			}
			// pods <PodName> -e
			if RuleJudgeLineHasWords(ctx.Line, ArgEvents.Text) {
				return this.exec.Events(pod)
			}
			// pods <PodName> -v
			if RuleJudgeLineHasWords(ctx.Line, ArgVolumes.Text) {
				return this.exec.Volumes(pod)
			}
			// pods <PodName> -a
			if RuleJudgeLineHasWords(ctx.Line, ArgServiceAccount.Text) {
				return this.exec.ServiceAccount(pod)
			}

		}
	}
	return ""
}

func (this *Pod) Suggestions(ctx *PromptCtx) []prompt.Suggest {
	Logger("Run pods suggestions.")

	if this.CanExecute(ctx) && RuleCanRemind(ctx.D) {
		this.suggest = NewPodsSuggestion(ctx)

		// pods <PodName>
		if RuleJudgeLineWordCount(ctx.Line, 2, Equal) {
			return this.suggest.ArgsPodName()
		}
		if RuleJudgeLineWordCount(ctx.Line, 2, Greater) {
			// pods <PodName> -o [yaml|desc|json]
			if RuleJudgeWordExists(ctx.D.GetWordBeforeCursorWithSpace(), "-o ") {
				return this.suggest.ArgsOutput()
			}
			// pods <PodName> -i <ContainerName>
			if RuleJudgeWordExists(ctx.D.GetWordBeforeCursorWithSpace(), "-i ") {
				return this.suggest.ArgsInteractive()
			}
			// pods <PodName> -i <ContainerName> -s sh
			if RuleJudgeWordExists(ctx.D.GetWordBeforeCursorWithSpace(), "-s ") {
				return this.suggest.ArgsShell()
			}
			// pods <PodName> -r

		}
		if RuleCanRemindHelper(ctx.D) {
			return this.suggest.Helper()
		}

	}

	return []prompt.Suggest{}
}

type PodsSuggestion struct {
	ctx *PromptCtx
}

func NewPodsSuggestion(ctx *PromptCtx) *PodsSuggestion {
	return &PodsSuggestion{ctx: ctx}
}

func (this PodsSuggestion) Helper() []prompt.Suggest {
	return []prompt.Suggest{
		ArgOutput,
		ArgRelationship,
		ArgInteractive,
		ArgShell,
		ArgEvents,
		ArgLabel,
		ArgServiceAccount,
		ArgVolumes,
	}
}

func (this PodsSuggestion) ArgsPodName() []prompt.Suggest {
	podList := FetchPods(GlobalNamespace)
	suggestions := make([]prompt.Suggest, 0)
	for _, pod := range podList {
		suggestions = append(suggestions, prompt.Suggest{
			Text: FormatResourceName(pod.Name, pod.Namespace),
			Description: fmt.Sprintf("Status: %s, PodIP: %s, HostIP: %s, Node: %v",
				pod.Status.Phase, pod.Status.PodIP, pod.Status.HostIP, pod.Spec.NodeName,
			),
		})
	}
	sort.Sort(SuggestSlice(suggestions))
	return prompt.FilterHasPrefix(suggestions, this.ctx.Word, true)
}

// o
func (this PodsSuggestion) ArgsOutput() []prompt.Suggest {
	return prompt.FilterHasPrefix([]prompt.Suggest{
		ModeDesc,
		ModeYAML,
		ModeJson,
	}, this.ctx.Word, false)
}

// i
func (this PodsSuggestion) ArgsInteractive() []prompt.Suggest {
	podName := FetchFirstArg(this.ctx.Line)
	pod := FetchPodWithName(ParserResourceName(podName))
	res := make([]prompt.Suggest, 0)

	for _, c := range pod.Spec.Containers {
		res = append(res, prompt.Suggest{
			Text:        c.Name,
			Description: "[Container] Image: " + c.Image,
		})
	}
	return res
}

// s
func (this PodsSuggestion) ArgsShell() []prompt.Suggest {

	return []prompt.Suggest{
		{Text: "sh", Description: "bash shell command"},
		{Text: "/bin/bash", Description: "bash shell command"},
	}
}

// r
func (this PodsSuggestion) ArgsRelationship() []prompt.Suggest {
	return []prompt.Suggest{}
}

type PodsExecutor struct {
	ctx *PromptCtx
}

func NewPodsExecutor(ctx *PromptCtx) *PodsExecutor {
	return &PodsExecutor{ctx: ctx}
}

// o
func (this PodsExecutor) Output(pod *v1.Pod, arg string) string {
	word := GetWordAfterArgWithSpace(this.ctx.Line, arg)
	pod.ManagedFields = []metav1.ManagedFieldsEntry{}
	if word == "yaml" {
		return RuntimeObject2Yaml(pod)
	}
	if word == "json" {
		return RuntimeObject2Json(pod)
	}
	return K8sCmdRun("describe pods  "+pod.Name, false)
}

// i
func (this PodsExecutor) Interactive(pod *v1.Pod, arg1, arg2 string) string {
	// container
	Logger("Pod interactive arg1: ", arg1, " arg2: ", arg2)
	container := GetWordAfterArgWithSpace(this.ctx.Line, arg1)
	var cmd string
	if arg2 == "" {
		cmd = "sh"
	} else {
		cmd = GetWordAfterArgWithSpace(this.ctx.Line, arg2)
	}
	Logger("Container: ", container, " Command: ", cmd)
	return PodExec(pod, container, cmd)
}

// r, search service, deployment, rs, daemonset, job
func (this PodsExecutor) Relationship(pod *v1.Pod) string {
	var res string
	pr := NewPodRelationship(pod)
	res = pr.Tool.FetchAll(pr.generate())
	return res
}

// l
func (this PodsExecutor) Labels(pod *v1.Pod) string {
	contents := []string{
		"Labels: ",
	}
	labelSlice := Map2Slice(pod.Labels)
	contents = append(contents, labelSlice...)
	return strings.Join(contents, "\n")
}

// a
func (this PodsExecutor) ServiceAccount(pod *v1.Pod) string {
	return "ServiceAccount: " + pod.Spec.ServiceAccountName
}

// v
func (this PodsExecutor) Volumes(pod *v1.Pod) string {
	// volumes
	volumes := []string{
		"##### Volumes #####",
	}
	bVolumes, err := yaml.Marshal(pod.Spec.Volumes)
	WrapError(err)
	volumes = append(volumes, string(bVolumes))
	// volume mounts
	volumeMounts := []string{
		"##### VolumeMounts #####",
	}
	for _, c := range pod.Spec.Containers {
		volumeMounts = append(volumeMounts, "[container] "+c.Name)
		bVolumesMount, err := yaml.Marshal(c.VolumeMounts)
		WrapError(err)
		volumeMounts = append(volumeMounts, string(bVolumesMount))

	}
	contents := []string{
		strings.Join(volumes, "\n"),
		strings.Join(volumeMounts, "\n"),
	}

	return strings.Join(contents, "\n")
}

// e
func (this PodsExecutor) Events(pod *v1.Pod) string {
	contents := []string{
		"Events: ",
	}
	events := FetchEvents(pod.UID)
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

type PodRelationship struct {
	Tool   *cmdResourceRelaTool
	pod    *v1.Pod
	ds     *appv1.DaemonSet
	rs     *appv1.ReplicaSet
	deploy *appv1.Deployment
	job    *batchv1.Job
	sts    *appv1.StatefulSet
	svc    []*v1.Service
}

func NewPodRelationship(pod *v1.Pod) *PodRelationship {
	return &PodRelationship{
		Tool: newCmdResourceRelaTool(),
		pod:  pod,
	}
}

func (this *PodRelationship) Service() string {
	var relaObjs []*v1.Service
	objects := FetchServices(GlobalNamespace)
	for _, obj := range objects {
		if RuleJudgeLabelSelectorMatch(obj.Spec.Selector, this.pod.Labels) {
			relaObjs = append(relaObjs, obj)
		}
	}
	this.svc = relaObjs
	Logger(fmt.Sprintf("Relevant services count: %+v", len(relaObjs)))
	return this.Tool.FormatView(SliceResource2SliceRuntimeObj(relaObjs))
}

func (this *PodRelationship) Deployment() string {
	var res string
	objects := FetchDeployments(GlobalNamespace)
	if this.rs == nil {
		return res
	}

	for _, owner := range this.rs.OwnerReferences {
		if owner.Kind == KindDeployment {
			for _, obj := range objects {
				if obj.Name == owner.Name {
					Logger("Relevant deploy name: ", obj.Name)
					res = this.Tool.FormatView(SliceResource2SliceRuntimeObj([]runtime.Object{obj}))
					this.deploy = obj
					return res
				}
			}
		}
	}
	return res
}

func (this *PodRelationship) DaemonSet() string {
	var res string
	for _, owner := range this.pod.OwnerReferences {
		if owner.Kind == KindDaemonSet {
			objs := FetchDaemonSets(GlobalNamespace)
			for _, obj := range objs {
				if obj.Name == owner.Name {
					Logger("Relevant ds name: ", obj.Name)
					res = this.Tool.FormatView(SliceResource2SliceRuntimeObj([]runtime.Object{obj}))
					this.ds = obj
					return res
				}
			}
		}
	}
	return res
}

func (this *PodRelationship) StatefulSet() string {
	var res string
	for _, owner := range this.pod.OwnerReferences {
		if owner.Kind == KindStatefulSet {
			objs := FetchStatefulSets(GlobalNamespace)
			for _, obj := range objs {
				if obj.Name == owner.Name {
					Logger("Relevant sts name: ", obj.Name)
					res = this.Tool.FormatView(SliceResource2SliceRuntimeObj([]runtime.Object{obj}))
					this.sts = obj
					return res
				}
			}
		}
	}
	return res
}

func (this *PodRelationship) ReplicaSet() string {
	var res string
	for _, owner := range this.pod.OwnerReferences {
		if owner.Kind == KindReplicaSet {
			objs := FetchReplicaSets(GlobalNamespace)
			for _, obj := range objs {
				if obj.Name == owner.Name {
					Logger("Relevant rs name: ", obj.Name)
					res = this.Tool.FormatView(SliceResource2SliceRuntimeObj([]runtime.Object{obj}))
					this.rs = obj
					return res
				}
			}
		}
	}
	return res
}

func (this *PodRelationship) Job() string {
	var res string
	for _, owner := range this.pod.OwnerReferences {
		if owner.Kind == KindJob {
			objs := FetchJobs(GlobalNamespace)
			for _, obj := range objs {
				if obj.Name == owner.Name {
					Logger("Relevant rs name: ", obj.Name)
					res = this.Tool.FormatView(SliceResource2SliceRuntimeObj([]runtime.Object{obj}))
					this.job = obj
					return res
				}
			}
		}
	}
	return res
}

func (this *PodRelationship) generate() []func() string {
	return []func() string{
		this.Service,
		this.ReplicaSet,
		this.Deployment,
		this.DaemonSet,
		this.StatefulSet,
		this.Job,
	}
}

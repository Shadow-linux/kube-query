package query

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
)

func init() {
	RegisterCmdResources(NewDaemonSets())
}

type DaemonSets struct {
	cmdResourceTool
	suggest *DaemonSetsSuggestion
	exec    *DaemonSetsExecutor
}

func NewDaemonSets() *DaemonSets {
	return &DaemonSets{cmdResourceTool: cmdResourceTool{}}
}

func (this *DaemonSets) Name() string {
	return "daemonsets"
}

func (this *DaemonSets) ShortName() string {
	return "ds"
}

func (this *DaemonSets) CanExecute(ctx *PromptCtx) bool {
	if RuleHasPrefix(
		ctx.Line,
		this.Name(),
		this.ShortName()) {
		return true
	}
	return false
}

func (this *DaemonSets) Execute(ctx *PromptCtx) string {
	if this.CanExecute(ctx) {
		var ds *appsv1.DaemonSet
		this.exec = NewDaemonSetsExecutor(ctx)

		if RuleJudgeLineWordCount(ctx.Line, 2, Equal) {
			ds = FetchDaemonSetWithName(ParserResourceName(this.FindResourceName(ctx)))
			if ds == nil {
				return ""
			}
			return this.exec.Output(ds, "")
		}

		if RuleJudgeLineWordCount(ctx.Line, 2, Greater) {

			ds = FetchDaemonSetWithName(ParserResourceName(this.FindResourceName(ctx)))
			if ds == nil {
				return ""
			}
			// ds <name> -o [yaml|desc|json]
			if RuleJudgeLineHasWords(ctx.Line, ArgOutput.Text) {
				return this.exec.Output(ds, ArgOutput.Text)
			}
			// ds <name> -l
			if RuleJudgeLineHasWords(ctx.Line, ArgLabel.Text) {
				return this.exec.Labels(ds)
			}
			// ds <name> -e
			if RuleJudgeLineHasWords(ctx.Line, ArgEvents.Text) {
				return this.exec.Events(ds)
			}
			if RuleJudgeLineHasWords(ctx.Line, ArgAnnotaions.Text) {
				return this.exec.Tool.Annotations(ds.Annotations)
			}
			// ds <name> -r
			if RuleJudgeLineHasWords(ctx.Line, ArgRelationship.Text) {
				return this.exec.Relationship(ds)
			}
		}
	}
	return ""
}

func (this *DaemonSets) Suggestions(ctx *PromptCtx) []prompt.Suggest {
	Logger("Run ds suggestions.")
	if this.CanExecute(ctx) && RuleCanRemind(ctx.D) {
		this.suggest = NewDaemonSetsSuggestion(ctx)

		// ds <dsName>
		if RuleJudgeLineWordCount(ctx.Line, 2, Equal) {
			return this.suggest.ArgDsName()
		}

		if RuleJudgeLineWordCount(ctx.Line, 2, Greater) {
			// ds <dsName> -o [yaml|desc|json]
			if RuleJudgeWordExists(ctx.D.GetWordBeforeCursorWithSpace(), "-o ") {
				return this.suggest.ArgsOutput()
			}
		}

		if RuleCanRemindHelper(ctx.D) {
			return this.suggest.Helper()
		}
	}
	return []prompt.Suggest{}
}

type DaemonSetsSuggestion struct {
	ctx *PromptCtx
}

func NewDaemonSetsSuggestion(ctx *PromptCtx) *DaemonSetsSuggestion {
	return &DaemonSetsSuggestion{ctx: ctx}
}

func (this DaemonSetsSuggestion) Helper() []prompt.Suggest {
	return []prompt.Suggest{
		ArgOutput,
		ArgRelationship,
		ArgLabel,
		ArgEvents,
		ArgAnnotaions,
	}
}

func (this DaemonSetsSuggestion) ArgDsName() []prompt.Suggest {
	dss := FetchDaemonSets(GlobalNamespace)
	suggestions := make([]prompt.Suggest, 0)
	for _, ds := range dss {
		suggestions = append(suggestions, prompt.Suggest{
			Text: FormatResourceName(ds.Name, ds.Namespace),
			Description: fmt.Sprintf(
				"Name: %s, Desired: %d, Avaiable: %d",
				ds.Name,
				ds.Status.DesiredNumberScheduled,
				ds.Status.NumberAvailable,
			),
		})
	}
	sort.Sort(SuggestSlice(suggestions))
	return prompt.FilterHasPrefix(suggestions, this.ctx.Word, true)
}

// o
func (this DaemonSetsSuggestion) ArgsOutput() []prompt.Suggest {
	return prompt.FilterHasPrefix([]prompt.Suggest{
		ModeDesc,
		ModeYAML,
		ModeJson,
	}, this.ctx.Word, false)
}

type DaemonSetsExecutor struct {
	ctx  *PromptCtx
	Tool *cmdResourceExecutorTool
}

func NewDaemonSetsExecutor(ctx *PromptCtx) *DaemonSetsExecutor {
	return &DaemonSetsExecutor{ctx: ctx, Tool: newCmdResourceExecutorTool()}
}

// o
func (this DaemonSetsExecutor) Output(ds *appsv1.DaemonSet, arg string) string {
	word := GetWordAfterArgWithSpace(this.ctx.Line, arg)
	ds.ManagedFields = []metav1.ManagedFieldsEntry{}
	return this.Tool.Output(ds.Name, ds.Namespace, word, "ds", ds)
}

func (this *DaemonSetsExecutor) Relationship(ds *appsv1.DaemonSet) string {
	var res string
	pr := NewDaemonSetRelationship(ds)
	res = pr.Tool.FetchAll(pr.generate())
	return res
}

// l
func (this DaemonSetsExecutor) Labels(ds *appsv1.DaemonSet) string {
	return this.Tool.Labels(ds.Labels)
}

// e
func (this DaemonSetsExecutor) Events(ds *appsv1.DaemonSet) string {
	return this.Tool.Events(ds.UID)
}

type DaemonSetRelationship struct {
	Tool *cmdResourceRelaTool
	ds   *appsv1.DaemonSet
	svc  []*v1.Service
	pods []*v1.Pod
}

func NewDaemonSetRelationship(ds *appsv1.DaemonSet) *DaemonSetRelationship {
	return &DaemonSetRelationship{ds: ds, Tool: newCmdResourceRelaTool()}
}

func (this *DaemonSetRelationship) Service() string {
	var relaObjs []*v1.Service
	objects := FetchServices(GlobalNamespace)
	for _, obj := range objects {
		if RuleJudgeLabelSelectorMatch(obj.Spec.Selector, this.ds.Labels) {
			relaObjs = append(relaObjs, obj)
		}
	}
	this.svc = relaObjs
	Logger(fmt.Sprintf("Relevant services count: %+v", len(relaObjs)))
	return this.Tool.FormatView(SliceResource2SliceRuntimeObj(relaObjs))
}

func (this *DaemonSetRelationship) Pod() string {
	var relaObjs []*v1.Pod
	objects := FetchPods(GlobalNamespace)
	for _, obj := range objects {
		if RuleJudgeLabelSelectorMatch(this.ds.Spec.Selector.MatchLabels, obj.Labels) {
			relaObjs = append(relaObjs, obj)
		}
	}
	this.pods = relaObjs
	Logger(fmt.Sprintf("Relevant pods count: %+v", len(this.pods)))
	return this.Tool.FormatView(SliceResource2SliceRuntimeObj(this.pods))
}

func (this *DaemonSetRelationship) generate() []func() string {
	return []func() string{
		this.Service,
		this.Pod,
	}
}

package query

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
)

func init() {
	RegisterCmdResources(NewService())
}

type Service struct {
	cmdResourceTool
	suggest *ServicesSuggestion
	exec    *ServiceExecutor
}

func NewService() *Service {
	return &Service{cmdResourceTool: cmdResourceTool{}}
}

func (this *Service) Name() string {
	return "services"
}

func (this *Service) ShortName() string {
	return "svc"
}

func (this *Service) CanExecute(ctx *PromptCtx) bool {
	if RuleHasPrefix(
		ctx.Line,
		this.Name(),
		this.ShortName()) {
		return true
	}
	return false
}

func (this *Service) Execute(ctx *PromptCtx) string {
	if this.CanExecute(ctx) {
		var svc *v1.Service
		this.exec = NewServiceExecutor(ctx)

		if RuleJudgeLineWordCount(ctx.Line, 2, Equal) {
			svc = FetchServiceWithName(ParserResourceName(this.FindResourceName(ctx)))
			if svc == nil {
				return ""
			}
			return this.exec.Output(svc, "")
		}

		if RuleJudgeLineWordCount(ctx.Line, 2, Greater) {

			svc = FetchServiceWithName(ParserResourceName(this.FindResourceName(ctx)))
			if svc == nil {
				return ""
			}
			// svc <ServiceName> -o [yaml|desc|json]
			if RuleJudgeLineHasWords(ctx.Line, ArgOutput.Text) {
				return this.exec.Output(svc, ArgOutput.Text)
			}
			// svc <ServiceName> -r
			if RuleJudgeLineHasWords(ctx.Line, ArgRelationship.Text) {
				return this.exec.Relationship(svc)
			}
			// svc <ServiceName> -l
			if RuleJudgeLineHasWords(ctx.Line, ArgLabel.Text) {
				return this.exec.Labels(svc)
			}
			if RuleJudgeLineHasWords(ctx.Line, ArgEvents.Text) {
				return this.exec.Events(svc)
			}

			if RuleJudgeLineHasWords(ctx.Line, ArgAnnotaions.Text) {
				return this.exec.Tool.Annotations(svc.Annotations)
			}
		}
	}
	return ""
}

func (this *Service) Suggestions(ctx *PromptCtx) []prompt.Suggest {
	Logger("Run svc suggestions.")
	if this.CanExecute(ctx) && RuleCanRemind(ctx.D) {
		this.suggest = NewServicesSuggestion(ctx)

		// svc <SvcName>
		if RuleJudgeLineWordCount(ctx.Line, 2, Equal) {
			return this.suggest.ArgSvcName()
		}

		if RuleJudgeLineWordCount(ctx.Line, 2, Greater) {
			// svc <SvcName> -o [yaml|desc|json]
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

type ServicesSuggestion struct {
	ctx *PromptCtx
}

func NewServicesSuggestion(ctx *PromptCtx) *ServicesSuggestion {
	return &ServicesSuggestion{ctx: ctx}
}

func (this ServicesSuggestion) Helper() []prompt.Suggest {
	return []prompt.Suggest{
		ArgOutput,
		ArgLabel,
		ArgRelationship,
		ArgEvents,
		ArgAnnotaions,
	}
}

func (this ServicesSuggestion) ArgSvcName() []prompt.Suggest {
	svcs := FetchServices(GlobalNamespace)
	suggestions := make([]prompt.Suggest, 0)
	for _, svc := range svcs {
		suggestions = append(suggestions, prompt.Suggest{
			Text: FormatResourceName(svc.Name, svc.Namespace),
			Description: fmt.Sprintf(
				"Type: %s, ClusterIP: %s, ExternalIP: %s, Ports: %s, Selector: %s",
				svc.Spec.Type,
				svc.Spec.ClusterIP,
				SliceString2String(svc.Spec.ExternalIPs),
				ParserServicePorts2String(svc.Spec.Ports),
				Map2String(svc.Spec.Selector),
			),
		})
	}
	sort.Sort(SuggestSlice(suggestions))
	return prompt.FilterHasPrefix(suggestions, this.ctx.Word, true)
}

// o
func (this ServicesSuggestion) ArgsOutput() []prompt.Suggest {
	return prompt.FilterHasPrefix([]prompt.Suggest{
		ModeDesc,
		ModeYAML,
		ModeJson,
	}, this.ctx.Word, false)
}

type ServiceExecutor struct {
	ctx  *PromptCtx
	Tool *cmdResourceExecutorTool
}

func NewServiceExecutor(ctx *PromptCtx) *ServiceExecutor {
	return &ServiceExecutor{ctx: ctx, Tool: newCmdResourceExecutorTool()}
}

// o
func (this ServiceExecutor) Output(svc *v1.Service, arg string) string {
	word := GetWordAfterArgWithSpace(this.ctx.Line, arg)
	svc.ManagedFields = []metav1.ManagedFieldsEntry{}
	return this.Tool.Output(svc.Name, svc.Namespace, word, "svc", svc)
}

// r, pods
func (this ServiceExecutor) Relationship(svc *v1.Service) string {
	var res string
	pr := NewServiceRelationship(svc)
	res = pr.Tool.FetchAll(pr.generate())
	return res
}

// l
func (this ServiceExecutor) Labels(svc *v1.Service) string {
	return this.Tool.Labels(svc.Labels)
}

// e
func (this ServiceExecutor) Events(svc *v1.Service) string {
	return this.Tool.Events(svc.UID)
}

type ServiceRelationship struct {
	Tool *cmdResourceRelaTool
	pods []*v1.Pod
	svc  *v1.Service
}

func NewServiceRelationship(svc *v1.Service) *ServiceRelationship {
	return &ServiceRelationship{svc: svc, Tool: newCmdResourceRelaTool()}
}

func (this *ServiceRelationship) Pods() string {
	this.pods = FetchPods(GlobalNamespace, this.svc.Spec.Selector)
	return this.Tool.FormatView(SliceResource2SliceRuntimeObj(this.pods))
}

func (this *ServiceRelationship) generate() []func() string {
	return []func() string{
		this.Pods,
	}
}

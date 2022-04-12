package query

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
	"sort"
	"strings"
)

func init() {
	RegisterCmdResources(NewDeploy())
}

type Deploy struct {
	cmdResourceTool
	suggest *DeploySuggestion
	exec    *DeployExecutor
}

func NewDeploy() *Deploy {
	return &Deploy{cmdResourceTool: cmdResourceTool{}}
}

func (this *Deploy) Name() string {
	return "deployments"
}

func (this *Deploy) ShortName() string {
	return "deploy"
}

func (this *Deploy) CanExecute(ctx *PromptCtx) bool {
	Logger("Line: ", ctx.Line)
	if RuleHasPrefix(
		ctx.Line,
		this.Name(),
		this.ShortName()) {
		return true
	}
	return false
}

func (this *Deploy) Execute(ctx *PromptCtx) string {
	if this.CanExecute(ctx) {
		var deploy *appv1.Deployment
		this.exec = NewDeployExecutor(ctx)
		// pods <PodName>
		if RuleJudgeLineWordCount(ctx.Line, 2, Equal) {
			deploy = FetchDeploymentWithName(ParserResourceName(this.FindResourceName(ctx)))
			if deploy == nil {
				return ""
			}
			return this.exec.Output(deploy, "")
		}

		if RuleJudgeLineWordCount(ctx.Line, 2, Greater) {

			deploy = FetchDeploymentWithName(ParserResourceName(this.FindResourceName(ctx)))
			if deploy == nil {
				return ""
			}

			// deploy <DeployName> -o [yaml|desc|json]
			if RuleJudgeLineHasWords(ctx.Line, ArgOutput.Text) {
				return this.exec.Output(deploy, ArgOutput.Text)
			}
			// deploy <DeployName> -r
			if RuleJudgeLineHasWords(ctx.Line, ArgRelationship.Text) {
				return this.exec.Relationship(deploy)
			}
			// deploy <DeployName> -l
			if RuleJudgeLineHasWords(ctx.Line, ArgLabel.Text) {
				return this.exec.Label(deploy)
			}
			// pods <PodName> -e
			if RuleJudgeLineHasWords(ctx.Line, ArgEvents.Text) {
				return this.exec.Event(deploy)
			}
			// pods <PodName> -v
			if RuleJudgeLineHasWords(ctx.Line, ArgVolumes.Text) {
				return this.exec.Volume(deploy)
			}
		}
	}

	return ""
}

func (this *Deploy) Suggestions(ctx *PromptCtx) []prompt.Suggest {
	Logger("Run pods suggestions.")

	if this.CanExecute(ctx) && RuleCanRemind(ctx.D) {
		this.suggest = NewDeploySuggestion(ctx)

		// deploy <DeployName>
		if RuleJudgeLineWordCount(ctx.Line, 2, Equal) {
			return this.suggest.ArgDeployName()
		}

		if RuleJudgeLineWordCount(ctx.Line, 2, Greater) {
			// deploy <DeployName> -o [yaml|desc|json]
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

type DeploySuggestion struct {
	ctx *PromptCtx
}

func NewDeploySuggestion(ctx *PromptCtx) *DeploySuggestion {
	return &DeploySuggestion{ctx: ctx}
}

func (this DeploySuggestion) Helper() []prompt.Suggest {
	return []prompt.Suggest{
		ArgOutput,
		ArgRelationship,
		ArgEvents,
		ArgVolumes,
		ArgLabel,
	}
}

func (this DeploySuggestion) ArgDeployName() []prompt.Suggest {
	deploys := FetchDeployments(GlobalNamespace)
	suggestions := make([]prompt.Suggest, 0)
	for _, d := range deploys {
		suggestions = append(suggestions, prompt.Suggest{
			Text: FormatResourceName(d.Name, d.Namespace),
			Description: fmt.Sprintf("Ready: %v/%v, Selector: %s",
				d.Status.ReadyReplicas,
				d.Status.Replicas,
				Map2String(d.Spec.Selector.MatchLabels),
			),
		})
	}
	sort.Sort(SuggestSlice(suggestions))
	return suggestions
}

// o
func (this DeploySuggestion) ArgsOutput() []prompt.Suggest {
	return prompt.FilterHasPrefix([]prompt.Suggest{
		ModeDesc,
		ModeYAML,
		ModeJson,
	}, this.ctx.Word, false)
}

type DeployExecutor struct {
	ctx  *PromptCtx
	Tool *cmdResourceExecutorTool
}

func NewDeployExecutor(ctx *PromptCtx) *DeployExecutor {
	return &DeployExecutor{ctx: ctx, Tool: newCmdResourceExecutorTool()}
}

func (this *DeployExecutor) Output(deploy *appv1.Deployment, arg string) string {
	word := GetWordAfterArgWithSpace(this.ctx.Line, arg)
	deploy.ManagedFields = []metav1.ManagedFieldsEntry{}
	return this.Tool.Output(deploy.Name, word, "deploy", deploy)
}

func (this *DeployExecutor) Relationship(deploy *appv1.Deployment) string {
	var res string
	pr := NewDeployRelationship(deploy)
	res = pr.Tool.FetchAll(pr.generate())
	return res
}

func (this *DeployExecutor) Event(deploy *appv1.Deployment) string {
	return this.Tool.Events(deploy.UID)
}

func (this *DeployExecutor) Label(deploy *appv1.Deployment) string {
	return this.Tool.Labels(deploy.Labels)
}

func (this *DeployExecutor) Volume(deploy *appv1.Deployment) string {
	volumes := []string{
		"##### Volumes #####",
	}
	bVolumes, err := yaml.Marshal(deploy.Spec.Template.Spec.Volumes)
	WrapError(err)
	volumes = append(volumes, string(bVolumes))
	return strings.Join(volumes, "\n")
}

type DeployRelationship struct {
	Tool   *cmdResourceRelaTool
	deploy *appv1.Deployment
	svc    []*v1.Service
	pod    *v1.Pod
	rs     *appv1.ReplicaSet
}

func NewDeployRelationship(deploy *appv1.Deployment) *DeployRelationship {
	return &DeployRelationship{deploy: deploy, Tool: newCmdResourceRelaTool()}
}

func (this *DeployRelationship) Service() string {
	var relaObjs []*v1.Service
	objects := FetchServices(GlobalNamespace)
	for _, obj := range objects {
		if RuleJudgeLabelSelectorMatch(obj.Spec.Selector, this.deploy.Labels) {
			relaObjs = append(relaObjs, obj)
		}
	}
	this.svc = relaObjs
	Logger(fmt.Sprintf("Relevant services count: %+v", len(relaObjs)))
	return this.Tool.FormatView(SliceResource2SliceRuntimeObj(relaObjs))
}

func (this *DeployRelationship) ReplicaSet() string {
	var res string
	replicaSets := FetchReplicaSets(GlobalNamespace)
	dRevision := FetchAnnotationsValue(this.deploy.Annotations,
		"deployment.kubernetes.io/revision")
	for _, rs := range replicaSets {
		for _, owner := range rs.OwnerReferences {
			if owner.Kind == KindDeployment && owner.Name == this.deploy.Name {
				revision := FetchAnnotationsValue(rs.Annotations,
					"deployment.kubernetes.io/revision")

				if revision != Empty && dRevision == revision {
					this.rs = rs
					res = this.Tool.FormatView([]runtime.Object{rs})
					return res
				}
			}
		}
	}
	return res
}

func (this *DeployRelationship) Pod() string {
	var res string
	if this.rs != nil {
		pods := FetchPods(GlobalNamespace, this.rs.Labels)
		res = this.Tool.FormatView(SliceResource2SliceRuntimeObj(pods))
	}
	return res
}

func (this *DeployRelationship) generate() []func() string {
	return []func() string{
		this.Service,
		this.ReplicaSet,
		this.Pod,
	}
}

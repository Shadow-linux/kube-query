package query

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
)

func init() {
	RegisterCmdResources(NewNodes())
}

type Nodes struct {
	cmdResourceTool
	suggest *NodesSuggestion
	exec    *NodesExecutor
}

func NewNodes() *Nodes {
	return &Nodes{cmdResourceTool: cmdResourceTool{}}
}

func (this *Nodes) Name() string {
	return "nodes"
}

func (this *Nodes) ShortName() string {
	return "no"
}

func (this *Nodes) CanExecute(ctx *PromptCtx) bool {
	if RuleHasPrefix(
		ctx.Line,
		this.Name(),
		this.ShortName()) {
		return true
	}
	return false
}

func (this *Nodes) Execute(ctx *PromptCtx) string {
	if this.CanExecute(ctx) {
		var node *v1.Node
		this.exec = NewNodesExecutor(ctx)

		if RuleJudgeLineWordCount(ctx.Line, 2, Equal) {
			node = FetchNodeWithName(this.FindResourceName(ctx))
			if node == nil {
				return ""
			}
			return this.exec.Output(node, "")
		}

		if RuleJudgeLineWordCount(ctx.Line, 2, Greater) {

			node = FetchNodeWithName(this.FindResourceName(ctx))
			if node == nil {
				return ""
			}
			// cm <ServiceName> -o [yaml|desc|json]
			if RuleJudgeLineHasWords(ctx.Line, ArgOutput.Text) {
				return this.exec.Output(node, ArgOutput.Text)
			}
			// cm <ServiceName> -l
			if RuleJudgeLineHasWords(ctx.Line, ArgLabel.Text) {
				return this.exec.Labels(node)
			}
			// pods <PodName> -e
			if RuleJudgeLineHasWords(ctx.Line, ArgEvents.Text) {
				return this.exec.Events(node)
			}

			// pods <PodName> --anno
			if RuleJudgeLineHasWords(ctx.Line, ArgAnnotaions.Text) {
				return this.exec.Annotations(node)
			}
		}
	}
	return ""
}

func (this *Nodes) Suggestions(ctx *PromptCtx) []prompt.Suggest {
	Logger("Run nodes suggestions.")
	if this.CanExecute(ctx) && RuleCanRemind(ctx.D) {
		this.suggest = NewNodesSuggestion(ctx)

		// node <nodeName>
		if RuleJudgeLineWordCount(ctx.Line, 2, Equal) {
			return this.suggest.ArgCmName()
		}

		if RuleJudgeLineWordCount(ctx.Line, 2, Greater) {
			// node <NodeNName> -o [yaml|desc|json]
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

type NodesSuggestion struct {
	ctx *PromptCtx
}

func NewNodesSuggestion(ctx *PromptCtx) *NodesSuggestion {
	return &NodesSuggestion{ctx: ctx}
}

func (this NodesSuggestion) Helper() []prompt.Suggest {
	return prompt.FilterHasPrefix([]prompt.Suggest{
		ArgOutput,
		ArgEvents,
		ArgLabel,
		ArgAnnotaions,
	}, this.ctx.Word, false)
}

func (this NodesSuggestion) ArgCmName() []prompt.Suggest {
	nodes := FetchNodes()
	suggestions := make([]prompt.Suggest, 0)
	for _, no := range nodes {
		suggestions = append(suggestions, prompt.Suggest{
			Text: no.Name,
			Description: fmt.Sprintf(
				"Name: %s, Status: %s, InternalIP: %s",
				no.Name,
				no.Status.Conditions,
				no.Status.Addresses[0].Address,
			),
		})
	}
	sort.Sort(SuggestSlice(suggestions))
	return prompt.FilterHasPrefix(suggestions, this.ctx.Word, true)
}

// o
func (this NodesSuggestion) ArgsOutput() []prompt.Suggest {
	return []prompt.Suggest{
		ModeDesc,
		ModeYAML,
		ModeJson,
	}
}

type NodesExecutor struct {
	ctx  *PromptCtx
	Tool *cmdResourceExecutorTool
}

func NewNodesExecutor(ctx *PromptCtx) *NodesExecutor {
	return &NodesExecutor{ctx: ctx, Tool: newCmdResourceExecutorTool()}
}

// o
func (this NodesExecutor) Output(no *v1.Node, arg string) string {
	word := GetWordAfterArgWithSpace(this.ctx.Line, arg)
	no.ManagedFields = []metav1.ManagedFieldsEntry{}
	return this.Tool.Output(no.Name, word, "nodes", no)
}

// l
func (this NodesExecutor) Labels(no *v1.Node) string {
	return this.Tool.Labels(no.Labels)
}

// e
func (this NodesExecutor) Events(no *v1.Node) string {
	return this.Tool.Events(no.UID)
}

// e
func (this NodesExecutor) Annotations(no *v1.Node) string {
	return this.Tool.Annotations(no.Annotations)
}

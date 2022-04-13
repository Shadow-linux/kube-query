package query

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
)

func init() {
	RegisterCmdResources(NewConfigMap())
}

type ConfigMap struct {
	cmdResourceTool
	suggest *ConfigMapSuggestion
	exec    *ConfigMapExecutor
}

func NewConfigMap() *ConfigMap {
	return &ConfigMap{cmdResourceTool: cmdResourceTool{}}
}

func (this *ConfigMap) Name() string {
	return "configmaps"
}

func (this *ConfigMap) ShortName() string {
	return "cm"
}

func (this *ConfigMap) CanExecute(ctx *PromptCtx) bool {
	if RuleHasPrefix(
		ctx.Line,
		this.Name(),
		this.ShortName()) {
		return true
	}
	return false
}

func (this *ConfigMap) Execute(ctx *PromptCtx) string {
	if this.CanExecute(ctx) {
		var cm *v1.ConfigMap
		this.exec = NewConfigMapExecutor(ctx)

		if RuleJudgeLineWordCount(ctx.Line, 2, Equal) {
			cm = FetchConfigMapWithName(ParserResourceName(this.FindResourceName(ctx)))
			if cm == nil {
				return ""
			}
			return this.exec.Output(cm, "")
		}

		if RuleJudgeLineWordCount(ctx.Line, 2, Greater) {

			cm = FetchConfigMapWithName(ParserResourceName(this.FindResourceName(ctx)))
			if cm == nil {
				return ""
			}
			// cm <ServiceName> -o [yaml|desc|json]
			if RuleJudgeLineHasWords(ctx.Line, ArgOutput.Text) {
				return this.exec.Output(cm, ArgOutput.Text)
			}
			// cm <ServiceName> -l
			if RuleJudgeLineHasWords(ctx.Line, ArgLabel.Text) {
				return this.exec.Labels(cm)
			}
			// cm <ServiceName> -e
			if RuleJudgeLineHasWords(ctx.Line, ArgEvents.Text) {
				return this.exec.Events(cm)
			}
			// cm <ServiceName> --anno
			if RuleJudgeLineHasWords(ctx.Line, ArgAnnotaions.Text) {
				return this.exec.Tool.Annotations(cm.Annotations)
			}
		}
	}
	return ""
}

func (this *ConfigMap) Suggestions(ctx *PromptCtx) []prompt.Suggest {
	Logger("Run cm suggestions.")
	if this.CanExecute(ctx) && RuleCanRemind(ctx.D) {
		this.suggest = NewConfigMapSuggestion(ctx)

		// cm <cmName>
		if RuleJudgeLineWordCount(ctx.Line, 2, Equal) {
			return this.suggest.ArgCmName()
		}

		if RuleJudgeLineWordCount(ctx.Line, 2, Greater) {
			// cm <cmName> -o [yaml|desc|json]
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

type ConfigMapSuggestion struct {
	ctx *PromptCtx
}

func NewConfigMapSuggestion(ctx *PromptCtx) *ConfigMapSuggestion {
	return &ConfigMapSuggestion{ctx: ctx}
}

func (this ConfigMapSuggestion) Helper() []prompt.Suggest {
	return prompt.FilterHasPrefix([]prompt.Suggest{
		ArgOutput,
		ArgEvents,
		ArgLabel,
		ArgAnnotaions,
	}, this.ctx.Word, false)
}

func (this ConfigMapSuggestion) ArgCmName() []prompt.Suggest {
	cms := FetchConfigMaps(GlobalNamespace)
	suggestions := make([]prompt.Suggest, 0)
	for _, cm := range cms {
		suggestions = append(suggestions, prompt.Suggest{
			Text: FormatResourceName(cm.Name, cm.Namespace),
			Description: fmt.Sprintf(
				"Name: %s, Data: %d,",
				cm.Name,
				len(cm.Data),
			),
		})
	}
	sort.Sort(SuggestSlice(suggestions))
	return prompt.FilterHasPrefix(suggestions, this.ctx.Word, true)
}

// o
func (this ConfigMapSuggestion) ArgsOutput() []prompt.Suggest {
	return []prompt.Suggest{
		ModeDesc,
		ModeYAML,
		ModeJson,
	}
}

type ConfigMapExecutor struct {
	ctx  *PromptCtx
	Tool *cmdResourceExecutorTool
}

func NewConfigMapExecutor(ctx *PromptCtx) *ConfigMapExecutor {
	return &ConfigMapExecutor{ctx: ctx, Tool: newCmdResourceExecutorTool()}
}

// o
func (this ConfigMapExecutor) Output(cm *v1.ConfigMap, arg string) string {
	word := GetWordAfterArgWithSpace(this.ctx.Line, arg)
	cm.ManagedFields = []metav1.ManagedFieldsEntry{}
	return this.Tool.Output(cm.Name, word, "cm", cm)
}

// l
func (this ConfigMapExecutor) Labels(cm *v1.ConfigMap) string {
	return this.Tool.Labels(cm.Labels)
}

// e
func (this ConfigMapExecutor) Events(cm *v1.ConfigMap) string {
	return this.Tool.Events(cm.UID)
}

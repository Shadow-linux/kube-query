package query

import (
	"github.com/c-bata/go-prompt"
	"k8s.io/utils/strings/slices"
	"sort"
)

func init() {
	RegisterCmdResources(NewNamespaces())
}

type Namespaces struct {
	cmdResourceTool
	suggest *NamespacesSuggestion
}

func NewNamespaces() *Namespaces {
	return &Namespaces{cmdResourceTool: cmdResourceTool{}}
}

func (this *Namespaces) Name() string {
	return "use"
}

func (this *Namespaces) ShortName() string {
	return "u"
}

func (this *Namespaces) CanExecute(ctx *PromptCtx) bool {
	if RuleHasPrefix(
		ctx.Line,
		this.Name(),
		this.ShortName()) {
		return true
	}
	return false
}

func (this *Namespaces) Execute(ctx *PromptCtx) string {
	if this.CanExecute(ctx) {
		if RuleJudgeLineWordCount(ctx.Line, 2, Greater) {
			return MsgExpectLineWord(2)
		}
		newNamespace := FormatLineWithSpace(ctx.Line)[1]
		var namespaces []string
		for _, ns := range FetchAllNamespace() {
			namespaces = append(namespaces, ns.Name)
		}
		if slices.Contains(namespaces, newNamespace) || newNamespace == AllNamespace {
			GlobalNamespace = newNamespace
			return "Set namespace " + GlobalNamespace

		}

		return "Not found namespace " + GlobalNamespace

	}
	return ""

}

func (this *Namespaces) Suggestions(ctx *PromptCtx) []prompt.Suggest {
	Logger("Run namespace suggestions.")

	if this.CanExecute(ctx) && RuleCanRemind(ctx.D) {
		this.suggest = NewNamespacesSuggestion(ctx)

		// use <Namespace>
		if RuleJudgeLineWordCount(ctx.Line, 2, Equal) {
			return this.suggest.Helper()
		}

	}
	return []prompt.Suggest{}
}

type NamespacesSuggestion struct {
	ctx *PromptCtx
}

func NewNamespacesSuggestion(ctx *PromptCtx) *NamespacesSuggestion {
	return &NamespacesSuggestion{ctx: ctx}
}

func (this *NamespacesSuggestion) Helper() []prompt.Suggest {
	namespaces := FetchAllNamespace()
	suggestions := []prompt.Suggest{prompt.Suggest{
		Text:        "all",
		Description: "All namespace",
	}}
	for _, n := range namespaces {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        n.Name,
			Description: "Namespace",
		})
	}
	sort.Sort(SuggestSlice(suggestions))
	return prompt.FilterHasPrefix(suggestions, this.ctx.Word, true)
}

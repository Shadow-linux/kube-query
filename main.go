package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Shadow-linux/kube-query/kube"
	"github.com/Shadow-linux/kube-query/query"
	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
	"os"
	"strings"
)

var (
	Version  string = "0.0.1"
	Revision string = ""
)

func Executor(line string) {
	line = strings.TrimSpace(line)
	query.Logger("Execute line: ", line)
	switch {
	case strings.TrimSpace(line) == "":
		return
	case line == "quit" || line == "exit":
		query.Print("Bye !")
		os.Exit(0)
		return
	case line == "clear":
		query.WrapError(query.ClearConsole())
		return
	case line == "help":
		query.Print(query.CmdHelpInfo)
		return
	case strings.HasPrefix(line, "@"):
		query.Logger("Shell mode.")
		query.PureCmdRun(line[1:], true)
		return
	}

	// kube-query
	// fetch line before pipe ('|')
	fmtLine := query.GetLineBeforePipe(line)
	query.Logger("FmtExecuteLine: ", fmtLine)
	ctx := query.NewPromptCtx(context.Background(), fmtLine, "", prompt.Document{Text: line})

	for _, resource := range query.ResourcesList {
		if resource.CanExecute(ctx) {
			if strings.Contains(line, "|") {
				info := resource.Execute(ctx)
				sepLines := strings.Split(line, "|")
				cmd := fmt.Sprintf("echo \"%s\" | %s ", info, strings.Join(sepLines[1:], "|"))
				fmt.Println(query.CmdRunWithFile(cmd, false))
				return
			}
			fmt.Println(resource.Execute(ctx))
			return
		}
	}

	// kube-prompt
	query.K8sCmdRun(line, true)
	return
}

func Completer(d prompt.Document) []prompt.Suggest {

	if d.TextBeforeCursor() == query.Empty {
		return []prompt.Suggest{}
	}
	query.Logger("Complete line: ", d.Text)
	w := d.GetWordBeforeCursor()
	query.Logger("Word: ", w)

	// kube query completer
	var suggestions = make([]prompt.Suggest, 0)
	for _, resource := range query.ResourcesList {
		fmtLine := query.GetLineBeforePipe(d.Text)
		ctx := query.NewPromptCtx(context.Background(), fmtLine, w, d)
		if query.RuleJudgeLineWordCount(fmtLine, 1, query.Equal) && !query.RuleJudgeLineHasSpace(fmtLine) {
			suggestions = resource.DefaultSuggestions(ctx, resource.Name(), resource.ShortName())
		} else {
			suggestions = resource.Suggestions(ctx)
		}

		if len(suggestions) != 0 {
			break
		}
	}

	// kube prompt completer
	if query.RuleJudgeLineWordCount(d.Text, 1, query.Equal) || len(suggestions) == 0 {
		c, err := kube.NewCompleter()
		query.WrapError(err)
		suggestions = append(suggestions, c.Complete(d)...)
	}

	// file path completer
	pathSuggestions := query.FilePathComplete(d.GetWordBeforeCursor())
	suggestions = append(suggestions, pathSuggestions...)

	query.Logger("suggestions", suggestions)
	return suggestions

}

func main() {
	defer func() {
		query.Print("Bye !")
		if e := recover(); e != nil {
			query.Logger("[Error]: ", e)
		}
	}()

	debug := flag.Bool("debug", false, "Turn on debug mode.")
	version := flag.Bool("version", false, "Print version.")
	flag.Parse()
	query.Debug = *debug
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	//init
	query.InitClient()
	query.InitInformerCache()

	// prompt
	fmt.Printf("kube-query %s (rev-%s)\n", Version, Revision)
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")
	query.ConsoleStdoutWriter = prompt.NewStdoutWriter()
	p := prompt.New(
		Executor,
		Completer,
		prompt.OptionTitle("kube-prompt: interactive kubernetes client"),
		prompt.OptionPrefix(">>> "),
		prompt.OptionInputTextColor(prompt.White),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
		prompt.OptionWriter(query.ConsoleStdoutWriter),
		prompt.OptionShowCompletionAtStart(),
	)
	p.Run()
}

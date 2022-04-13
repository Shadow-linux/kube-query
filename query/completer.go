package query

import (
	"github.com/c-bata/go-prompt"
	"path"
	"strings"
)

var (
	ArgOutput         = prompt.Suggest{Text: "-o", Description: "output one of mode: [yaml|desc|json]. default: desc."}
	ArgRelationship   = prompt.Suggest{Text: "-r", Description: "relationship, display resource label relationship."}
	ArgLabel          = prompt.Suggest{Text: "-l", Description: "label, display labels."}
	ArgServiceAccount = prompt.Suggest{Text: "-a", Description: "account (service), display service account."}
	ArgVolumes        = prompt.Suggest{Text: "-v", Description: "volumes, display volumes & volume mounts."}
	ArgEvents         = prompt.Suggest{Text: "-e", Description: "events, display events."}
	ArgInteractive    = prompt.Suggest{Text: "-i", Description: "interactive, interactive with container. default: the first container."}
	ArgShell          = prompt.Suggest{Text: "-s", Description: "shell command, specify shell command like: sh, /bin/bash. [default: sh]"}

	// output mode
	ModeYAML = prompt.Suggest{Text: "yaml", Description: "Output for yaml mode."}
	ModeDesc = prompt.Suggest{Text: "desc", Description: "Output for desc mode."}
	ModeJson = prompt.Suggest{Text: "json", Description: "Output for json mode."}
)

func FilePathComplete(p string) []prompt.Suggest {
	var res []prompt.Suggest
	if strings.HasPrefix(p, "/") || strings.HasPrefix(p, "./") {
		fileDirs := CurrentTierPath(path.Dir(p))
		for _, fd := range fileDirs {
			res = append(res, prompt.Suggest{
				Text:        fd,
				Description: "path",
			})
		}
		words := strings.Split(p, "/")
		w := words[len(words)-1]
		return prompt.FilterContains(res, w, true)
	}
	return res
}

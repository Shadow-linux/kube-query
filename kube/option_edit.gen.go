package kube

import (
	prompt "github.com/c-bata/go-prompt"
)

var editOptions = []prompt.Suggest{
	prompt.Suggest{Text: "-f", Description: "Filename, directory, or URL to files to use to edit the resource"},
	prompt.Suggest{Text: "--filename", Description: "Filename, directory, or URL to files to use to edit the resource"},
	prompt.Suggest{Text: "--include-extended-apis", Description: "If true, include definitions of new APIs via calls to the API server. [default true]"},
	prompt.Suggest{Text: "--include-uninitialized", Description: "If true, the kubectl command applies to uninitialized objects. If explicitly set to false, this flag overrides other flags that make the kubectl commands apply to uninitialized objects, e.g., \"--all\". Objects with empty metadata.initializers are regarded as initialized."},
	prompt.Suggest{Text: "-o", Description: "Output format. One of: yaml|json."},
	prompt.Suggest{Text: "--output", Description: "Output format. One of: yaml|json."},
	prompt.Suggest{Text: "--output-patch", Description: "Output the patch if the resource is edited."},
	prompt.Suggest{Text: "--record", Description: "Record current kubectl command in the resource annotation. If set to false, do not record the command. If set to true, record the command. If not set, default to updating the existing annotation value only if one already exists."},
	prompt.Suggest{Text: "-R", Description: "Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory."},
	prompt.Suggest{Text: "--recursive", Description: "Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory."},
	prompt.Suggest{Text: "--save-config", Description: "If true, the configuration of current object will be saved in its annotation. Otherwise, the annotation will be unchanged. This flag is useful when you want to perform kubectl apply on this object in the future."},
	prompt.Suggest{Text: "--validate", Description: "If true, use a schema to validate the input before sending it"},
	prompt.Suggest{Text: "--windows-line-endings", Description: "Defaults to the line ending native to your platform."},
}

package query

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/strings/slices"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
)

func Print(contents ...interface{}) {
	l := ""
	for _, c := range contents {
		l = fmt.Sprintf("%+v %+v", l, c)
	}
	fmt.Println(l)
}

func Logger(contents ...interface{}) {
	if !Debug {
		return
	}

	l := ""
	for _, c := range contents {
		l = fmt.Sprintf("%+v %+v", l, c)
	}
	log.Println(l)
}

func WrapError(err error) {
	if err != nil {
		log.Println(err.Error())
		if Debug {
			panic(err.Error())
		}
	}
}

func ClearConsole() error {
	ConsoleStdoutWriter.EraseScreen()
	ConsoleStdoutWriter.CursorGoTo(0, 0)
	return ConsoleStdoutWriter.Flush()
}

func NewTable(writer io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(writer)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	return table
}

func FormatLineWithSpace(line string) []string {
	return strings.Split(regexp.MustCompile("\\s+").ReplaceAllString(line, Space), Space)
}

func FetchFirstArg(line string) string {
	lines := FormatLineWithSpace(line)
	if len(lines) >= 2 {
		return lines[1]
	}
	return ""
}

func Map2String(m map[string]string) string {
	var res string
	for k, v := range m {
		if res == "" {
			res = fmt.Sprintf("%s=%v", k, v)
			continue
		}
		res = res + ";" + fmt.Sprintf("%s=%v", k, v)
	}
	return res
}

func Map2Slice(m map[string]string) []string {
	var res []string
	for k, v := range m {
		res = append(res, fmt.Sprintf("%s=%v", k, v))
	}
	return res
}

func SliceString2String(s []string) string {
	var res string
	for _, v := range s {
		if res == "" {
			return v
		}
		res = res + ";" + v
	}
	return res
}

func SliceResource2SliceRuntimeObj(resources interface{}) []runtime.Object {
	var res []runtime.Object
	if reflect.TypeOf(resources).Kind() == reflect.Slice {
		s := reflect.ValueOf(resources)
		for i := 0; i < s.Len(); i++ {
			elem := s.Index(i)
			res = append(res, elem.Interface().(runtime.Object))
		}
	}
	return res
}

func SliceResource2SliceMetav1Obj(resources interface{}) []metav1.Object {
	var res []metav1.Object
	if reflect.TypeOf(resources).Kind() == reflect.Slice {
		s := reflect.ValueOf(resources)
		for i := 0; i < s.Len(); i++ {
			elem := s.Index(i)
			res = append(res, elem.Interface().(metav1.Object))
		}
	}
	return res
}

func FormatResourceName(name, ns string) string {
	return name + "." + ns
}

func ParserResourceName(fmtName string) (name, ns string) {
	if !strings.Contains(fmtName, ".") {
		return fmtName, ""
	}
	s := strings.Split(fmtName, ".")
	name, ns = s[0], s[1]
	return
}

func ParserServicePorts2String(ports []v1.ServicePort) string {
	var res string
	for i, p := range ports {
		var tmp string
		if p.NodePort != 0 {
			tmp = fmt.Sprintf("%d:%d/%s", p.Port, p.NodePort, p.Protocol)
		} else {
			tmp = fmt.Sprintf("%d/%s", p.Port, p.Protocol)
		}
		if i == 0 {
			res = tmp
			continue
		}
		res = res + ";" + tmp
	}
	return res
}

func GetWordAfterArgWithSpace(line string, arg string) string {
	lineSlice := FormatLineWithSpace(line)
	idx := slices.Index(lineSlice, arg)
	length := len(lineSlice)
	if idx >= 0 {
		if idx+1 <= length-1 {
			return lineSlice[idx+1]
		}
	}
	return ""
}

func GetLineBeforePipe(line string) string {
	if strings.Contains(line, "|") {
		sepLines := strings.Split(line, "|")
		return strings.TrimSpace(sepLines[0])
	}
	return line
}

func GetLineAfterPipe(line string) string {
	if strings.Contains(line, "|") {
		sepLines := strings.Split(line, "|")
		return strings.TrimSpace(strings.Join(sepLines[1:], "|"))
	}
	return line
}

func Int2String(d interface{}) string {
	return fmt.Sprintf("%d", d)
}

func FetchAnnotationsValue(annos map[string]string, targetKey string) string {
	for k, v := range annos {
		if k == targetKey {
			return v
		}
	}
	return ""
}

func CurrentTierPath(fPath string) []string {
	var res []string
	// get file state
	if fState, err := os.Stat(fPath); err == nil {
		if !fState.IsDir() {
			res = append(res, fPath)
		} else {
			f, _ := os.Open(fPath)
			defer f.Close()
			names, _ := f.Readdirnames(0)
			for _, name := range names {
				res = append(res, name)
			}
		}
	}
	return res
}

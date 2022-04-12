package query

import (
	"github.com/c-bata/go-prompt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// sort
type ResourceObjSlice []metav1.Object

func (this ResourceObjSlice) Len() int {
	return len(this)
}

func (this ResourceObjSlice) Less(i, j int) bool {
	return this[j].GetCreationTimestamp().Time.After(this[i].GetCreationTimestamp().Time)
}

func (this ResourceObjSlice) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

type StringSlice []string

func (this StringSlice) Len() int {
	return len(this)
}

func (this StringSlice) Less(i, j int) bool {
	return this[i] < this[j]
}

func (this StringSlice) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

type SuggestSlice []prompt.Suggest

func (this SuggestSlice) Len() int {
	return len(this)
}

func (this SuggestSlice) Less(i, j int) bool {
	return this[i].Text < this[j].Text
}

func (this SuggestSlice) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

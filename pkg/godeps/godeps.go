package godeps

import (
	"strings"
)

type Dep struct {
	ImportPath string `json:"ImportPath"`
	Comment    string `json:"Comment,omitempty"`
	Rev        string `json:"Rev"`
}

type Godeps struct {
	ImportPath   string   `json:"ImportPath"`
	GoVersion    string   `json:"GoVersion"`
	GodepVersion string   `json:"GodepVersion"`
	Packages     []string `json:"Packages"`
	Deps         []Dep    `json:"Deps"`
}

type DepDiff map[string]Dep

type GoDepsDiff struct {
	Old, New   *Godeps
	ImportPath map[string]DepDiff
}

type SortByImportPath []Dep
func (a SortByImportPath) Len() int           { return len(a) }
func (a SortByImportPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByImportPath) Less(i, j int) bool { return strings.Compare(a[i].ImportPath, a[j].ImportPath) == -1 }

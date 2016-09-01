package godeps

import (
	"fmt"
	"os"
	"sort"
	"syscall"

	"github.com/legionus/godep-utils/pkg/utils"
)

const (
	DepsHookChooseOldExitStatus = 127 + 1
	DepsHookChooseNewExitStatus = 127 + 2
	DepsHookChooseIgnExitStatus = 127 + 3
)

type MergeHooks struct {
	PreHook, DepHook, PostHook string
}

func runHook(hook string, env []string) (int, error) {
	if hook == "" {
		return 0, nil
	}

	env = append(env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
	env = append(env, fmt.Sprintf("HOME=%s", os.Getenv("HOME")))
	env = append(env, fmt.Sprintf("TMPDIR=%s", os.Getenv("TMPDIR")))
	env = append(env, fmt.Sprintf("GOPATH=%s", os.Getenv("GOPATH")))
	env = append(env, fmt.Sprintf("EXIT_CHOOSE_OLD=%d", DepsHookChooseOldExitStatus))
	env = append(env, fmt.Sprintf("EXIT_CHOOSE_NEW=%d", DepsHookChooseNewExitStatus))
	env = append(env, fmt.Sprintf("EXIT_CHOOSE_IGN=%d", DepsHookChooseIgnExitStatus))

	proc, err := os.StartProcess(hook, []string{hook},
		&os.ProcAttr{
			Dir:   os.Getenv("PWD"),
			Env:   env,
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		},
	)
	if err != nil {
		return 1, err
	}

	state, err := proc.Wait()
	if err != nil {
		return 1, err
	}

	pstate := state.Sys().(syscall.WaitStatus)
	exitStatus := pstate.ExitStatus()

	if exitStatus > 0 && exitStatus < 127 {
		return exitStatus, fmt.Errorf("hook failed: rc=%d", exitStatus)
	}
	return exitStatus, nil
}

func Merge(diff *GoDepsDiff, hooks *MergeHooks) (*Godeps, error) {
	deps := &Godeps{}

	if _, err := runHook(hooks.PreHook, []string{}); err != nil {
		return nil, err
	}

	for _, v := range diff.Old.Packages {
		if !utils.InSliceString(v, deps.Packages) {
			deps.Packages = append(deps.Packages, v)
		}
	}
	for _, v := range diff.New.Packages {
		if !utils.InSliceString(v, deps.Packages) {
			deps.Packages = append(deps.Packages, v)
		}
	}

	deps.ImportPath = diff.New.ImportPath
	deps.GoVersion = diff.New.GoVersion
	deps.GodepVersion = diff.New.GodepVersion

	var keys sort.StringSlice
	for k := range diff.ImportPath {
		keys = append(keys, k)
	}
	keys.Sort()

	for _, path := range keys {
		hookChoose, err := runHook(hooks.DepHook,
			[]string{
				fmt.Sprintf("GODEP_IMPORT_PATH=%s", path),
				fmt.Sprintf("GODEP_OLD_COMMENT=%s", diff.ImportPath[path]["old"].Comment),
				fmt.Sprintf("GODEP_OLD_REV=%s", diff.ImportPath[path]["old"].Rev),
				fmt.Sprintf("GODEP_NEW_COMMENT=%s", diff.ImportPath[path]["new"].Comment),
				fmt.Sprintf("GODEP_NEW_REV=%s", diff.ImportPath[path]["new"].Rev),
			},
		)
		if err != nil {
			return nil, err
		}

		switch hookChoose {
		case DepsHookChooseOldExitStatus:
			deps.Deps = append(deps.Deps, diff.ImportPath[path]["old"])
		case DepsHookChooseNewExitStatus:
			deps.Deps = append(deps.Deps, diff.ImportPath[path]["new"])
		case DepsHookChooseIgnExitStatus:
		default:
			filled := 0
			if diff.ImportPath[path]["old"].Rev != "" {
				filled += 1
			}
			if diff.ImportPath[path]["new"].Rev != "" {
				filled += 2
			}
			switch filled {
			case 1:
				deps.Deps = append(deps.Deps, diff.ImportPath[path]["old"])
			case 2, 3:
				deps.Deps = append(deps.Deps, diff.ImportPath[path]["new"])
			}
		}
	}

	if _, err := runHook(hooks.PostHook, []string{}); err != nil {
		return nil, err
	}

	return deps, nil
}

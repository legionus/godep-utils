package godeps

import (
	"fmt"
	"os"
)

const (
	InspectHookEnv = "GODEP_INSPECT_HOOK"
)

func Inspect(deps *Godeps) error {
	for _, dep := range deps.Deps {
		_, err := runHook(os.Getenv(InspectHookEnv),
			[]string{
				fmt.Sprintf("GODEP_IMPORT_PATH=%s", dep.ImportPath),
				fmt.Sprintf("GODEP_COMMENT=%s", dep.Comment),
				fmt.Sprintf("GODEP_REV=%s", dep.Rev),
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

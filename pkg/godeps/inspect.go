package godeps

func Inspect(deps *Godeps, hook string) error {
	for _, dep := range deps.Deps {
		_, err := runHook(hook,
			[]string{
				"GODEP_IMPORT_PATH=" + dep.ImportPath,
				"GODEP_COMMENT=" + dep.Comment,
				"GODEP_REV=" + dep.Rev,
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

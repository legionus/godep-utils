package godeps

func MakeDiff(oldGodeps, newGodeps *Godeps) *GoDepsDiff {
	diff := &GoDepsDiff{
		Old:        oldGodeps,
		New:        newGodeps,
		ImportPath: make(map[string]DepDiff),
	}

	for i := range oldGodeps.Deps {
		oldDep := oldGodeps.Deps[i]

		if _, ok := diff.ImportPath[oldDep.ImportPath]; ok {
			// Already seen
			continue
		}

		found := false
		for j := range newGodeps.Deps {
			newDep := newGodeps.Deps[j]

			if newDep.ImportPath == oldDep.ImportPath {
				diff.ImportPath[newDep.ImportPath] = DepDiff{
					"old": oldDep,
					"new": newDep,
				}
				found = true
				break
			}
		}

		if found {
			continue
		}

		diff.ImportPath[oldDep.ImportPath] = DepDiff{
			"old": oldDep,
		}
	}

	for i := range newGodeps.Deps {
		newDep := newGodeps.Deps[i]

		if _, ok := diff.ImportPath[newDep.ImportPath]; ok {
			// Already seen
			continue
		}

		found := false
		for j := range oldGodeps.Deps {
			oldDep := oldGodeps.Deps[j]

			if newDep.ImportPath == oldDep.ImportPath {
				found = true
				break
			}
		}

		if found {
			continue
		}

		diff.ImportPath[newDep.ImportPath] = DepDiff{
			"new": newDep,
		}
	}

	return diff
}

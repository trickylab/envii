package model

// ImportSummary describes the result of importing variables into an env.
type ImportSummary struct {
	Added       int
	Overwritten int
	Skipped     int
}

// ImportVars merges vars into env. Existing keys are skipped unless overwrite is true.
func ImportVars(env *Env, vars []*Var, overwrite bool) ImportSummary {
	index := make(map[string]int, len(env.Vars))
	for i, v := range env.Vars {
		index[v.Key] = i
	}

	var summary ImportSummary
	for _, incoming := range vars {
		if i, ok := index[incoming.Key]; ok {
			if !overwrite {
				summary.Skipped++
				continue
			}
			env.Vars[i] = incoming
			summary.Overwritten++
			continue
		}
		env.Vars = append(env.Vars, incoming)
		index[incoming.Key] = len(env.Vars) - 1
		summary.Added++
	}
	return summary
}

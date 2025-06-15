package model

type OptionalVars []*OptionalVar

func (vs OptionalVars) existByName(name string) bool {
	for _, v := range vs {
		if v.Name == name {
			return true
		}
	}

	return false
}

func (vs OptionalVars) Distinct() OptionalVars {
	vars := make(OptionalVars, 0, len(vs))
	for _, v := range vs {
		if vars.existByName(v.Name) {
			continue
		}

		vars = append(vars, v)
	}

	return vars
}

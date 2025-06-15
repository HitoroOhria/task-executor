package model

type RequiredVars []*RequiredVar

func (vs RequiredVars) existByName(name string) bool {
	for _, v := range vs {
		if v.Name == name {
			return true
		}
	}

	return false
}

func (vs RequiredVars) Distinct() RequiredVars {
	vars := make(RequiredVars, 0, len(vs))
	for _, v := range vs {
		if vars.existByName(v.Name) {
			continue
		}

		vars = append(vars, v)
	}

	return vars
}

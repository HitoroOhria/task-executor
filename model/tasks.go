package model

type Tasks []*Task

func (ts Tasks) FindByName(name string) *Task {
	for _, t := range ts {
		if t.Name == name {
			return t
		}
	}

	return nil
}

func (ts Tasks) FindByFullName(fullName FullTaskName) *Task {
	for _, t := range ts {
		if t.FullName == fullName {
			return t
		}
	}

	return nil
}

func (ts Tasks) FindSelected() *Task {
	for _, t := range ts {
		if t.Selected {
			return t
		}
	}

	return nil
}

package main

type Tasks []*Task

func (ts Tasks) FindByName(name string) *Task {
	for _, t := range ts {
		if t.Name == name {
			return t
		}
	}

	return nil
}

package domain

type Todo struct {
	ID          int
	Description string
	Done        bool
}

func (t *Todo) ToggleComplete() {
	t.Done = !t.Done
}

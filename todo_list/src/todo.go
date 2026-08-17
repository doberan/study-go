package main

type Todo struct {
	ID          int
	Description string
	Done        bool
}

func addTodo(repository *TodoRepository, description string) error {
	newTodo := Todo{Description: description, Done: false}

	return repository.Create(newTodo)
}

func deleteTodo(repository *TodoRepository, id int) error {
	return repository.Delete(id)
}

// Todo一覧を表示
func renderTodoList(repository *TodoRepository) ([]Todo, error) {

	return repository.FindAll()

}

func toggleTodoStatus(repository *TodoRepository, id int) (Todo, error) {
	todo, err := repository.FindByID(id)
	if err != nil {
		return Todo{}, err
	}

	todo.Done = !todo.Done

	if err := repository.Update(todo); err != nil {
		return Todo{}, err
	}
	return todo, nil
}

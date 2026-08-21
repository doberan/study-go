package main

type TodoService interface {
	List() ([]Todo, error)
	Create(description string) (Todo, error)
	Update(todo Todo) (Todo, error)
	Delete(id int) error
}

type todoService struct {
	repository TodoRepository
}

func NewTodoService(repository TodoRepository) *todoService {
	return &todoService{repository: repository}
}

func (s *todoService) List() ([]Todo, error) {
	return s.repository.FindAll()
}

func (s *todoService) Create(description string) (Todo, error) {
	todo := Todo{Description: description, Done: false}

	if err := s.repository.Create(&todo); err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (s *todoService) Update(todo Todo) (Todo, error) {
	if err := s.repository.Update(todo); err != nil {
		return Todo{}, err
	}
	return todo, nil
}

func (s *todoService) Delete(id int) error {
	return s.repository.Delete(id)
}

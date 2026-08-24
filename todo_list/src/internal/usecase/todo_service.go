package usecase

import (
	"todo_list/internal/domain"
)

type TodoService interface {
	List() ([]domain.Todo, error)
	Create(description string) (domain.Todo, error)
	Update(todo domain.Todo) (domain.Todo, error)
	Delete(id int) error
}

type todoService struct {
	repository domain.TodoRepository
}

func NewTodoService(repository domain.TodoRepository) *todoService {
	return &todoService{repository: repository}
}

func (s *todoService) List() ([]domain.Todo, error) {
	return s.repository.FindAll()
}

func (s *todoService) Create(description string) (domain.Todo, error) {
	todo := domain.Todo{Description: description, Done: false}

	if err := s.repository.Create(&todo); err != nil {
		return domain.Todo{}, err
	}

	return todo, nil
}

func (s *todoService) Update(todo domain.Todo) (domain.Todo, error) {
	if err := s.repository.Update(todo); err != nil {
		return domain.Todo{}, err
	}
	return todo, nil
}

func (s *todoService) Delete(id int) error {
	return s.repository.Delete(id)
}

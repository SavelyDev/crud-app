package service

import (
	"github.com/SavelyDev/crud-app/internal/domain"
)

type TodoList interface {
	CreateList(userId int, todoList domain.TodoList) (int, error)
	GetAllLists(userId int) ([]domain.TodoList, error)
	GetListById(userId, listId int) (domain.TodoList, error)
	DeleteList(userId, listId int) error
	UpdateList(userId, listId int, input domain.UpdateListInput) error
}

type TodoListService struct {
	repo TodoList
}

func NewTodoListService(repo TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) CreateList(userId int, todoList domain.TodoList) (int, error) {
	return s.repo.CreateList(userId, todoList)
}

func (s *TodoListService) GetAllLists(userId int) ([]domain.TodoList, error) {
	return s.repo.GetAllLists(userId)
}

func (s *TodoListService) GetListById(userId, listId int) (domain.TodoList, error) {
	return s.repo.GetListById(userId, listId)
}

func (s *TodoListService) DeleteList(userId, listId int) error {
	return s.repo.DeleteList(userId, listId)
}

func (s *TodoListService) UpdateList(userId, listId int, input domain.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.UpdateList(userId, listId, input)
}

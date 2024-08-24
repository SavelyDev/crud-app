package service

import (
	"github.com/SavelyDev/crud-app/internal/domain"
)

type TodoItem interface {
	CreateItem(listId int, input domain.TodoItem) (int, error)
	GetAllItems(userId, listId int) ([]domain.TodoItem, error)
	GetItemById(userId, itemId int) (domain.TodoItem, error)
	DeleteItem(userId, itemId int) error
	UpdateItem(userId, itemId int, input domain.UpdateItemInput) error
}

type TodoItemService struct {
	repo     TodoItem
	listRepo TodoList
}

func NewTodoItemService(repo TodoItem, listRepo TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (s *TodoItemService) CreateItem(userId, listId int, input domain.TodoItem) (int, error) {
	_, err := s.listRepo.GetListById(userId, listId)
	if err != nil {
		return 0, err
	}

	return s.repo.CreateItem(listId, input)
}

func (s *TodoItemService) GetAllItems(userId, listId int) ([]domain.TodoItem, error) {
	return s.repo.GetAllItems(userId, listId)
}

func (s *TodoItemService) GetItemById(userId, itemId int) (domain.TodoItem, error) {
	return s.repo.GetItemById(userId, itemId)
}

func (s *TodoItemService) DeleteItem(userId, itemId int) error {
	return s.repo.DeleteItem(userId, itemId)
}

func (s *TodoItemService) UpdateItem(userId, itemId int, input domain.UpdateItemInput) error {
	return s.repo.UpdateItem(userId, itemId, input)
}

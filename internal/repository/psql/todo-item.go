package psql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/SavelyDev/crud-app/internal/domain"
)

type TodoItemRepo struct {
	db *sql.DB
}

func NewTodoItemRepo(db *sql.DB) *TodoItemRepo {
	return &TodoItemRepo{db: db}
}

func (r *TodoItemRepo) CreateItem(listId int, todoItem domain.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	row := tx.QueryRow("INSERT INTO todo_items (title, description) VALUES ($1, $2) RETURNING id",
		todoItem.Title, todoItem.Description)
	if err := row.Scan(&itemId); err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec("INSERT INTO lists_items (item_id, list_id) VALUES ($1, $2)",
		itemId, listId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItemRepo) GetAllItems(userId, listId int) ([]domain.TodoItem, error) {
	var items []domain.TodoItem

	rows, err := r.db.Query(`SELECT ti.id, ti.title, ti.description, ti.done FROM todo_items ti 
	JOIN lists_items li ON ti.id = li.item_id 
	JOIN users_lists ul ON li.list_id = ul.list_id
	WHERE ul.user_id=$1 AND li.list_id=$2`, userId, listId)
	if err != nil {
		return items, err
	}

	var item domain.TodoItem
	for rows.Next() {
		if err := rows.Scan(&item.Id, &item.Title, &item.Description, &item.Done); err != nil {
			return items, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *TodoItemRepo) GetItemById(userId, itemId int) (domain.TodoItem, error) {
	var item domain.TodoItem

	row := r.db.QueryRow(`SELECT ti.id, ti.title, ti.description, ti.done FROM todo_items ti 
	JOIN lists_items li ON ti.id = li.item_id 
	JOIN users_lists ul ON li.list_id = ul.list_id
	WHERE ul.user_id=$1 AND ti.id=$2`, userId, itemId)
	if err := row.Scan(&item.Id, &item.Title, &item.Description, &item.Done); err != nil {
		return item, err
	}

	return item, nil
}

func (r *TodoItemRepo) UpdateItem(userId, itemId int, input domain.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	args = append(args, userId, itemId)

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE todo_items ti SET %s FROM lists_items li, users_lists ul
	WHERE ti.id = li.item_id AND li.list_id = ul.list_id 
	AND ul.user_id = $%d AND ti.id = $%d`, setQuery, argId, argId+1)

	_, err := r.db.Exec(query, args...)

	return err
}

func (r *TodoItemRepo) DeleteItem(userId, itemId int) error {
	_, err := r.db.Exec(`DELETE FROM todo_items ti
	USING lists_items li, users_lists ul
	WHERE ti.id = li.item_id AND li.list_id = ul.list_id 
	AND ul.user_id=$1 AND ti.id=$2 `, userId, itemId)

	return err
}

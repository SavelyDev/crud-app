package psql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/SavelyDev/crud-app/internal/domain"
)

type TodoListRepo struct {
	db *sql.DB
}

func NewTodoListRepo(db *sql.DB) *TodoListRepo {
	return &TodoListRepo{db: db}
}

func (r *TodoListRepo) CreateList(userId int, todoList domain.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var listId int
	row := tx.QueryRow("INSERT INTO todo_lists (title, description) VALUES ($1, $2) RETURNING id",
		todoList.Title, todoList.Description)
	if err := row.Scan(&listId); err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec("INSERT INTO users_lists (user_id, list_id) VALUES ($1, $2)",
		userId, listId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return listId, tx.Commit()
}

func (r *TodoListRepo) GetAllLists(userId int) ([]domain.TodoList, error) {
	var lists []domain.TodoList

	rows, err := r.db.Query(`SELECT tl.id, tl.title, tl.description FROM todo_lists tl 
							JOIN users_lists ul ON tl.id = ul.list_id 
							WHERE ul.user_id = $1`, userId)
	if err != nil {
		return lists, err
	}

	for rows.Next() {
		var list domain.TodoList

		if err := rows.Scan(&list.Id, &list.Title, &list.Description); err != nil {
			return lists, err
		}

		lists = append(lists, list)
	}

	return lists, nil
}

func (r *TodoListRepo) GetListById(userId, listId int) (domain.TodoList, error) {
	var list domain.TodoList

	row := r.db.QueryRow(`SELECT tl.id, tl.title, tl.description FROM todo_lists tl 
							JOIN users_lists ul ON tl.id = ul.list_id 
							WHERE ul.user_id = $1 AND ul.list_id = $2`, userId, listId)
	if err := row.Scan(&list.Id, &list.Title, &list.Description); err != nil {
		return list, err
	}

	return list, nil
}

func (r *TodoListRepo) UpdateList(userId, listId int, input domain.UpdateListInput) error {
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

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE todo_lists tl SET %s FROM users_lists ul 
	WHERE tl.id = ul.list_id AND ul.user_id = $%d AND ul.list_id = $%d`, setQuery, argId, argId+1)

	args = append(args, userId, listId)

	_, err := r.db.Exec(query, args...)

	return err
}

func (r *TodoListRepo) DeleteList(userId, listId int) error {
	_, err := r.db.Exec(`DELETE FROM todo_lists tl USING users_lists ul 
						WHERE tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2`, userId, listId)

	return err
}

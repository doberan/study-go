package main

import (
	"database/sql"
)

type TodoRepository struct {
	db *sql.DB
}

func (r *TodoRepository) FindAll() ([]Todo, error) {
	rows, err := r.db.Query(`
		SELECT id, description, done
		FROM todos
		ORDER BY id ASC
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var todo Todo

		err := rows.Scan(
			&todo.ID,
			&todo.Description,
			&todo.Done,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepository) FindByID(id int) (Todo, error) {
	var todo Todo
	err := r.db.QueryRow(`SELECT id, description, done FROM todos WHERE id = ?`, id).Scan(
		&todo.ID,
		&todo.Description,
		&todo.Done,
	)
	if err != nil {
		return Todo{}, err
	}
	return todo, nil
}

func (r *TodoRepository) Create(todo Todo) error {
	_, err := r.db.Exec(`INSERT INTO todos (description, done) VALUES (?, ?)`,
		todo.Description, todo.Done)
	if err != nil {
		return err
	}
	return nil
}

func (r *TodoRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM todos WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *TodoRepository) Update(todo Todo) error {
	result, err := r.db.Exec(`UPDATE todos SET description = ?, done = ? WHERE id = ?`, todo.Description, todo.Done, todo.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

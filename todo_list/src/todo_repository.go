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

package domain

import (
	"database/sql"
)

type TodoRepository interface {
	FindAll() ([]Todo, error)
	FindByID(id int) (Todo, error)
	Create(todo *Todo) error
	Delete(id int) error
	Update(todo Todo) error
}

type MySQLTodoRepository struct {
	db *sql.DB
}

func NewMySQLTodoRepository(db *sql.DB) *MySQLTodoRepository {
	return &MySQLTodoRepository{db: db}
}

func (r *MySQLTodoRepository) FindAll() ([]Todo, error) {
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

func (r *MySQLTodoRepository) FindByID(id int) (Todo, error) {
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

func (r *MySQLTodoRepository) Create(todo *Todo) error {
	result, err := r.db.Exec(`INSERT INTO todos (description, done) VALUES (?, ?)`,
		todo.Description, todo.Done)
	if err != nil {
		return err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	todo.ID = int(lastID)

	return nil
}

func (r *MySQLTodoRepository) Delete(id int) error {
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

func (r *MySQLTodoRepository) Update(todo Todo) error {
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

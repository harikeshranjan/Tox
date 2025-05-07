package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Todo represents a task in our todo list
type Todo struct {
	ID        int64
	Task      string
	Done      bool
	CreatedAt time.Time
	DoneAt    *time.Time
}

// Manager handles database operations
type Manager struct {
	db *sql.DB
}

// NewManager creates a new database manager
func NewManager() (*Manager, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get home directory: %w", err)
	}

	// Create .tox directory if it doesn't exist
	toxDir := filepath.Join(homeDir, ".tox")
	if err := os.MkdirAll(toxDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create tox directory: %w", err)
	}

	// Connect to SQLite database
	dbPath := filepath.Join(toxDir, "todos.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}

	// Create table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			task TEXT NOT NULL,
			done BOOLEAN DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			done_at TIMESTAMP NULL
		)
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("could not create todos table: %w", err)
	}

	return &Manager{db: db}, nil
}

// Close closes the database connection
func (m *Manager) Close() error {
	return m.db.Close()
}

// AddTodo adds a new todo to the database
func (m *Manager) AddTodo(task string) (int64, error) {
	result, err := m.db.Exec("INSERT INTO todos (task) VALUES (?)", task)
	if err != nil {
		return 0, fmt.Errorf("could not add todo: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("could not get last insert id: %w", err)
	}
	return id, nil
}

// ListTodos returns all todos from the database
func (m *Manager) ListTodos(showDone bool) ([]Todo, error) {
	query := "SELECT id, task, done, created_at, done_at FROM todos"
	if !showDone {
		query += " WHERE done = 0"
	}
	query += " ORDER BY done, id"

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not query todos: %w", err)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		var doneAt sql.NullTime
		err := rows.Scan(&todo.ID, &todo.Task, &todo.Done, &todo.CreatedAt, &doneAt)
		if err != nil {
			return nil, fmt.Errorf("could not scan todo: %w", err)
		}
		if doneAt.Valid {
			todo.DoneAt = &doneAt.Time
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// MarkDone marks a todo as done
func (m *Manager) MarkDone(id int64) error {
	now := time.Now()
	result, err := m.db.Exec("UPDATE todos SET done = 1, done_at = ? WHERE id = ?", now, id)
	if err != nil {
		return fmt.Errorf("could not mark todo as done: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("todo with id %d not found", id)
	}
	return nil
}

// ReindexAll resets all todo IDs to be sequential
func (m *Manager) ReindexAll() error {
	// Begin a transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback() // Will be a no-op if transaction succeeds

	// Reindex the todos
	if err := m.reindexTodos(tx); err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

// reindexTodos resets the IDs of all todos to be sequential
func (m *Manager) reindexTodos(tx *sql.Tx) error {
	// Create a temporary table
	_, err := tx.Exec(`
		CREATE TEMPORARY TABLE todos_temp AS
		SELECT NULL as id, task, done, created_at, done_at
		FROM todos
		ORDER BY id
	`)
	if err != nil {
		return err
	}

	// Delete all rows from the original table
	_, err = tx.Exec("DELETE FROM todos")
	if err != nil {
		return err
	}

	// Reset the AUTOINCREMENT counter
	_, err = tx.Exec("DELETE FROM sqlite_sequence WHERE name='todos'")
	if err != nil {
		return err
	}

	// Insert rows from temporary table back to original table
	_, err = tx.Exec(`
		INSERT INTO todos (task, done, created_at, done_at)
		SELECT task, done, created_at, done_at
		FROM todos_temp
	`)
	if err != nil {
		return err
	}

	// Drop the temporary table
	_, err = tx.Exec("DROP TABLE todos_temp")
	if err != nil {
		return err
	}

	return nil
}

// DeleteTodo removes a todo from the database
func (m *Manager) DeleteTodo(id int64) error {
	// Begin a transaction to ensure atomicity
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback() // Will be a no-op if transaction succeeds

	// Delete the todo
	result, err := tx.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("could not delete todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("todo with id %d not found", id)
	}

	// Reindex the todos
	if err := m.reindexTodos(tx); err != nil {
		return fmt.Errorf("could not reindex todos: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

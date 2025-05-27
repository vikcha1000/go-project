package model

import (
	"time"
)

type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Task        string    `gorm:"type:text;not null" json:"task"`
	AuthorID    uint      `gorm:"not null" json:"author_id"`
	ExecutorID  uint      `gorm:"not null" json:"executor_id"`
	Deadline    time.Time `gorm:"not null" json:"deadline"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsDone      bool      `gorm:"default:false" json:"is_done"`
}


// type TaskResponse struct {
// 	ID          uint      `gorm:"primaryKey" json:"id"`
// 	Name        string    `gorm:"size:255;not null" json:"name"`
// 	Task        string    `gorm:"type:text;not null" json:"task" binding:"required"`
// 	AuthorID    uint      `gorm:"not null" json:"-"`
// 	Author      User      `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"author" binding:"required"`
// 	ExecutorID  uint      `gorm:"not null" json:"-"`
// 	Executor    User      `gorm:"foreignKey:ExecutorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"executor" binding:"required"`
// 	Deadline    time.Time `gorm:"not null" json:"deadline" binding:"required"`
// 	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
// 	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"-"`
// 	IsDone      bool      `gorm:"default:false" json:"isDone"`
// }

// type TaskRequest struct {
// 	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
// 	Name       string    `gorm:"size:255;not null" json:"name"`
// 	Task       string    `gorm:"type:text;not null" json:"task" binding:"required"`
// 	Author     string    `gorm:"size:255;not null" json:"author" binding:"required"`
// 	Executor   string    `gorm:"size:255;not null" json:"executor" binding:"required"`
// 	Deadline   time.Time `gorm:"not null" json:"deadline" binding:"required"`
// 	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
// 	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
// 	IsDone     bool      `gorm:"default:false" json:"isDone"`
// }

// type User struct {
// 		ID        int       `gorm:"primaryKey"`
// 	Name             string `gorm:"size:255" json:"name"`
// 	TelegramUsername string `gorm:"size:255" json:"telegram_username"`
// }

// func CreateTable(db *sqlx.DB) error {
// 	tasksQuery := `
//     CREATE TABLE IF NOT EXISTS tasks (
//         id SERIAL PRIMARY KEY,
//         name VARCHAR(255) NOT NULL,
//         task VARCHAR(255) NOT NULL,
//         author VARCHAR(255) NOT NULL,
//         executor VARCHAR(255) NOT NULL,
//         deadline TIMESTAMP WITH TIME ZONE,
//         is_done BOOLEAN DEFAULT FALSE,
//         created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
//         CONSTRAINT fk_author FOREIGN KEY (author) REFERENCES users(telegram_username) ON UPDATE CASCADE,
//         CONSTRAINT fk_executor FOREIGN KEY (executor) REFERENCES users(telegram_username) ON UPDATE CASCADE
//     )`

// 	usersQuery := `
//     CREATE TABLE IF NOT EXISTS users (
//         id SERIAL PRIMARY KEY,
//         name VARCHAR(255) NOT NULL,
//         telegram_username VARCHAR(255) NOT NULL UNIQUE
//     )`

// 	tx, err := db.Beginx()
// 	if err != nil {
// 		return fmt.Errorf("failed to begin transaction: %w", err)
// 	}

// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	if _, err = tx.Exec(usersQuery); err != nil {
// 		return fmt.Errorf("failed to create users table: %w", err)
// 	}

// 	if _, err = tx.Exec(tasksQuery); err != nil {
// 		return fmt.Errorf("failed to create tasks table: %w", err)
// 	}

// 	return tx.Commit()

// }

// func CreateTask(db *sqlx.DB, task *TaskRequest) error {
// 	query := `INSERT INTO tasks (name, task, author, executor, deadline) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, is_done`
// 	return db.QueryRowx(query, task.Name, task.Task, task.Author, task.Executor, task.Deadline).StructScan(task)
// }

// func GetAllTasks(db *sqlx.DB) ([]TaskResponse, error) {
// 	var tasks []TaskResponse

// 	query := `
//         SELECT 
//         t.id,
//         t.name,
//         t.task,
//         t.deadline,
//         t.created_at,
//         t.is_done,
//         -- Данные автора
//         author.id as "author.id",
//         author.name as "author.name",
//         author.telegram_username as "author.telegram_username",
//         -- Данные исполнителя
//         executor.id as "executor.id",
//         executor.name as "executor.name",
//         executor.telegram_username as "executor.telegram_username"
//     FROM tasks t
//     JOIN users author ON t.author = author.telegram_username
//     JOIN users executor ON t.executor = executor.telegram_username
//     ORDER BY t.created_at DESC`
// 	err := db.Select(&tasks, query)
// 	return tasks, err
// }

// func CreateUser(db *sqlx.DB, user *User) error {
// 	query := `INSERT INTO users (name, telegram_username) VALUES ($1, $2) RETURNING id`
// 	return db.QueryRowx(query, user.Name, user.TelegramUsername).StructScan(user)
// }

// // TODO почему id string
// func GetUserByID(db *sqlx.DB, id string) (*User, error) {
// 	var user User
// 	query := `SELECT * FROM users WHERE id = $1`
// 	err := db.Get(&user, query, id)

// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, nil // Пользователь не найден
// 		}
// 		return nil, err
// 	}

// 	return &user, nil
// }

// func UpdateUserByID(db *sqlx.DB, id string, updatedUser *User) (*User, error) {
// 	query := `
// 		UPDATE users 
// 		SET name = $1, telegram_username = $2 
// 		WHERE id = $3
// 		RETURNING id, name, telegram_username`

// 	var user User
// 	err := db.QueryRowx(query,
// 		updatedUser.Name,
// 		updatedUser.TelegramUsername,
// 		id).StructScan(&user)

// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, nil // Пользователь не найден
// 		}
// 		return nil, err
// 	}

// 	return &user, nil
// }

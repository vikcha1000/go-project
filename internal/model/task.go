package model

import (
	"time"
)

type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text;not null" json:"description"`
	AuthorID    uint      `gorm:"not null" json:"-"`
	Author      User      `gorm:"foreignKey:AuthorID" json:"author"`
	ExecutorID  uint      `gorm:"not null" json:"-"`
	Executor    User      `gorm:"foreignKey:ExecutorID" json:"executor"`
	Deadline    time.Time `gorm:"not null" json:"deadline"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsDone      bool      `gorm:"default:false" json:"is_done"`
}

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

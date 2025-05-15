package models

import (
    "time"
    "github.com/jmoiron/sqlx"
)

type Item struct {
    ID        int       `db:"id" json:"id"`
    Name      string    `db:"name" json:"name"`
    Value     int       `db:"value" json:"value"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func CreateTable(db *sqlx.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS items (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        value INTEGER NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    )`
    _, err := db.Exec(query)
    return err
}

func CreateItem(db *sqlx.DB, item *Item) error {
    query := `INSERT INTO items (name, value) VALUES ($1, $2) RETURNING id, created_at`
    return db.QueryRowx(query, item.Name, item.Value).StructScan(item)
}

func GetAllItems(db *sqlx.DB) ([]Item, error) {
    var items []Item
    query := `SELECT id, name, value, created_at FROM items ORDER BY created_at DESC`
    err := db.Select(&items, query)
    return items, err
}
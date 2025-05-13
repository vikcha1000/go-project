package model

import "time"

// DataInput - модель для входящих данных
type DataInput struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// DataOutput - модель для исходящих данных
type DataOutput struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Value     int       `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

// GetAllData - получает все данные (заглушка)
func GetAllData() ([]DataOutput, error) {
	// Здесь должна быть логика получения данных из БД
	// Сейчас это заглушка
	return []DataOutput{
		{ID: 1, Name: "Тест", Value: 100, CreatedAt: time.Now()},
	}, nil
}

// CreateData - создает новые данные (заглушка)
func CreateData(input DataInput) (*DataOutput, error) {
	// Здесь должна быть логика сохранения в БД
	// Сейчас это заглушка
	return &DataOutput{
		ID:        2,
		Name:      input.Name,
		Value:     input.Value,
		CreatedAt: time.Now(),
	}, nil
}

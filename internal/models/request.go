package models

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

// Request - базовый интерфейс для всех запросов
type Request interface {
	Validate() error
}

// ParseAndValidate - общая функция для парсинга и валидации запросов
func ParseAndValidate(r *http.Request, req Request) error {
	// Проверяем, есть ли тело запроса
	if r.Body == nil {
		return errors.New("тело запроса отсутствует")
	}

	// Декодируем тело запроса в структуру
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		if err == io.EOF {
			return errors.New("тело запроса пустое")
		}
		return err
	}

	// Валидируем структуру
	return req.Validate()
}

// UserCreateRequest - модель запроса для создания пользователя
type UserCreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Validate реализует интерфейс Request
func (r *UserCreateRequest) Validate() error {
	// Просто проверяем наличие username
	if strings.TrimSpace(r.Username) == "" {
		return errors.New("имя пользователя обязательно")
	}
	return nil
}

// UserUpdateRequest - модель запроса для обновления пользователя
type UserUpdateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Validate реализует интерфейс Request
func (r *UserUpdateRequest) Validate() error {
	// Просто проверяем наличие username
	if strings.TrimSpace(r.Username) == "" {
		return errors.New("имя пользователя обязательно")
	}
	return nil
}

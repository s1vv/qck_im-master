package errors

import (
	"fmt"
	"runtime"
)

// Базовая ошибка с контекстом
type BaseError struct {
	Op   string // Операция
	Err  error  // Оригинальная ошибка
	File string // Файл
	Line int    // Строка
}

func (e *BaseError) Error() string {
	return fmt.Sprintf("[%s:%d] в файле %s: %v", e.File, e.Line, e.Op, e.Err)
}

// Ошибка для работы с БД
type DBError struct {
	BaseError
}

func NewDBError(op string, err error) error {
	_, file, line, _ := runtime.Caller(2)
	return &DBError{BaseError{Op: op, Err: err, File: file, Line: line}}
}

// Ошибка для HTTP-запросов
type HTTPError struct {
	BaseError
	StatusCode int
}

func NewHTTPError(op string, statusCode int, err error) error {
	_, file, line, _ := runtime.Caller(2)
	return &HTTPError{BaseError{Op: op, Err: err, File: file, Line: line}, statusCode}
}

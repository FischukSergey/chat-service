package errhandler

import (
	clientv1 "github.com/FischukSergey/chat-service/internal/server-client/v1"
)

type Response struct {
	Error clientv1.Error `json:"error"`
}

var ResponseBuilder = func(code int, msg string, details string) any {
	// Фикс: возвращаем Response с ошибкой, которая соответствует коду и сообщению.
	// Фикс: Если детали не пустые, то добавляем их в ошибку.
	// Фикс: Если детали пустые, то оставляем ошибку пустой.
	if details != "" {
		return Response{
			Error: clientv1.Error{
				Code:    code,
				Message: msg,
				Details: &details,
			},
		}
	}
	return Response{
		Error: clientv1.Error{
			Code:    code,
			Message: msg,
		},
	}
}

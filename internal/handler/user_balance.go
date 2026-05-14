package handler

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/repository"
)

/*
Хендлер: GET /api/user/balance

Хендлер доступен только авторизованному пользователю. В ответе должны содержаться данные о текущей сумме баллов лояльности, а также сумме использованных за весь период регистрации баллов.

Формат запроса:

GET /api/user/balance HTTP/1.1
Content-Length: 0
Возможные коды ответа:

200 — успешная обработка запроса.

Формат ответа:

200 OK HTTP/1.1
Content-Type: application/json
...

{
    "current": 500.5,
    "withdrawn": 42
}
401 — пользователь не авторизован;

500 — внутренняя ошибка сервера.
*/

func (h *Handler) UserBalanceHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		h.writeMethodNotAllowed(w, "only GET method is allowed")
		return
	}

	mimeType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mimeType != "application/json" {
		h.writeBadRequest(w, "only application/json content type is supported")
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeBadRequest(w, "failed to read request body: "+err.Error())
		return
	}

	var requestBodySerialised domain.UserLoginRequest
	err = json.Unmarshal(requestBody, &requestBodySerialised)
	if err != nil {
		h.writeBadRequest(w, "couldn't parse requestBody:"+err.Error())
		return
	}

	if requestBodySerialised.Login == "" {
		h.writeBadRequest(w, "login cannot be empty")
		return
	}

	if requestBodySerialised.Password == "" {
		h.writeBadRequest(w, "password cannot be empty")
		return
	}

	authResult, err := h.service.LoginUser(r.Context(), requestBodySerialised)
	if errors.Is(err, repository.ErrUserNotFound) {
		h.writeUnauthorised(w, err.Error())
		return
	}

	if err != nil {
		h.writeInternalServerError(w, err.Error())
		return
	}

	// write an auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    authResult.Token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
}

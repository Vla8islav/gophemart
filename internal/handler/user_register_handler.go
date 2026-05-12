package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Vla8islav/gophemart/internal/models"
	"go.uber.org/zap"
)

/*
Регистрация пользователя
Хендлер: POST /api/user/register

Регистрация производится по паре логин/пароль. Каждый логин должен быть уникальным.

После успешной регистрации должна происходить автоматическая аутентификация пользователя.

Для передачи аутентификационных данных используйте механизм cookies или HTTP-заголовок Authorization.

Формат запроса:

POST /api/user/register HTTP/1.1
Content-Type: application/json
...

{
    "login": "<login>",
    "password": "<password>"
}
Возможные коды ответа:

200 — пользователь успешно зарегистрирован и аутентифицирован;
400 — неверный формат запроса;
409 — логин уже занят;
500 — внутренняя ошибка сервера.
*/

func (h *Handler) writeBadRequest(w http.ResponseWriter, msg string) {
	h.logger.Error("bad request", zap.String("msg", msg))
	http.Error(w, msg, http.StatusBadRequest)
}

func (h *Handler) writeMethodNotAllowed(w http.ResponseWriter, msg string) {
	h.logger.Error("method not allowed: ", zap.String("msg", msg))
	http.Error(w, msg, http.StatusMethodNotAllowed)
}

func (h *Handler) UserRegisterHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusNotImplemented)

	if r.Method != http.MethodPost {
		h.writeMethodNotAllowed(w, "only POST method is allowed")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		h.writeBadRequest(w, "only application/json content type is supported")
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeBadRequest(w, "failed to read request body: "+err.Error())
		return
	}

	var requestBodySerialised []models.User
	err = json.Unmarshal(requestBody, &requestBodySerialised)
	if err != nil {
		h.writeBadRequest(w, "couldn't parse requestBody with metrics :"+err.Error())
		return
	}

	h.service.Ping()

}

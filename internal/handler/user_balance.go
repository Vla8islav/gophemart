package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/helpers"
	"github.com/Vla8islav/gophemart/internal/middlewares"
)

/*
Хендлер: GET /api/user/balance

Хендлер доступен только авторизованному пользователю.
В ответе должны содержаться данные о текущей сумме баллов лояльности,
а также сумме использованных за весь период регистрации баллов.

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

	userID, ok := middlewares.UserIDFromContext(r.Context())
	if !ok {
		h.writeUnauthorised(w, "unauthorized")
		return
	}

	balance, err := h.service.GetUserBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	responseStruct := domain.UserBalanceResponse{
		Current:   helpers.CentsToFloat(balance.Current),
		Withdrawn: helpers.CentsToFloat(balance.Withdrawn),
	}
	response, err := json.Marshal(&responseStruct)
	if err != nil {
		h.writeInternalServerError(w, "couldn't encode the response "+err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		h.writeInternalServerError(w, "couldn't write response"+err.Error())
		return
	}

}

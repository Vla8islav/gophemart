package handler

import (
	"encoding/json"
	"errors"
	"io"
	"math"
	"mime"
	"net/http"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/middlewares"
)

/* Запрос на списание средств
Хендлер: POST /api/user/balance/withdraw

Хендлер доступен только авторизованному пользователю. Номер заказа представляет собой гипотетический номер нового заказа пользователя, в счёт оплаты которого списываются баллы.

Примечание: для успешного списания достаточно успешной регистрации запроса, никаких внешних систем начисления не предусмотрено и не требуется реализовывать.

Формат запроса:

POST /api/user/balance/withdraw HTTP/1.1
Content-Type: application/json

	{
	    "order": "2377225624",
	    "sum": 751
	}

Здесь order — номер заказа, а sum — сумма баллов к списанию в счёт оплаты.

Возможные коды ответа:

200 — успешная обработка запроса;
401 — пользователь не авторизован;
402 — на счету недостаточно средств;
422 — неверный номер заказа;
500 — внутренняя ошибка сервера.
*/

func (h *Handler) UserBalanceWithdrawHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		h.writeMethodNotAllowed(w, "only POST method is allowed")
		return
	}

	userID, ok := middlewares.UserIDFromContext(r.Context())
	if !ok {
		h.writeUnauthorised(w, "unauthorized")
		return
	}

	mimeType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		h.writeInternalServerError(w, "couldn't parse content type"+err.Error())
		return
	}
	if mimeType != "application/json" {
		h.writeBadRequest(w, "only application/json content type is supported")
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeBadRequest(w, "failed to read request body: "+err.Error())
		return
	}

	var requestBodySerialised domain.UserBalanceWithdrawRequest
	err = json.Unmarshal(requestBody, &requestBodySerialised)
	if err != nil {
		h.writeBadRequest(w, "couldn't parse requestBody:"+err.Error())
		return
	}

	if requestBodySerialised.Order == "" {
		h.writeUnprocessableEntity(w, "invalid order number")
		return
	}

	sumCents := int64(math.Round(requestBodySerialised.Sum * 100))
	if sumCents <= 0 {
		h.writeUnprocessableEntity(w, "sum must be positive")
		return
	}

	err = h.service.WithdrawFromUserBalance(r.Context(), userID,
		domain.UserBalanceWithdraw{Sum: sumCents, Order: requestBodySerialised.Order})

	if err != nil {

		switch {
		case errors.Is(err, domain.ErrNotEnoughMoney):
			h.writePaymentRequired(w, "not enough money on balance")

		case errors.Is(err, domain.ErrInvalidOrderNumber):
			h.writeUnprocessableEntity(w, "couldn't find order number: "+err.Error())

		default:
			h.writeInternalServerError(w, "failed to create order: "+err.Error())
		}

		return
	}

	w.WriteHeader(http.StatusOK)

}

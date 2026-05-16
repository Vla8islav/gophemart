package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/helpers"
	"github.com/Vla8islav/gophemart/internal/middlewares"
	"go.uber.org/zap"
)

/*
Получение информации о выводе средств
Хендлер: GET /api/user/withdrawals

Хендлер доступен только авторизованному пользователю. Факты выводов в выдаче должны быть отсортированы по времени вывода от самых новых к самым старым. Формат даты — RFC3339.

Формат запроса:

GET /api/user/withdrawals HTTP/1.1
Content-Length: 0
Возможные коды ответа:

200 — успешная обработка запроса.

Формат ответа:

200 OK HTTP/1.1
Content-Type: application/json
...

[
    {
        "order": "2377225624",
        "sum": 500,
        "processed_at": "2020-12-09T16:09:57+03:00"
    }
]

204 — нет ни одного списания;
401 — пользователь не авторизован;
500 — внутренняя ошибка сервера.
*/

func (h *Handler) UserWithdrawalsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		h.writeMethodNotAllowed(w, "only GET method is allowed")
		return
	}

	userID, ok := middlewares.UserIDFromContext(r.Context())
	if !ok {
		h.writeUnauthorised(w, "unauthorized")
		return
	}

	withdrawals, err := h.service.GetUserWithdrawals(r.Context(), userID)
	if err != nil {
		h.writeInternalServerError(w, "couldn't get user withdrawals "+strconv.FormatInt(userID, 10)+": "+err.Error())
		return
	}
	if len(withdrawals) == 0 {
		h.writeNoContent(w, "no user withdrawals found")
		return
	}

	// constructing response
	response := make(domain.WithdrawResponse, 0, len(withdrawals))
	for _, withdrawal := range withdrawals {
		response = append(response, domain.WithdrawResponseItem{
			Order:       withdrawal.OrderNumber,
			Sum:         helpers.CentsToFloat(withdrawal.Amount),
			ProcessedAt: withdrawal.ProcessedAt,
		})
	}

	marshal, err := json.Marshal(response)

	if err != nil {
		h.writeInternalServerError(w, "couldn't marshal user withdrawals "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(marshal)

	if err != nil {
		h.logger.Error("couldn't write response", zap.Error(err))
		return
	}

}

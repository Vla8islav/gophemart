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
Получение списка загруженных номеров заказов
Хендлер: GET /api/user/orders

Хендлер доступен только авторизованному пользователю. Номера заказа в выдаче должны быть отсортированы по времени загрузки от самых новых к самым старым. Формат даты — RFC3339.

Доступные статусы обработки расчётов:

NEW — заказ загружен в систему, но не попал в обработку;
PROCESSING — вознаграждение за заказ рассчитывается;
INVALID — система расчёта вознаграждений отказала в расчёте;
PROCESSED — данные по заказу проверены и информация о расчёте успешно получена.
Формат запроса:

GET /api/user/orders HTTP/1.1
Content-Length: 0
Возможные коды ответа:

200 — успешная обработка запроса.

Формат ответа:

200 OK HTTP/1.1
Content-Type: application/json
...

[
    {
        "number": "9278923470",
        "status": "PROCESSED",
        "accrual": 500,
        "uploaded_at": "2020-12-10T15:15:45+03:00"
    },
    {
        "number": "12345678903",
        "status": "PROCESSING",
        "uploaded_at": "2020-12-10T15:12:01+03:00"
    },
    {
        "number": "346436439",
        "status": "INVALID",
        "uploaded_at": "2020-12-09T16:09:53+03:00"
    }
]
204 — нет данных для ответа;

401 — пользователь не авторизован;

500 — внутренняя ошибка сервера.
*/

func (h *Handler) UserOrdersGetHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		h.writeMethodNotAllowed(w, "only GET method is allowed")
		return
	}

	userID, ok := middlewares.UserIDFromContext(r.Context())
	if !ok {
		h.writeUnauthorised(w, "unauthorized")
		return
	}

	orders, err := h.service.GetUserOrders(r.Context(), userID)
	if err != nil {
		h.writeInternalServerError(w, "couldn't get user orders "+strconv.FormatInt(userID, 10)+": "+err.Error())
		return
	}
	if len(orders) == 0 {
		h.writeNoContent(w, "no user orders found")
		return
	}

	var response domain.UserOrderGetResponse

	for _, order := range orders {
		var accrualFloatPtr *float64
		if order.Accrual != nil {
			accrualFloat := helpers.CentsToFloat(*order.Accrual)
			accrualFloatPtr = &accrualFloat
		}
		response = append(response, domain.UserOrderGetResponseItem{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    accrualFloatPtr,
			UploadedAt: order.UploadedAt,
		})
	}

	marshal, err := json.Marshal(orders)

	if err != nil {

		h.writeInternalServerError(w, "couldn't marshal user orders "+err.Error())
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

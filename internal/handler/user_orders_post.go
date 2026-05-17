package handler

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/middlewares"
)

/*
Загрузка номера заказа
Хендлер: POST /api/user/orders

Хендлер доступен только аутентифицированным пользователям. Номером заказа является последовательность цифр произвольной длины.

Номер заказа может быть проверен на корректность ввода с помощью алгоритма Луна.

Формат запроса:

POST /api/user/orders HTTP/1.1
Content-Type: text/plain
...

12345678903
Возможные коды ответа:

200 — номер заказа уже был загружен этим пользователем;
202 — новый номер заказа принят в обработку;
400 — неверный формат запроса;
401 — пользователь не аутентифицирован;
409 — номер заказа уже был загружен другим пользователем;
422 — неверный формат номера заказа;
500 — внутренняя ошибка сервера.

*/

func (h *Handler) UserOrdersPostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		h.writeMethodNotAllowed(w, "only POST method is allowed")
		return
	}

	mimeType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mimeType != "text/plain" {
		h.writeBadRequest(w, "only text/plain content type is supported")
		return
	}

	userID, ok := middlewares.UserIDFromContext(r.Context())
	if !ok {
		h.writeUnauthorised(w, "unauthorized")
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeBadRequest(w, "failed to read request body: "+err.Error())
		return
	}

	orderNumberString := strings.TrimSpace(string(requestBody))

	err = h.service.CreateOrder(r.Context(), userID, orderNumberString)
	if err != nil {

		switch {
		case errors.Is(err, domain.ErrOrderAlreadyUploadedByUser):
			w.WriteHeader(http.StatusOK)

		case errors.Is(err, domain.ErrOrderAlreadyUploadedByAnotherUser):
			h.writeAlreadyExists(w, "order number already uploaded by another user")

		case errors.Is(err, domain.ErrInvalidOrderNumber):
			h.writeUnprocessableEntity(w, "order number didn't pass the luhn checksum check")

		case errors.Is(err, domain.ErrOrderEmpty):
			h.writeBadRequest(w, "empty order number")

		default:
			h.writeInternalServerError(w, "failed to create order: "+err.Error())
		}

		return
	}

	w.WriteHeader(http.StatusAccepted)
}

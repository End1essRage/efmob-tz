package http

import (
	"net/http"

	"github.com/end1essrage/efmob-tz/pkg/common/interfaces/http/utils"
)

// CreateSubscription godoc
// @Summary Create subscription
// @Description Create a new subscription
// @Tags subs
// @Accept json
// @Produce json
// @Param request body SubscriptionCreateRequest true "Subscription data"
// @Success 201 {object} Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions [post]
func (h *SubsHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var req SubscriptionCreateRequest
	if !utils.DecodeJSONBody(w, r, &req) {
		return
	}

	writeNotImplemented(w)
}

// GetSubscription godoc
// @Summary Get subscription
// @Description Get subscription by ID
// @Tags subs
// @Produce json
// @Param id path string true "Subscription ID (UUID)"
// @Success 200 {object} Subscription
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/{id} [get]
func (h *SubsHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {

	//id := chi.URLParam(r, "id")

	writeNotImplemented(w)
}

// UpdateSubscription godoc
// @Summary Update subscription
// @Description Update subscription by ID
// @Tags subs
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param request body SubscriptionUpdateRequest true "Updated data"
// @Success 200 {object} Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/{id} [put]
func (h *SubsHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	var req SubscriptionUpdateRequest
	if !utils.DecodeJSONBody(w, r, &req) {
		return
	}

	//id := chi.URLParam(r, "id")
	writeNotImplemented(w)
}

// DeleteSubscription godoc
// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subs
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/{id} [delete]
func (h *SubsHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	//id := chi.URLParam(r, "id")

	writeNotImplemented(w)
}

// ListSubscriptions godoc
// @Summary List subscriptions
// @Description List subscriptions with filters
// @Tags subs
// @Produce json
// @Param user_id query string false "User ID (UUID)"
// @Param service_name query string false "Service name"
// @Param from query string false "Start period (MM-YYYY)"
// @Param to query string false "End period (MM-YYYY)"
// @Success 200 {array} Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions [get]
func (h *SubsHandler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	var req SubscriptionQueryRequest
	if !utils.ParseQuery(w, r, &req) {
		return
	}

	writeNotImplemented(w)
}

// GetTotalCost godoc
// @Summary Calculate total subscription cost
// @Description Calculate total cost for selected period
// @Tags subs
// @Produce json
// @Param user_id query string true "User ID (UUID)"
// @Param service_name query string false "Service name"
// @Param from query string true "Period start (MM-YYYY)"
// @Param to query string true "Period end (MM-YYYY)"
// @Success 200 {object} TotalCostResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/total [get]
func (h *SubsHandler) GetTotalCost(w http.ResponseWriter, r *http.Request) {
	var req TotalCostRequest
	if !utils.ParseQuery(w, r, &req) {
		return
	}

	writeNotImplemented(w)
}
func writeNotImplemented(w http.ResponseWriter) {
	utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
		"msg": "NOT_IMPLEMENTED",
	})
}

/*
func (h *SubsHandler) writeAppError(w http.ResponseWriter, err error) {
	appErr := app.MapDomainError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.HTTPStatus)

	// Безопасно формируем сообщение
	msg := appErr.Code
	if appErr.HTTPStatus >= 500 {
		// Для внутренних ошибок можно не раскрывать детали
		msg = "INTERNAL_ERROR"
	}

	if err := json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
		"code":  appErr.Code,
	}); err != nil {
		logger.Logger().Log("AuthHandler", "writeAppError").Error(err)
	}
}


func formatOptionalDate(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("01-2006")
	return &s
}
*/

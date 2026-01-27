package http

import (
	"net/http"
	"time"

	"github.com/end1essrage/efmob-tz/pkg/common/interfaces/http/utils"
	"github.com/end1essrage/efmob-tz/pkg/common/persistance"
	"github.com/end1essrage/efmob-tz/pkg/subs/application/commands"
	"github.com/end1essrage/efmob-tz/pkg/subs/application/queries"
	"github.com/end1essrage/efmob-tz/pkg/subs/domain"
)

var mapSubscriptionFromDomain = func(record *domain.Subscription) *Subscription {
	return &Subscription{
		ID:          record.ID(),
		UserID:      record.UserID(),
		ServiceName: record.ServiceName(),
		Price:       record.Price(),
		StartDate:   formatDate(record.StartDate()),
		EndDate:     formatOptionalDate(record.EndDate()),
	}
}

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
	sD, err := parseDate(w, req.StartDate)
	if err != nil {
		return
	}
	eD, err := parseOptionalDate(w, req.EndDate)
	if err != nil {
		return
	}

	cmd := commands.CreateSubscriptionCommand{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   sD,
		EndDate:     eD,
	}

	record, err := h.container.CreateSubscriptionHandler.Handle(r.Context(), cmd)
	if err != nil {
		// оборачиваем ошибку
		h.writeAppError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, mapSubscriptionFromDomain(record))
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
	uid, err := extrudeUidFromQuery(w, r)
	if err != nil {
		return
	}

	record, err := h.container.GetSubscriptionHandler.Handle(r.Context(), queries.GetSubscriptionQuery{ID: uid})
	if err != nil {
		// оборачиваем ошибку
		h.writeAppError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, mapSubscriptionFromDomain(record))
}

// UpdateSubscription godoc
// @Summary Update subscription
// @Description Update subscription by ID
// @Tags subs
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param request body SubscriptionUpdateRequest true "Updated data"
// @Success 202 {object} Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/{id} [patch]
func (h *SubsHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	uid, err := extrudeUidFromQuery(w, r)
	if err != nil {
		return
	}

	var req SubscriptionUpdateRequest
	if !utils.DecodeJSONBody(w, r, &req) {
		return
	}

	sD, err := parseDate(w, req.StartDate)
	if err != nil {
		return
	}

	eD, err := parseOptionalDate(w, req.EndDate)
	if err != nil {
		return
	}

	record, err := h.container.UpdateSubscriptionHandler.Handle(r.Context(), commands.UpdateSubscriptionCommand{
		ID:        uid,
		Price:     req.Price,
		StartDate: sD,
		EndDate:   eD,
	})
	if err != nil {
		h.writeAppError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, mapSubscriptionFromDomain(record))
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
	uid, err := extrudeUidFromQuery(w, r)
	if err != nil {
		return
	}

	if err := h.container.DeleteSubscriptionHandler.Handle(r.Context(), commands.DeleteSubscriptionCommand{ID: uid}); err != nil {
		h.writeAppError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
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

	//собираем квери
	// Парсим optional dates
	var period *domain.Period
	if req.From != nil || req.To != nil {
		sD, err := parseOptionalDate(w, req.From)
		if err != nil {
			return
		}

		eD, err := parseOptionalDate(w, req.To)
		if err != nil {
			return
		}

		// Создаем период только если у нас есть хотя бы одна дата
		if sD != nil || eD != nil {
			var start, end time.Time
			if sD != nil {
				start = *sD
			}
			if eD != nil {
				end = *eD
			}

			p, err := domain.NewPeriod(start, end)
			if err != nil {
				h.writeAppError(w, err)
				return
			}
			period = p
		}
	}

	// собираем пагинацию
	var pagination *persistance.Pagination
	if req.PageSize != nil {
		pagination = &persistance.Pagination{Limit: *req.PageSize}
		if req.Page != nil {
			page := *req.Page
			if page < 2 {
				pagination.Offset = 0
			} else {
				pagination.Offset = pagination.Limit * (page - 1)
			}
		} else {
			pagination.Offset = 0
		}
	}

	// собираем сортировку
	var sorting *persistance.Sorting
	if req.OrderBy != nil {
		sorting = &persistance.Sorting{OrderBy: *req.OrderBy}
		if req.Direction != nil {
			sorting.Direction = persistance.SortingDirection(*req.Direction)
		} else {
			sorting.Direction = persistance.DefaultDirection
		}
	}

	// исполняем квери
	records, err := h.container.ListSubscriptionsHandler.Handle(r.Context(), queries.ListSubscriptionsQuery{
		Query:      domain.NewSubscriptionQuery(req.UserID, req.ServiceName, period),
		Pagination: pagination,
		Sorting:    sorting,
	})
	if err != nil {
		h.writeAppError(w, err)
		return
	}

	// маппим ответ
	resp := make([]*Subscription, len(records))
	for i, r := range records {
		resp[i] = mapSubscriptionFromDomain(r)
	}

	utils.WriteJSON(w, http.StatusOK, resp)
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

	//собираем квери
	sD, err := parseDate(w, req.From)
	if err != nil {
		return
	}

	eD, err := parseDate(w, req.To)
	if err != nil {
		return
	}

	period, err := domain.NewPeriod(sD, eD)
	if err != nil {
		h.writeAppError(w, err)
		return
	}

	result, err := h.container.TotalCostHandler.Handle(r.Context(), queries.TotalCostQuery{
		Query: domain.NewSubscriptionQuery(&req.UserID, req.ServiceName, period),
	})
	if err != nil {
		h.writeAppError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, TotalCostResponse{Total: result})
}

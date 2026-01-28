package http

import (
	"net/http"
	"time"

	"github.com/end1essrage/efmob-tz/pkg/common/interfaces/http/utils"
	"github.com/end1essrage/efmob-tz/pkg/common/logger"
	"github.com/end1essrage/efmob-tz/pkg/common/persistance"
	"github.com/end1essrage/efmob-tz/pkg/subs/application/commands"
	"github.com/end1essrage/efmob-tz/pkg/subs/application/queries"
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
	log := logger.Logger().WithFields(logger.LogOptions{
		Pkg:  "SubsHandler",
		Func: "CreateSubscription",
		Ctx:  r.Context(),
	})

	var req SubscriptionCreateRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		log.Errorf("ошибка парсинга тела запроса: %v", err)
		return
	}
	sD, err := parseDate(w, req.StartDate)
	if err != nil {
		log.Errorf("ошибка парсинга даты: %v", err)
		return
	}
	eD, err := parseOptionalDate(w, req.EndDate)
	if err != nil {
		log.Errorf("ошибка парсинга опциональной даты: %v", err)
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
		log.Errorf("ошибка выполнения: %v", err)
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
	log := logger.Logger().WithFields(logger.LogOptions{
		Pkg:  "SubsHandler",
		Func: "GetSubscription",
		Ctx:  r.Context(),
	})

	uid, err := extrudeUidFromQuery(w, r)
	if err != nil {
		return
	}

	record, err := h.container.GetSubscriptionHandler.Handle(r.Context(), queries.GetSubscriptionQuery{ID: uid})
	if err != nil {
		log.Errorf("ошибка выполнения: %v", err)
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
// @Success 202 "Accepted"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/{id} [patch]
func (h *SubsHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	log := logger.Logger().WithFields(logger.LogOptions{
		Pkg:  "SubsHandler",
		Func: "UpdateSubscription",
		Ctx:  r.Context(),
	})

	uid, err := extrudeUidFromQuery(w, r)
	if err != nil {
		return
	}

	var req SubscriptionUpdateRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		log.Errorf("ошибка парсинга тела запроса: %v", err)
		return
	}

	sD, err := parseOptionalDate(w, req.StartDate)
	if err != nil {
		return
	}

	// Обработка EndDate с тремя состояниями
	var endDate *time.Time
	setEndDateNull := false // Флаг нужно ли занулять

	// Поле было в запросе
	if req.EndDate.IsSet() {
		// Занулить
		if req.EndDate.IsNull() {
			endDate = nil         // означает занулить
			setEndDateNull = true // Флаг что нужно занулить
		} else {
			// есть значение
			parsedTime, err := parseDate(w, req.EndDate.Value())
			if err != nil {
				return
			}
			endDate = &parsedTime
		}
	} else {
		// если не было остальных полей
		if req.Price == nil && sD == nil {
			utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: "no fields to update",
				Code:  "UNNECESSARY_UPDATE",
			})
		}
	}

	if err := h.container.UpdateSubscriptionHandler.Handle(r.Context(), commands.UpdateSubscriptionCommand{
		ID:             uid,
		Price:          req.Price,
		StartDate:      sD,
		EndDate:        endDate,
		SetEndDateNull: setEndDateNull,
	}); err != nil {
		log.Errorf("ошибка выполнения: %v", err)
		// оборачиваем ошибку
		h.writeAppError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, nil)
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
	log := logger.Logger().WithFields(logger.LogOptions{
		Pkg:  "SubsHandler",
		Func: "DeleteSubscription",
		Ctx:  r.Context(),
	})

	uid, err := extrudeUidFromQuery(w, r)
	if err != nil {
		return
	}

	if err := h.container.DeleteSubscriptionHandler.Handle(r.Context(), commands.DeleteSubscriptionCommand{ID: uid}); err != nil {
		log.Errorf("ошибка выполнения: %v", err)
		// оборачиваем ошибку
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
// @Param start_from query string false "Start period from (MM-YYYY)"
// @Param start_to query string false "Start period to  (MM-YYYY)"
// @Param end_from query string false "End period from (MM-YYYY)"
// @Param end_to query string false "End period to  (MM-YYYY)"
// @Param nil_end query bool false "includes items with empty end_date(by default - true: if end_to != nil - false)"
// @Param page query int false "Page num can use 0 or 1 for first"
// @Param page_size query int false "Page Size / Limit"
// @Param order_by query string false "Sorting field name"
// @Param direction query string false "Sorting direction use 'asc'(default) or 'desc'"
// @Success 200 {array} Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions [get]
func (h *SubsHandler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	log := logger.Logger().WithFields(logger.LogOptions{
		Pkg:  "SubsHandler",
		Func: "ListSubscriptions",
		Ctx:  r.Context(),
	})

	var req SubscriptionQueryRequest
	if !utils.ParseQuery(w, r, &req) {
		return
	}

	//собираем квери
	// Парсим optional dates
	sfD, err := parseOptionalDate(w, req.StartFrom)
	if err != nil {
		return
	}
	stD, err := parseOptionalDate(w, req.StartTo)
	if err != nil {
		return
	}

	efD, err := parseOptionalDate(w, req.EndFrom)
	if err != nil {
		return
	}
	etD, err := parseOptionalDate(w, req.EndTo)
	if err != nil {
		return
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
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		StartFrom:   sfD,
		StartTo:     stD,
		EndFrom:     efD,
		EndTo:       etD,
		WithNilEnd:  req.NilEnd,

		Pagination: pagination,
		Sorting:    sorting,
	})
	if err != nil {
		log.Errorf("ошибка выполнения: %v", err)
		// оборачиваем ошибку
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
// @Param start_from query string false "Start period from (MM-YYYY)"
// @Param start_to query string false "Start period to  (MM-YYYY)"
// @Param end_from query string false "End period from (MM-YYYY)"
// @Param end_to query string false "End period to  (MM-YYYY)"
// @Param nil_end query bool false "includes items with empty end_date(by default - true: if end_to != nil - false)"
// @Success 200 {object} TotalCostResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/total [get]
func (h *SubsHandler) GetTotalCost(w http.ResponseWriter, r *http.Request) {
	log := logger.Logger().WithFields(logger.LogOptions{
		Pkg:  "SubsHandler",
		Func: "GetTotalCost",
		Ctx:  r.Context(),
	})

	var req TotalCostRequest
	if !utils.ParseQuery(w, r, &req) {
		return
	}

	// Парсим optional dates
	sfD, err := parseOptionalDate(w, req.StartFrom)
	if err != nil {
		return
	}
	stD, err := parseOptionalDate(w, req.StartTo)
	if err != nil {
		return
	}

	efD, err := parseOptionalDate(w, req.EndFrom)
	if err != nil {
		return
	}
	etD, err := parseOptionalDate(w, req.EndTo)
	if err != nil {
		return
	}

	result, err := h.container.TotalCostHandler.Handle(r.Context(), queries.TotalCostQuery{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		StartFrom:   sfD,
		StartTo:     stD,
		EndFrom:     efD,
		EndTo:       etD,
		WithNilEnd:  req.NilEnd,
	})
	if err != nil {
		log.Errorf("ошибка выполнения: %v", err)
		// оборачиваем ошибку
		h.writeAppError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, TotalCostResponse{Total: result})
}

package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/end1essrage/efmob-tz/pkg/common/interfaces/http/utils"
	"github.com/end1essrage/efmob-tz/pkg/common/logger"
	app "github.com/end1essrage/efmob-tz/pkg/subs/application"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

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

	if err := json.NewEncoder(w).Encode(ErrorResponse{
		Error: msg,
		Code:  appErr.Code,
	}); err != nil {
		logger.Logger().Log("AuthHandler", "writeAppError").Error(err)
	}
}

func extrudeUidFromQuery(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	id := chi.URLParam(r, "id")

	uid, err := uuid.Parse(id)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "invalid id format, should be uuid",
			Code:  "INVALID_QUERY",
		})

		return uuid.Nil, err
	}

	return uid, nil
}

func formatOptionalDate(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("01-2006")
	return &s
}

func formatDate(t time.Time) string {
	return t.Format("01-2006")
}

func parseOptionalDate(w http.ResponseWriter, s *string) (*time.Time, error) {
	if s == nil {
		return nil, nil
	}
	t, err := time.Parse("01-2006", *s)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
			Code:  "INVALID_DATE",
		})
		return &t, err
	}
	return &t, nil
}

func parseDate(w http.ResponseWriter, s string) (time.Time, error) {
	t, err := time.Parse("01-2006", s)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
			Code:  "INVALID_DATE",
		})
		return t, err
	}
	return t, nil
}

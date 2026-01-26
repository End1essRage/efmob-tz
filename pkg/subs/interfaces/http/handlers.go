package http

import (
	"net/http"

	utils "github.com/end1essrage/efmob-tz/pkg/common/interfaces/http/utils"
)

// Test godoc
// @Summary Test
// @Description SimpleTest
// @Tags subs
// #@Accept json
// #@Produce json
// #@Param request body RegisterRequest true "Registration data"
// #@Success 201 {object} RegisterResponse "User created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /test [get]
func (h *SubsHandler) Test(w http.ResponseWriter, r *http.Request) {

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "OKED",
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
*/

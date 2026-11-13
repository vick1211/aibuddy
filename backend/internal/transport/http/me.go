package httptransport

import (
	"aibuddy/internal/auth"
	"net/http"

	"github.com/google/uuid"
)

type meResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string `json:"role"`
	Team   string `json:"team"`
}

func Me(w http.ResponseWriter, r *http.Request) {
	id, ok := auth.IdentityFrom(r.Context())
	if !ok {
		unauthorized(w)
		return
	}
	writeJSON(w, http.StatusOK, meResponse{UserID: id.UserID, Role: string(id.Role), Team: id.Team})
}
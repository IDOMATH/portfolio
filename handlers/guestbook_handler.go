package handlers

import (
	"github.com/IDOMATH/portfolio/db"
	"github.com/IDOMATH/portfolio/render"
	"github.com/IDOMATH/portfolio/types"
	"github.com/IDOMATH/portfolio/util"
	"net/http"
)

type GuestbookHandler struct {
	guestbookStore db.PostgresGuestbookStore
}

func NewGuestbookHandler(guestbookStore db.PostgresGuestbookStore) *GuestbookHandler {
	return &GuestbookHandler{
		guestbookStore: guestbookStore,
	}
}

func (h *GuestbookHandler) HandleGetApprovedGuestbookSignatures() http.HandlerFunc {

	signatures, err := h.guestbookStore.GetApprovedGuestbookSignatures()
	objects := make(map[string]interface{})
	objects["signatures"] = signatures

	if err != nil {
		return func(w http.ResponseWriter, r *http.Request) {
			util.WriteError(w, http.StatusInternalServerError, err)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		render.Template(w, r, "guestbook.go.html", &types.TemplateData{
			PageTitle: "Guestbook",
			ObjectMap: objects,
		})
	}
}

func (h *GuestbookHandler) HandlePostGuestbookSignature() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (h *GuestbookHandler) HandleGetAllGuestbookSignature() http.HandlerFunc {

	signatures, err := h.guestbookStore.GetAllGuestbookSignatures()
	objects := make(map[string]interface{})
	objects["signatures"] = signatures

	if err != nil {
		return func(w http.ResponseWriter, r *http.Request) {
			util.WriteError(w, http.StatusInternalServerError, err)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		render.Template(w, r, "guestbook_admin.go.html", &types.TemplateData{
			PageTitle: "Guestbook Admin",
			ObjectMap: objects,
		})
	}
}

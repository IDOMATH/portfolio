package handlers

import (
	"github.com/IDOMATH/portfolio/db"
	"github.com/IDOMATH/portfolio/render"
	"github.com/IDOMATH/portfolio/types"
	"github.com/IDOMATH/portfolio/util"
	"net/http"
	"strconv"
	"time"
)

type GuestbookHandler struct {
	guestbookStore db.PostgresGuestbookStore
}

func NewGuestbookHandler(guestbookStore db.PostgresGuestbookStore) *GuestbookHandler {
	return &GuestbookHandler{
		guestbookStore: guestbookStore,
	}
}

func (h *GuestbookHandler) HandleGetApprovedGuestbookSignatures(w http.ResponseWriter, r *http.Request) {
	signatures, err := h.guestbookStore.GetApprovedGuestbookSignatures()
	objects := make(map[string]interface{})
	objects["signatures"] = signatures

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	render.Template(w, r, "guestbook.go.html", &types.TemplateData{
		PageTitle: "Guestbook",
		ObjectMap: objects,
	})

}

func (h *GuestbookHandler) HandlePostGuestbookSignature(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		name := r.FormValue("name")

		signature := types.GuestbookSignature{
			Name:       name,
			IsApproved: false,
			CreatedAt:  time.Now(),
		}
		_, err := h.guestbookStore.InsertGuestbookSignature(signature)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		objects := make(map[string]interface{})
		objects["signed"] = name

		render.Template(w, r, "guestbook-signed-successfully.go.html", &types.TemplateData{
			PageTitle: "Guestbook Signed",
			ObjectMap: objects,
		})
	}
}

func (h *GuestbookHandler) HandleGetAllGuestbookSignature(w http.ResponseWriter, r *http.Request) {
	signatures, err := h.guestbookStore.GetAllGuestbookSignatures()
	objects := make(map[string]interface{})
	objects["signatures"] = signatures

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	render.Template(w, r, "guestbook-admin.go.html", &types.TemplateData{
		PageTitle: "Guestbook Admin",
		ObjectMap: objects,
	})

}

func (h *GuestbookHandler) HandleApproveGuestbookSignature(w http.ResponseWriter, r *http.Request) {
	reqId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = h.guestbookStore.ApproveGuestbookSignature(reqId)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

}

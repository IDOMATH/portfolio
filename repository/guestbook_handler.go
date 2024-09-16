package repository

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/IDOMATH/portfolio/db"
	"github.com/IDOMATH/portfolio/render"
	"github.com/IDOMATH/portfolio/types"
	"github.com/IDOMATH/portfolio/util"
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

	render.Template(w, r, "guestbook.go.html", &render.TemplateData{
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

		render.Template(w, r, "guestbook-signed-successfully.go.html", &render.TemplateData{
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

	render.Template(w, r, "guestbook-admin.go.html", &render.TemplateData{
		PageTitle: "Guestbook Admin",
		ObjectMap: objects,
	})

}

func (h *GuestbookHandler) HandleApproveGuestbookSignature(w http.ResponseWriter, r *http.Request) {
	reqId, err := strconv.Atoi(r.FormValue("id"))
	fmt.Println("approving id: ", reqId)
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

func (h *GuestbookHandler) HandleDenyGuestbookSignature(w http.ResponseWriter, r *http.Request) {
	//var ids int[]
	reqId, err := strconv.Atoi(r.FormValue("id"))
	//for _, reqId := range r.Form {
	//	id, err := strconv.Atoi(reqId)
	//	ids = append(ids, id)
	//}
	//reqId := r.FormValue("ids")
	//var err error
	//fmt.Printf("ID: %v", id)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	fmt.Println("ID: ", reqId)
	//err = h.guestbookStore.DenyGuestbookSignature(reqId)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

}

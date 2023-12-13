package handlers

import (
	"github.com/IDOMATH/portfolio/db"
	"github.com/IDOMATH/portfolio/util"
	"net/http"
)

type FitnessHandler struct {
	fitnessStore db.PostgresFitnessStore
}

func NewFintessHandler(store db.PostgresFitnessStore) *FitnessHandler {
	return &FitnessHandler{
		fitnessStore: store,
	}
}

func (h *FitnessHandler) HandleGetFitness(w http.ResponseWriter, r *http.Request) {
	fitnessEntries, err := h.fitnessStore.GetAllFitnessRecaps()
	objects := make(map[string]interface{})
	objects["fitnessEntries"] = fitnessEntries

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
	}
}

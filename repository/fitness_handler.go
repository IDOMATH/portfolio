package repository

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/IDOMATH/portfolio/db"
	"github.com/IDOMATH/portfolio/render"
	"github.com/IDOMATH/portfolio/types"
	"github.com/IDOMATH/portfolio/util"
)

type FitnessHandler struct {
	fitnessStore db.PostgresFitnessStore
}

func NewFitnessHandler(store db.PostgresFitnessStore) *FitnessHandler {
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

func (h *FitnessHandler) HandleGetFitnessForm(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "fitness-form.go.html",
		&render.TemplateData{PageTitle: "Fitness"})
}

func (h *FitnessHandler) HandlePostFitnessForm(w http.ResponseWriter, r *http.Request) {
	// TODO: figure out how this will be input on the front end and translated here
	weight, err := strconv.Atoi(r.PostForm.Get("weight"))
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// TODO: figure out how this will be input on the front end and translated here
	distance, err := strconv.Atoi(r.PostForm.Get("distance"))
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	date := r.PostForm.Get("date")
	ymdString := strings.Split(date, "-")
	var ymd []int
	for _, str := range ymdString {
		i, err := strconv.Atoi(str)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		ymd = append(ymd, i)
	}

	recap := types.FitnessRecap{
		HundredthsOfAMile: distance,
		TenthsOfAPound:    weight,
		Date:              time.Date(ymd[0], time.Month(ymd[1]), ymd[2], 0, 0, 0, 0, time.UTC)}
	_, err = h.fitnessStore.InsertFitnessRecap(recap)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

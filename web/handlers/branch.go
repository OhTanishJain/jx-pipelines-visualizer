package handlers

import (
	"net/http"
	"strings"

	visualizer "github.com/jenkins-x/jx-pipelines-visualizer"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/render"
)

type BranchHandler struct {
	Store  *visualizer.Store
	Render *render.Render
	Logger *logrus.Logger
}

func (h *BranchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	owner := vars["owner"]
	repository := vars["repo"]
	branch := vars["branch"]
	if strings.HasPrefix(branch, "pr-") {
		branch = strings.ToUpper(branch)
	}

	pipelines, err := h.Store.Query(visualizer.Query{
		Owner:      owner,
		Repository: repository,
		Branch:     branch,
		Query:      r.URL.Query().Get("q"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.Render.HTML(w, http.StatusOK, "home", struct {
		Owner      string
		Repository string
		Branch     string
		Query      string
		Pipelines  *visualizer.Pipelines
	}{
		owner,
		repository,
		branch,
		r.URL.Query().Get("q"),
		pipelines,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

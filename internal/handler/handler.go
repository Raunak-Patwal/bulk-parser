package handler

import (
	"encoding/json"
	"net/http"

	"bulk-parser/internal/limiter"
	"bulk-parser/internal/model"
	"bulk-parser/internal/orchestrator"
	"sync/atomic"
)

var nextJobID int64

type Handler struct {
	limiter      *limiter.RequestLimiter
	orchestrator *orchestrator.Orchestrator
}

func NewHandler(
	limiter *limiter.RequestLimiter,
	orchestrator *orchestrator.Orchestrator,
) *Handler {
	return &Handler{
		limiter:      limiter,
		orchestrator: orchestrator,
	}
}

func (h *Handler) Process(w http.ResponseWriter, r *http.Request) {

	h.limiter.Acquire()
	defer h.limiter.Release()

	var request model.BulkRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	jobs := make([]model.Job, 0, len(request.Sources))

	for _, source := range request.Sources {
    jobs = append(jobs, model.Job{
        ID:       int(atomic.AddInt64(&nextJobID, 1)),
        Language: source.Language,
        Data:     source.Data,
    })
}

	results := h.orchestrator.Process(jobs)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(results)
}

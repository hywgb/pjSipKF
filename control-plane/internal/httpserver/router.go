package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"

	"github.com/hywgb/pjSipKF/control-plane/internal/config"
	"github.com/hywgb/pjSipKF/control-plane/internal/mediacore"
	"github.com/hywgb/pjSipKF/control-plane/internal/registrar"
)

type apiServer struct {
	logger   *zap.Logger
	cfg      config.Config
	reg      registrar.Service
	mcClient mediacore.Client
}

func RegisterRoutes(r *chi.Mux, logger *zap.Logger, cfg config.Config) {
	mc, err := mediacore.NewClientFromConfig(cfg)
	if err != nil {
		logger.Fatal("failed to init mediacore client", zap.Error(err), zap.String("uds", cfg.MediaCoreUDS))
	}
	server := &apiServer{
		logger:   logger,
		cfg:      cfg,
		reg:      registrar.NewInMemory(),
		mcClient: mc,
	}

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Post("/v1/register", server.handleRegister)
	r.Delete("/v1/register", server.handleDeregister)
	r.Get("/v1/lookup", server.handleLookup)

	r.Post("/v1/sessions", server.handleCreateSession)
	r.Delete("/v1/sessions/{id}", server.handleTerminateSession)
}

type registerRequest struct {
	User    string `json:"user"`
	Contact string `json:"contact"`
	Expires int    `json:"expires"`
}

func (s *apiServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.User == "" || req.Contact == "" || req.Expires <= 0 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := s.reg.Register(r.Context(), req.User, req.Contact, req.Expires); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *apiServer) handleDeregister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.User == "" || req.Contact == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := s.reg.Deregister(r.Context(), req.User, req.Contact); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *apiServer) handleLookup(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if user == "" {
		http.Error(w, "missing user", http.StatusBadRequest)
		return
	}
	contacts, err := s.reg.Lookup(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"user": user, "contacts": contacts})
}

type createSessionRequest struct {
	SDPOffer string            `json:"sdp_offer"`
	Meta     map[string]string `json:"metadata"`
}

func (s *apiServer) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	var req createSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, answer, err := s.mcClient.CreateSession(r.Context(), req.SDPOffer, req.Meta)
	if err != nil {
		st, _ := status.FromError(err)
		s.logger.Error("CreateSession failed", zap.Error(err), zap.String("grpc_code", st.Code().String()), zap.String("grpc_msg", st.Message()))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"session_id": id, "sdp_answer": answer})
}

func (s *apiServer) handleTerminateSession(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	if err := s.mcClient.TerminateSession(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
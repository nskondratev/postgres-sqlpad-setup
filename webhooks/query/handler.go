package query

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nskondratev/postgres-sqlpad-setup/webhooks/sqlpad"
)

type payload struct {
	Query struct {
		ID   string `json:"id"`
		User struct {
			Email string `json:"email"`
		} `json:"use"`
	} `json:"query"`
}

// Handler обрабатывает вебхуки на создание запроса и удаляет его
type Handler struct {
	cli        *sqlpad.Client
	adminEmail string
}

// NewHandler ...
func NewHandler(cli *sqlpad.Client, adminEmail string) *Handler {
	return &Handler{cli: cli, adminEmail: adminEmail}
}

// ServeHTTP реализует интерфейс http.Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[query_created_handler] method is not post: %s", r.Method)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	var payload payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("[query_created_handler] failed to decode payload: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if payload.Query.User.Email == h.adminEmail {
		log.Printf("[query_created_handler] admin created query. Skip it")
		w.WriteHeader(http.StatusOK)
		return
	}

	query, err := h.cli.GetQuery(r.Context(), payload.Query.ID)
	if err != nil {
		log.Printf("[query_created_handler] failed to get query by id: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(query.ACL) > 0 {
		err = h.cli.ClearQueryACL(r.Context(), payload.Query.ID)
		if err != nil {
			log.Printf("[query_created_handler] failed to clear query acl: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

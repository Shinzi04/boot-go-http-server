package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (c *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type paramaters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := paramaters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while decoding response body", err)
		return
	}

	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Need to provide body email!", nil)
		return
	}

	user, err := c.db.CreateUser(context.Background(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while creating user", err)
		return
	}

	type userResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	respondWithJSON(w, 201, userResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}

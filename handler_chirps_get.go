package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while fetching users data", err)
		return
	}

	response := make([]Chrip, len(chirps))
	for i, chrip := range chirps {
		response[i] = Chrip{
			ID:        chrip.ID,
			CreatedAt: chrip.CreatedAt,
			UpdatedAt: chrip.UpdatedAt,
			Body:      chrip.Body,
			UserID:    chrip.UserID,
		}
	}

	data, err := json.Marshal(response)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failing to marshal response", err)
		return
	}

	w.Write(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) handlerChripsGet(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(context.Background(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error while fetching data", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chrip{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

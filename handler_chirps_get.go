package main

import (
	"encoding/json"
	"gohttp/internal/database"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	sortBy := r.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error

	if s != "" {
		authorID, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
			return
		}
		chirps, err = cfg.db.GetChirpByUserID(r.Context(), authorID)
	} else {
		chirps, err = cfg.db.GetChirps(r.Context())
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while fetching chirps data", err)
		return
	}

	response := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		response[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	// Sorting in memory
	switch sortBy {
	case "desc":
		sort.Slice(response, func(i, j int) bool {
			return response[i].CreatedAt.After(response[j].CreatedAt)
		})
	default: // "asc" or empty
		sort.Slice(response, func(i, j int) bool {
			return response[i].CreatedAt.Before(response[j].CreatedAt)
		})
	}

	data, err := json.Marshal(response)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failing to marshal response", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error while fetching data", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

package main

import (
	"context"
	"gohttp/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not Authorization Bearer Detected", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(context.Background(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Refresh Token", nil)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expirationTimeToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not Authorization Bearer Detected", err)
		return
	}

	if err := cfg.db.RevokeRefreshToken(context.Background(), refreshToken); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Refresh Token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

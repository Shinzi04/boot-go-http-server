package main

import (
	"context"
	"encoding/json"
	"gohttp/internal/auth"
	"gohttp/internal/database"
	"net/http"
	"time"
)

const (
	// NOTE: DEFAULT TOKEN EXPIRATION DATE
	expirationTimeToken = time.Hour

	// NOTE: DEFAULT REFRESH TOKEN EXPIRATION DATE
	expirationTimeRefreshToken = time.Hour * 24 * 60
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type paramaters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := paramaters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	refreshToken, err := cfg.db.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		Token:     auth.MakeRefreshToken(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(expirationTimeRefreshToken).UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while creating refresh token", err)
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expirationTimeToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while creating access token", err)
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken.Token,
	})
}

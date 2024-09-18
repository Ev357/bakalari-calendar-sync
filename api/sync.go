package handler

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/Ev357/bakalari-calendar-sync/utils"
	"github.com/joho/godotenv"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	godotenv.Overload()

	config, err := utils.GetConfig()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authHeader := r.Header.Get("Authorization")
	expectedAuthHeader := fmt.Sprintf("Bearer %s", config.CronSecret)

	if subtle.ConstantTimeCompare([]byte(authHeader), []byte(expectedAuthHeader)) != 1 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = utils.Sync(config)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte("ok"))
}

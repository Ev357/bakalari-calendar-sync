package handler

import (
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

	if r.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", config.CronSecret) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = utils.Sync(config)

	if err != nil {
		panic(err)
	}

	w.Write([]byte("ok"))
}

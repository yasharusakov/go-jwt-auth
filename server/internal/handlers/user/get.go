package user

import (
	"encoding/json"
	"net/http"
	"server/internal/database/postgresql"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := postgresql.Pool.Query(r.Context(), "SELECT id, email FROM users")
	if err != nil {
		http.Error(w, "error getting users", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var users []map[string]string

	for rows.Next() {
		var id, email string
		err = rows.Scan(&id, &email)
		if err != nil {
			http.Error(w, "error scanning row", http.StatusInternalServerError)
			return
		}
		users = append(users, map[string]string{"id": id, "email": email})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

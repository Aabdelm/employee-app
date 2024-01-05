package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	employeedb "github.com/Aabdelm/employee-app/database"
	"github.com/Aabdelm/employee-app/security"
)

type SecurityHandler struct {
	L   *log.Logger
	Sec employeedb.Sec
}

func NewSecurityHandler(l *log.Logger, sec employeedb.Sec) *SecurityHandler {
	return &SecurityHandler{
		L:   l,
		Sec: sec,
	}
}

func (sh *SecurityHandler) PostInfo(rw http.ResponseWriter, r *http.Request) {
	sh.L.Printf("[INFO] Starting function PostInfo \n")

	var err error

	user := employeedb.NewUser()
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(user); err != nil {
		sh.L.Printf("[ERROR] Failed to decode JSON. Error %s", err)
		http.Error(rw, "Failed to decode JSON. This might be due to a corrupt JSON.", http.StatusBadRequest)
		return
	}

	hashedInfo, err := sh.Sec.FetchInfo(user.UserName)
	if err != nil {
		sh.L.Printf("[ERROR] Failed to retrieve user from Database. Error %s", err)
		if err == sql.ErrNoRows {
			http.Error(rw, "User does not exist", http.StatusUnauthorized)
		}
		http.Error(rw, "Failed to retrieve user information", http.StatusInternalServerError)
	}

	err = security.VerifyPass(user.UserPass, hashedInfo, sh.L)

	if err != nil {
		sh.L.Printf("[ERROR] Failed to verify password. Error %s", err)
		if err == security.ErrWrongPass {
			http.Error(rw, "Wrong password", http.StatusUnauthorized)
			return
		}
		http.Error(rw, "Failed to verify password", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(rw, "You're logged in :D")

}

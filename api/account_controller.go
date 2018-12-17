package api

import (
	"encoding/json"
	"github.com/flohero/Spongebot/database/model"
	"net/http"
)

func (c *Controller) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var temp *model.Account
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		badRequest(w, err)
		return
	}
	err, temp = c.persistence.CreateAccount(temp)
	if err != nil {
		badRequest(w, err)
		return
	}
	created(w)
	writeJson(w, temp)
}

func (c *Controller) Authenticate(w http.ResponseWriter, r *http.Request) {
	account := &model.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		badRequest(w, err)
		return
	}

	err, account = c.persistence.Login(account.Email, account.Password)
	if err != nil {
		forbidden(w, err)
		return
	}
	writeJson(w, account)
}
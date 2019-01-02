package database

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/flohero/Spongebot/database/model"
	"golang.org/x/crypto/bcrypt"
	"os"
)

const JWT_PASSWORD string = "JWT_PASSWORD"

func (p *Persistence) IsValid(acc *model.Account) (bool, error) {
	if p.FindByUsername(acc.Username).Username != "" {
		return false, errors.New("")
	}
	return true, nil
}

func (p *Persistence) CreateAccount(acc *model.Account) (error, *model.Account) {
	if ok, err := p.IsValid(acc); !ok {
		return err, nil
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(acc.Password), bcrypt.DefaultCost)
	acc.Password = string(hashedPassword)
	p.db.Create(acc)
	if acc.Id <= 0 {
		return errors.New("Failed to create account"), nil
	}
	tk := &model.Token{UserId: acc.Id}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv(JWT_PASSWORD)))
	acc.Token = tokenString

	acc.Password = "" //delete password
	return nil, acc
}

func (p *Persistence) Login(username, password string) (error, *model.Account) {

	account := &model.Account{}
	if account = p.FindByUsername(username); account.Username == "" {
		return errors.New("Username not found"), nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return errors.New("Invalid login credentials. Please try again"), nil
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	tk := &model.Token{UserId: account.Id}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv(JWT_PASSWORD)))
	account.Token = tokenString //Store the token in the response

	return nil, account
}

func (p *Persistence) FindByUsername(username string) (acc *model.Account) {
	acc = &model.Account{}
	p.db.Where(&model.Account{Username: username}).First(acc)
	acc.Token = ""
	return acc
}

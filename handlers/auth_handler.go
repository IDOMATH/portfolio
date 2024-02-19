package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/IDOMATH/portfolio/db"
	"github.com/IDOMATH/portfolio/types"
	"github.com/IDOMATH/portfolio/util"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/mail"
	"time"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

func (h *AuthHandler) HandleUserSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// TODO: write a JSON error and return a status code
		return
	}
	err := r.ParseForm()
	if err != nil {
		// TODO: Render a template with an error message
		w.Write([]byte(err.Error()))
		return
	}

	// TODO: some server side validation

	username := r.PostForm.Get("username")
	email := r.PostForm.Get("email")
	// TODO: Look up some password best practices
	password := r.PostForm.Get("password")

	_, err = mail.ParseAddress(email)
	if err != nil {
		fmt.Println("invalid email address", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		// TODO: Render a template with an error message
		w.Write([]byte(err.Error()))
		return
	}

	user := &types.User{
		Username: username,
		Email:    email,
		Password: string(passwordHash),
	}
	user.CreatedAt = time.Now()
	//blog.ImageName :=
	insertedUser, err := h.userStore.InsertUser(context.TODO(), user)
	if err != nil {
		// TODO: Render a template with an error message
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println(insertedUser)
}

func (h *AuthHandler) HandleUserLogIn(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		// TODO: write a JSON error and return a status code
		return
	}
	err := r.ParseForm()
	if err != nil {
		// TODO: Render a template with an error message
		w.Write([]byte(err.Error()))
		return
	}

	// TODO: some server side validation

	username := r.PostForm.Get("username")
	// TODO: Look up some password best practices
	password := r.PostForm.Get("password")

	if !util.IsValidPassword(password) {
		// TODO: Render a template with an error message
		err = errors.New("password does not meet requirements")
		w.Write([]byte(err.Error()))
		return
	}

	user, err := h.userStore.GetUser(context.Background(), username)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		// TODO: Render a template with an error message for failed login
		w.Write([]byte(err.Error()))
		return
	}

	// TODO: add entry to session and token DB
}

package handlers

import (
	"context"
	"fmt"
	"github.com/IDOMATH/portfolio/db"
	"github.com/IDOMATH/portfolio/types"
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
		fmt.Println(err)
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
		fmt.Println(err)
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
		fmt.Println(err)
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
		fmt.Println(err)
		return
	}

	// TODO: some server side validation

	username := r.PostForm.Get("username")
	// TODO: Look up some password best practices
	password := r.PostForm.Get("password")

	user, err := h.userStore.GetUser(context.Background(), username)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		// TODO: Render a template with an error message for failed login
		fmt.Println("username or password incorrect", err)
		return
	}

	// TODO: add entry to session and token DB
}

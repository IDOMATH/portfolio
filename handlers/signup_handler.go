package handlers

import (
	"context"
	"fmt"
	"github.com/IDOMATH/portfolio/db"
	"github.com/IDOMATH/portfolio/types"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlePostUser(w http.ResponseWriter, r *http.Request) {
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

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		// TODO: Render a template with an error message
		fmt.Println(err)
		return
	}

	// TODO: store the uploaded file somewhere.
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

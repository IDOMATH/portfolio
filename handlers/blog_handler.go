package handlers

import (
	"context"
	"fmt"
	"github.com/IDOMATH/portfolio/db"
	"github.com/IDOMATH/portfolio/render"
	"github.com/IDOMATH/portfolio/types"
	"github.com/IDOMATH/portfolio/util"
	"net/http"
	"strings"
	"time"
)

type BlogHandler struct {
	blogStore db.BlogStore
}

func NewBlogHandler(blogStore db.BlogStore) *BlogHandler {
	return &BlogHandler{
		blogStore: blogStore,
	}
}

func (h *BlogHandler) HandleBlog(w http.ResponseWriter, r *http.Request) {
	c, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	blogCards, err := h.blogStore.GetBlogCards(c)
	objects := make(map[string]interface{})
	objects["blog_posts"] = blogCards

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)

	} else {

		render.Template(w, r, "all-blogs.go.html", &types.TemplateData{
			PageTitle: "All Blogs",
			ObjectMap: objects,
		})
	}
}

func (h *BlogHandler) HandleGetBlogById(w http.ResponseWriter, r *http.Request) {
	// Path to get here is	 /blog/{id}
	// 					  [0]/[1] /[2]
	url := strings.Split(r.URL.Path, "/")
	fmt.Println(url)
	id := url[2]
	blog, err := h.blogStore.GetBlogById(context.Background(), id)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
	}
	if blog == nil {
		http.Redirect(w, r, "/blog", http.StatusTemporaryRedirect)
	}
	// TODO: Make a template for singular blogs
	w.Write([]byte(blog.Title))
}

func (h *BlogHandler) HandleNewBlog(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		render.Template(w, r, "blog-form.go.html", &types.TemplateData{PageTitle: "New Blog"})
	}
	if r.Method == "POST" {

		var blog types.BlogPost

		title := r.FormValue("title")
		// TODO: get the author from the logged in user.
		body := r.FormValue("body")
		imageName := r.FormValue("image")

		blog.Title = title
		blog.Body = body
		blog.ImageName = imageName
		blog.PublishedAt = time.Now()
		_, err := h.blogStore.InsertBlogPost(context.Background(), &blog)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, err)
		}

		// TODO: Make this actually load the page
		http.Redirect(w, r, "/blog", http.StatusCreated)
	}
}

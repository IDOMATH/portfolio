package handlers

import (
	"context"
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

func (h *BlogHandler) HandleGetBlogs(ctx context.Context) http.HandlerFunc {
	c, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	blogCards, err := h.blogStore.GetBlogCards(c)
	objects := make(map[string]interface{})
	objects["blog_posts"] = blogCards

	if err != nil {
		return func(w http.ResponseWriter, r *http.Request) {
			util.WriteError(w, http.StatusInternalServerError, err)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		render.Template(w, r, "all-blogs.go.html", &types.TemplateData{
			PageTitle: "All Blogs",
			ObjectMap: objects,
		})
	}
}

func (h *BlogHandler) HandleGetBlogById(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	id := url[1]
	blog, err := h.blogStore.GetBlogById(context.Background(), id)
	if err != nil {
		// TODO: make this return a handlerfunc
		util.WriteError(w, http.StatusInternalServerError, err)
	}

	w.Write([]byte(blog.Title))
}

func (h *BlogHandler) HandlePostBlog(w http.ResponseWriter, r *http.Request) {
	//pic, err := c.FormFile("thumbnail")
	//if err != nil {
	//	return err
	//}
	// TODO: store the uploaded file somewhere.
	var blog types.BlogPost

	title := r.PostForm.Get("title")
	// TODO: get the author from the logged in user.
	body := r.PostForm.Get("body")
	imageName := r.PostForm.Get("image")

	blog.Title = title
	blog.Body = body
	blog.ImageName = imageName
	blog.PublishedAt = time.Now()
	//blog.ImageName :=
	insertedBlog, err := h.blogStore.InsertBlogPost(context.Background(), &blog)
	if err != nil {
		// TODO: make this return a handlerfunc
		util.WriteError(w, http.StatusInternalServerError, err)
	}

	w.Write([]byte(insertedBlog.Title))
}

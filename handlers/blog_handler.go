package handlers

import (
	"context"
	"github.com/IDOMATH/portfolio/db"
	"github.com/IDOMATH/portfolio/render"
	"github.com/IDOMATH/portfolio/types"
	"github.com/IDOMATH/portfolio/util"
	"github.com/gofiber/fiber/v2"
	"net/http"
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

//func (h *BlogHandler) HandleGetBlogById(w http.ResponseWriter, r *http.Request, ctx context.Context) error {
//	fmt.Println(ctx.Params("id"))
//	blog, err := h.blogStore.GetBlogById(c.Context(), c.Params("id"))
//	if err != nil {
//		return err
//	}
//	return c.JSON(blog)
//}

func (h *BlogHandler) HandlePostBlog(c *fiber.Ctx) error {
	//pic, err := c.FormFile("thumbnail")
	//if err != nil {
	//	return err
	//}
	// TODO: store the uploaded file somewhere.
	var blog types.BlogPost
	if err := c.BodyParser(&blog); err != nil {
		return err
	}
	blog.PublishedAt = time.Now()
	//blog.ImageName :=
	insertedBlog, err := h.blogStore.InsertBlogPost(c.Context(), &blog)
	if err != nil {
		return err
	}
	return c.JSON(insertedBlog)
}

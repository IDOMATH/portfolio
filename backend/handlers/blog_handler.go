package handlers

import (
	"fmt"
	"github.com/IDOMATH/portfolio/db"
	"github.com/IDOMATH/portfolio/types"
	"github.com/gofiber/fiber/v2"
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

func (h *BlogHandler) HandleGetBlogs(c *fiber.Ctx) error {
	//util.EnableCors(c)
	blogCards, err := h.blogStore.GetBlogCards(c.Context())
	if err != nil {
		return err
	}
	return c.Render("all-blogs", fiber.Map{
		"PageTitle": "All Blogs",
		"BlogPosts": blogCards,
	}, "layouts/base")
}

func (h *BlogHandler) HandleGetBlogById(c *fiber.Ctx) error {
	fmt.Println(c.Params("id"))
	blog, err := h.blogStore.GetBlogById(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	return c.JSON(blog)
}

func (h *BlogHandler) HandlePostBlog(c *fiber.Ctx) error {
	//util.EnableCors(c)
	c.Body()
	var blog types.BlogPost
	if err := c.BodyParser(&blog); err != nil {
		return err
	}
	blog.PublishedAt = time.Now()
	insertedBlog, err := h.blogStore.InsertBlogPost(c.Context(), &blog)
	if err != nil {
		return err
	}
	return c.JSON(insertedBlog)
}

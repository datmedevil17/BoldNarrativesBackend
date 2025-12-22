package blog

import (
	"github.com/datmedevil17/BoldNarrativesBackend/internal/services/blog"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service   *blog.Service
	jwtSecret string
}

func NewHandler(service *blog.Service, jwtSecret string) *Handler {
	return &Handler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

func (h *Handler) CreateBlog(c *gin.Context) {

}
func (h *Handler) GetBlogById(c *gin.Context) {

}
func (h *Handler) UpdateBlog(c *gin.Context) {

}
func (h *Handler) DeleteBlog(c *gin.Context) {

}
func (h *Handler) GetTotalCount(c *gin.Context) {

}
func (h *Handler) SortByTime(c *gin.Context) {

}
func (h *Handler) SortByViews(c *gin.Context) {

}
func (h *Handler) GetTrending(c *gin.Context) {

}
func (h *Handler) IncrementViews(c *gin.Context) {

}
func (h *Handler) CheckVote(c *gin.Context) {

}
func (h *Handler) ToggleVote(c *gin.Context) {

}
func (h *Handler) GetCommentsByBlogId(c *gin.Context) {

}
func (h *Handler) CreateComment(c *gin.Context) {

}
func (h *Handler) DeleteComment(c *gin.Context) {

}


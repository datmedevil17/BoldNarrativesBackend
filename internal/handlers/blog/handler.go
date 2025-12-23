package blog

import (
	"net/http"
	"strconv"

	"github.com/datmedevil17/BoldNarrativesBackend/internal/services/blog"
	"github.com/datmedevil17/BoldNarrativesBackend/internal/utils"
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
	var req CreateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	userId, ok := c.Get("userID")
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error occured in blog creation")
		return
	}
	authorId := userId.(uint)
	blog, err := h.service.CreateBlog(authorId, req.Title, req.Content, req.Genre)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error occured in blog creation")
		return
	}
	c.JSON(http.StatusOK, blog)

}
func (h *Handler) GetBlogById(c *gin.Context) {
	blogId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid blog id")
		return
	}
	blog, err := h.service.GetBlogById(uint(blogId))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Blog not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"blog": blog,
	})

}
func (h *Handler) UpdateBlog(c *gin.Context) {
	blogId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid blog id")
		return
	}
	var req UpdateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	userID, _ := c.Get("userID")
	currentUserID := userID.(uint)
	blog, err := h.service.UpdateBlog(uint(blogId), req.Title, req.Content, req.Genre, currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error occured in blog update")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"blog": blog,
	})
}
func (h *Handler) DeleteBlog(c *gin.Context) {
	blogId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid blog id")
		return
	}
	userID, _ := c.Get("userID")
	currentUserID := userID.(uint)
	err = h.service.DeleteBlog(uint(blogId), currentUserID)
	if err != nil {
		if err.Error() == "unauthorized: you can only delete your own blogs" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error occured in blog deletion")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Blog deleted successfully",
	})
}
func (h *Handler) GetTotalCount(c *gin.Context) {
	var req FilterRequest
	if err:=c.ShouldBindJSON(&req);err!=nil{
		req=FilterRequest{}		
	}
	opts:=blog.Filter{
		Genre:    req.Genre,
		AuthorID: req.AuthorID,
		Search:   req.Search,
	}

	total,err:=h.service.GetBlogsCount(opts)
	if err!=nil{
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error occured in getting blogs count")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total": total,
	})
	
	

}
func (h *Handler) SortByTime(c *gin.Context) {
	sortOrder := c.Query("sortOrder")
	var req FilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
	
		return
	}
	opts := blog.Filter{
		Genre:    req.Genre,
		AuthorID: req.AuthorID,
		Search:   req.Search,
		Skip:     req.Skip,
		Limit:    10,
	}
	ascending := sortOrder == "asc"

	blogs, err := h.service.GetBlogsSortedByTime(opts, ascending)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error occured in getting blogs")
		return
	}
	c.JSON(http.StatusOK, blogs)

}
func (h *Handler) SortByViews(c *gin.Context) {
	var req FilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		return
	}
	opts := blog.Filter{
		Genre:    req.Genre,
		AuthorID: req.AuthorID,
		Search:   req.Search,
		Skip:     req.Skip,
		Limit:    10,
	}

	blogs, err := h.service.GetBlogsSortedByViews(opts)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error occured in getting blogs")
		return
	}
	c.JSON(http.StatusOK, blogs)

}
func (h *Handler) GetTrending(c *gin.Context) {
	blogs, err := h.service.GetTrendingBlogs()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error occured in getting trending blogs")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"blogs": blogs,
	})

}
func (h *Handler) IncrementViews(c *gin.Context) {
	var req ViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	err := h.service.IncrementViews(req.ID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error occured in incrementing views")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Views incremented successfully",
	})

}
func (h *Handler) CheckVote(c *gin.Context) {
	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	userID, _ := c.Get("userID")
	currentUserID := userID.(uint)
	vote, err := h.service.CheckVote(req.ID, currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error occured in checking vote")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"vote": vote,
	})

}
func (h *Handler) ToggleVote(c *gin.Context) {
	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	userID, _ := c.Get("userID")
	currentUserID := userID.(uint)
	vote, err := h.service.ToggleVote(req.ID, currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error occured in toggling vote")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"vote": vote,
	})

}
func (h *Handler) GetCommentsByBlogId(c *gin.Context) {

	blogId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid blog id")
		return
	}
	comments, err := h.service.GetCommentsByBlogId(uint(blogId))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error occured in getting comments")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}
func (h *Handler) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	userID, _ := c.Get("userID")
	currentUserID := userID.(uint)
	comment, err := h.service.CreateComment(req.BlogID, currentUserID, req.Comment)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error occured in creating comment")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"comment": comment,
	})

}
func (h *Handler) DeleteComment(c *gin.Context) {
	commentId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid comment id")
		return
	}
	userID, _ := c.Get("userID")
	currentUserID := userID.(uint)
	err = h.service.DeleteComment(uint(commentId), currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error occured in deleting comment")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})
}

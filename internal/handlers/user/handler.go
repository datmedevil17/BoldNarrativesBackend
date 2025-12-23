package user

import (
	"net/http"
	"strconv"

	"github.com/datmedevil17/BoldNarrativesBackend/internal/services/user"
	"github.com/datmedevil17/BoldNarrativesBackend/internal/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service   *user.Service
	jwtSecret string
}

func NewHandler(service *user.Service, jwtSecret string) *Handler {
	return &Handler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

func (h *Handler) SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Enter valid email")
	}
	user, err := h.service.CreateUser(req.Email, req.Name, req.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	token, err := utils.GenerateToken(user.Email, user.ID, h.jwtSecret)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error generating token")
		return
	}
	c.String(http.StatusOK, token)
}

func (h *Handler) SignIn(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Enter valid email")
	}
	user, err := h.service.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
	}
	token, err := utils.GenerateToken(user.Email, user.ID, h.jwtSecret)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error generating token")
		return
	}
	c.String(http.StatusOK, token)
}

func (h *Handler) GetUserById(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid userId")
		return
	}

	user, err := h.service.GetUserById(uint(userID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	c.JSON(http.StatusOK, user)

}

func (h *Handler) GetCurrentUserId(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid userId")
		return
	}
	user, err := h.service.GetUserById(uint(userID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}
	c.JSON(http.StatusOK, user)

}

func (h *Handler) ViewProfile(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid userId")
		return
	}
	user, err := h.service.GetUserProfile(uint(userID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}
	c.JSON(http.StatusOK, user)

}

func (h *Handler) GetProfile(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Error occured in profile.")
		return
	}
	user, err := h.service.GetUserProfile(uint(userID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}
	c.JSON(http.StatusOK, user)

}
func (h *Handler) CheckFollowStatus(c *gin.Context) {
	userId, _ := c.Get("userID")
	currentUserId := userId.(uint)

	var req FollowRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Enter valid user id")
	}
	if currentUserId == req.TargetUserIdParam {
		utils.ErrorResponse(c, http.StatusBadRequest, "You cannot follow yourself")
		return
	}
	followStatus, err := h.service.CheckIfFollowing(req.TargetUserIdParam, currentUserId)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
	}
	c.JSON(http.StatusOK, gin.H{"followStatus": followStatus})

}
func (h *Handler) FollowUser(c *gin.Context) {
	userId, _ := c.Get("userID")
	currentUserId := userId.(uint)

	var req FollowRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Enter valid user id")
	}
	err := h.service.FollowUser(req.TargetUserIdParam, currentUserId)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
	}
	c.JSON(http.StatusOK, gin.H{"message": "Followed successfully"})

}
func (h *Handler) UnFollowUser(c *gin.Context) {
	userId, _ := c.Get("userID")
	currentUserId := userId.(uint)

	var req FollowRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Enter valid user id")
	}
	err := h.service.UnFollowUser(req.TargetUserIdParam, currentUserId)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
	}
	c.JSON(http.StatusOK, gin.H{"message": "Unfollowed successfully"})
}

func (h *Handler) GetFollowers(c *gin.Context) {
	userId, _ := c.Get("userID")
	currentUserId := userId.(uint)
	followers, err := h.service.GetFollowers(currentUserId)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
	}
	c.JSON(http.StatusOK, followers)

}
func (h *Handler) GetFollowing(c *gin.Context) {
	userId, _ := c.Get("userID")
	currentUserId := userId.(uint)
	following, err := h.service.GetFollowing(currentUserId)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
	}
	c.JSON(http.StatusOK, following)

}

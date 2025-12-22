package user

import (
	"net/http"

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
		utils.ErrorResponse(c,http.StatusBadRequest,"Enter valid email")
	}
	user,err:=h.service.CreateUser(req.Email,req.Name,req.Password)
	if err!=nil{
		if err.Error() == "user already exists" {
            utils.ErrorResponse(c, http.StatusConflict, err.Error())
            return
        }
		utils.ErrorResponse(c,http.StatusInternalServerError,"Internal Server Error")
	}
	token,err:=utils.GenerateToken(user.Email,user.ID,h.jwtSecret)
	if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Error generating token")
        return
	}
	c.String(http.StatusOK,token)
}

func (h *Handler) SignIn(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c,http.StatusBadRequest,"Enter valid email")
	}
	user,err:=h.service.AuthenticateUser(req.Email,req.Password)
	if err!=nil{
		utils.ErrorResponse(c,http.StatusInternalServerError,"Internal Server Error")
	}
	token,err:=utils.GenerateToken(user.Email,user.ID,h.jwtSecret)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error generating token")
		return
	}
	c.String(http.StatusOK,token)
}

func (h *Handler)GetUserById(c *gin.Context){
	
}

func (h *Handler)GetCurrentUserId(c *gin.Context){

}

func (h *Handler)ViewProfile(c *gin.Context){
	
}

func(h *Handler)GetProfile(c *gin.Context){
	
}
func(h *Handler)CheckFollowStatus(c *gin.Context){
	
}
func(h *Handler)FollowUser(c *gin.Context){
	
}
func(h *Handler)UnFollowUser(c *gin.Context){
	
}

func(h *Handler)GetFollowers(c *gin.Context){
	
}
func(h *Handler)GetFollowing(c *gin.Context){
	
}
	

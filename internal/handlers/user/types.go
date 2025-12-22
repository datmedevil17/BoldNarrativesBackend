package user

type SignUpRequest struct{
	Email string `json:"email" binding:"required"`
	Name string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`

}

type SignInRequest struct{
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type FollowRequest struct{
	TargetUserIdParam uint `json:"targetUserIdParam" binding:"required"`
}
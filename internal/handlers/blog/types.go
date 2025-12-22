package blog

type CreateBlogRequest struct{
	Title string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Genre string `json:"genre" binding:"required"`
}

type UpdateBlogRequest struct{
	Title string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Genre string `json:"genre" binding:"required"`
}

type CreateCommentRequest struct{
	Comment string `json:"comment" binding:"required"`
	BlogID uint `json:"blog_id" binding:"required"`
}

type ViewRequest struct{
	ID uint `json:"blog_id" binding:"required"`
}

type VoteRequest struct{
	ID uint `json:"blog_id" binding:"required"`
}

type FilterRequest struct{
	Genre string `json:"genre"`
	AuthorID uint `json:"author_id"`
	Search string `json:"search"`
	Skip int `json:"skip"`
}
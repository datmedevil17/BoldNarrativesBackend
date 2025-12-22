package blog

import (
	"errors"
	"sort"
	"time"

	"github.com/datmedevil17/BoldNarrativesBackend/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

type Filter struct {
	Genre    string
	AuthorID *uint
	Search   string
	Skip     int
	Limit    int
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateBlog(authorId uint, title, content, genre string) (*models.Blog, error) {
	blog := &models.Blog{
		AuthorID: authorId,
		Title:    title,
		Content:  content,
		Genre:    genre,
	}
	if err := s.db.Create(blog).Error; err != nil {
		return nil, err
	}
	s.db.Preload("Author").First(&blog)

	return blog, nil
}

func (s *Service) GetBlogById(blogId uint) (*models.Blog, error) {
	var blog models.Blog
	err := s.db.Preload("Author").Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Preload("Author")
	}).First(&blog, blogId).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (s *Service) UpdateBlog(blogId uint, title, content, genre string, userId uint) (*models.Blog, error) {
	var blog models.Blog
	err := s.db.First(&blog, blogId).Error
	if err != nil {
		return nil, err
	}
	if blog.AuthorID != userId {
		return nil, errors.New("Unauthorized")
	}
	blog.Title = title
	blog.Content = content
	blog.Genre = genre
	if err := s.db.Save(&blog).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}

func (s *Service) DeleteBlog(blogId uint, userId uint) error {
	var blog models.Blog
	err := s.db.First(&blog, blogId).Error
	if err != nil {
		return err
	}
	if blog.AuthorID != userId {
		return errors.New("Unauthorized")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("blog_id=?", blogId).Delete(&models.Comment{}).Error; err != nil {
			return err
		}
		if err := tx.Where("blog_id=?", blogId).Delete(&models.Vote{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&blog).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *Service) GetBlogsCount(opts Filter) (int64, error) {
	query := s.db.Model(&models.Blog{})

	if opts.Genre != "" && opts.Genre != "All" {
		query.Where("genre=?", opts.Genre)
	}
	if opts.AuthorID != nil {
		query.Where("author_id=?", *opts.AuthorID)
	}
	if opts.Search != "" {
		query.Where("title ILIKE ? ", "%"+opts.Search+"%")
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Service) GetBlogsSortedByTime(opts Filter, ascending bool) ([]models.BlogListResponse, error) {
	var blogs []models.Blog
	query := s.db.Preload("Author")
	if opts.Genre != "" && opts.Genre != "All" {
		query.Where("genre=?", opts.Genre)
	}
	if opts.AuthorID != nil {
		query.Where("author_id=?", *opts.AuthorID)
	}
	if opts.Search != "" {
		query.Where("title ILIKE ? ", "%"+opts.Search+"%")
	}
	if ascending {
		query.Order("created_at ASC")
	} else {
		query.Order("created_at DESC")
	}

	if opts.Limit > 0 {
		query = query.Limit(opts.Limit)
	}

	if opts.Skip > 0 {
		query = query.Offset(opts.Skip)
	}
	if err := query.Find(&blogs).Error; err != nil {
		return nil, err
	}
	return s.toBlogListResponse(blogs)

}

func (s *Service) GetBlogsSortedByViews(opts Filter) ([]models.BlogListResponse, error) {
	var blogs []models.Blog
	query := s.db.Preload("Author")
	if opts.Genre != "" && opts.Genre != "All" {
		query.Where("genre=?", opts.Genre)
	}
	if opts.AuthorID != nil {
		query.Where("author_id=?", *opts.AuthorID)
	}
	if opts.Search != "" {
		query.Where("title ILIKE ? ", "%"+opts.Search+"%")
	}
	query = query.Order("views DESC,created_at DESC")

	if opts.Limit > 0 {
		query = query.Limit(opts.Limit)
	}

	if opts.Skip > 0 {
		query = query.Offset(opts.Skip)
	}
	if err := query.Find(&blogs).Error; err != nil {
		return nil, err
	}
	return s.toBlogListResponse(blogs)
}

type TrendingBlog struct {
	models.BlogListResponse
	Score float64 `json:"score"`
}

func (s *Service) GetTrendingBlogs() ([]TrendingBlog, error) {
	var blogs []models.Blog

	// Get top blogs by views
	err := s.db.Preload("Author").
		Order("views DESC").
		Limit(10).
		Find(&blogs).Error

	if err != nil {
		return nil, err
	}

	var trendingBlogs []TrendingBlog

	for _, blog := range blogs {
		// Count votes
		var voteCount int64
		s.db.Model(&models.Vote{}).Where("blog_id = ?", blog.ID).Count(&voteCount)

		// Calculate age in days
		age := time.Since(blog.CreatedAt).Hours() / 24
		if age < 1 {
			age = 1 // Minimum 1 day to avoid division by zero
		}

		// Calculate score: (views + votes * 2) / ageInDays
		score := (float64(blog.Views) + float64(voteCount)*2) / age

		trendingBlogs = append(trendingBlogs, TrendingBlog{
			BlogListResponse: models.BlogListResponse{
				ID:        blog.ID,
				Title:     blog.Title,
				Genre:     blog.Genre,
				Views:     blog.Views,
				VoteCount: voteCount,
				AuthorID:  blog.AuthorID,
				Author:    blog.Author.ToResponse(),
				CreatedAt: blog.CreatedAt,
			},
			Score: score,
		})
	}

	// Sort by score
	sort.Slice(trendingBlogs, func(i, j int) bool {
		return trendingBlogs[i].Score > trendingBlogs[j].Score
	})

	// Return top 10
	if len(trendingBlogs) > 10 {
		trendingBlogs = trendingBlogs[:10]
	}

	return trendingBlogs, nil

}

func (s *Service) IncrementViews(blogId uint) error {
	return s.db.Model(&models.Blog{}).Where("id=?", blogId).Update("views", gorm.Expr("views + 1")).Error
	// blogs++ and blogs.save nhi use kr rhe kyunki Race condition hojati agr do saath me use krte //database khud increment krta hai
}

func (s *Service) ToggleVote(blogId, userId uint) (bool, error) {
	var vote models.Vote
	err := s.db.Where("blog_id=? and user_id=?", blogId, userId).First(&vote).Error
	if err == nil {
		if err := s.db.Delete(&vote).Error; err != nil {
			return false, err
		}
		return false, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newVote := models.Vote{
			BlogID: blogId,
			UserID: userId,
		}
		if err := s.db.Create(&newVote).Error; err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (s *Service) CheckVote(blogId, userId uint) (bool, error) {
	var count int64
	err := s.db.Model(&models.Vote{}).Where("blog_id=? and user_id=?", blogId, userId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Service) CreateComment(blogId, authorId uint, comment string) (*models.Comment, error) {
	newComment := &models.Comment{
		BlogID:   blogId,
		AuthorID: authorId,
		Comment:  comment,
	}
	err := s.db.Create(&newComment).Error
	if err != nil {
		return nil, err
	}
	s.db.Preload("Author").First(&newComment)
	return newComment, nil
}

func (s *Service) GetCommentsByBlogId(blogId uint) ([]models.CommentResponse, error) {
	var comments []models.Comment
	err := s.db.Where("blog_id=?", blogId).Preload("Author").Order("created_at DESC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	var response []models.CommentResponse
	for _, comment := range comments {
		response = append(response, models.CommentResponse{
			ID:        comment.ID,
			Comment:   comment.Comment,
			AuthorID:  comment.AuthorID,
			Author:    comment.Author.ToResponse(),
			CreatedAt: comment.CreatedAt,
		})
	}
	return response, nil
}

func (s *Service) DeleteComment(commentId, userId uint) error {
	var comment models.Comment
	err := s.db.First(&comment, commentId).Error
	if err != nil {
		return err
	}
	if comment.AuthorID != userId {
		return errors.New("Unauthorized")
	}
	return s.db.Delete(&comment).Error
}

func (s *Service) toBlogListResponse(blogs []models.Blog) ([]models.BlogListResponse, error) {
	var response []models.BlogListResponse
	for _, blog := range blogs {
		var voteCount int64
		s.db.Model(&models.Vote{}).Where("blog_id=?", blog.ID).Count(&voteCount)

		response = append(response, models.BlogListResponse{
			ID:        blog.ID,
			Title:     blog.Title,
			Genre:     blog.Genre,
			Views:     blog.Views,
			VoteCount: voteCount,
			AuthorID:  blog.AuthorID,
			Author:    blog.Author.ToResponse(),
			CreatedAt: blog.CreatedAt,
		})
	}
	return response, nil
}

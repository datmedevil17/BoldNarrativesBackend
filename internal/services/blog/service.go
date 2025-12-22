package blog

import (
	"errors"

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

func (s *Service) GetBlogsSortedByTime(filter Filter) ([]models.Blog, error) {}

func (s *Service) GetBlogsSortedByViews(filter Filter) ([]models.Blog, error) {}

type TrendingBlog struct {
	models.BlogListResponse
	Score float64 `json:"score"`
}

func (s *Service) GetTrendingBlogs() ([]TrendingBlog, error) {}

func (s *Service) IncrementViews(blogId uint) error {}

func (s *Service) CheckVote(blogId, userId uint) (bool, error) {}

func (s *Service) CreateComment(blogId, authorId uint, comment string) (*models.Comment, error) {}

func (s *Service) GetCommentsByBlogId(blogId uint) ([]models.CommentResponse, error) {}

func (s *Service) DeleteComment(commentId, userId uint) error {}

func (s *Service) toBlogListResponse(blog *models.Blog) (models.BlogListResponse, error) {}

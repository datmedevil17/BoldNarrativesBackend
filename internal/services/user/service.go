package user

import (
	"errors"

	"github.com/datmedevil17/BoldNarrativesBackend/internal/models"
	"github.com/datmedevil17/BoldNarrativesBackend/internal/utils"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateUser(name, email, password string) (*models.User, error) {
	var existingUser models.User
	if err := s.db.Where("email=?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPasswod(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) AuthenticateUser(email, password string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email=?", email).First(&user).Error; err == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Record Not Found")
		}
	}
	if !utils.CheckPasswordHash(password,user.Password){
		return nil,errors.New("Invalid Password")
	}
	return &user,nil
}

func (s *Service)GetUserById(id string)(*models.User,error){}

func (s *Service)GetUserProfile(id string)(*models.User,error){}

func (s *Service)FollowUser(followerId,followingId uint)error{}

func (s *Service)UnFollowUser(followerId,followingId uint)error{}

func (s *Service)GetFollowers(userId uint)([]models.User,error){}

func (s *Service)GetFollowing(userId uint)([]models.FollowResponse,error){
	var follows []models.Follows
	err :=s.db.Where("follower_id=?",userId).Preload("Following").Find(&follows).Error
	if err != nil {
		return nil,err
	}
	var following []models.FollowResponse
	for _,follow := range follows {
		following = append(following,models.FollowResponse{
			ID:follow.Following.ID,
			Name:follow.Following.Name,
		})
	}
	return following,nil
	
}

	
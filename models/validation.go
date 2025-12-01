package models

import (
	"errors"
	"net/mail"
	"strings"
)

// ValidateUserSignup validates user signup request
func (u *UserSignupRequest) Validate() error {
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("email is required")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("invalid email format")
	}

	if strings.TrimSpace(u.Name) == "" {
		return errors.New("name is required")
	}

	if len(u.Name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}

	if strings.TrimSpace(u.Password) == "" {
		return errors.New("password is required")
	}

	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}

// ValidateUserCredentials validates user login credentials
func (u *UserCredentials) Validate() error {
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("email is required")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("invalid email format")
	}

	if strings.TrimSpace(u.Password) == "" {
		return errors.New("password is required")
	}

	return nil
}

// ValidateBlogRequest validates blog creation request
func (b *BlogRequest) Validate() error {
	if strings.TrimSpace(b.Title) == "" {
		return errors.New("title is required")
	}

	if len(b.Title) < 3 {
		return errors.New("title must be at least 3 characters long")
	}

	if len(b.Title) > 200 {
		return errors.New("title cannot exceed 200 characters")
	}

	if strings.TrimSpace(b.Content) == "" {
		return errors.New("content is required")
	}

	if len(b.Content) < 10 {
		return errors.New("content must be at least 10 characters long")
	}

	if strings.TrimSpace(b.Genre) == "" {
		return errors.New("genre is required")
	}

	return nil
}

// ValidateCommentRequest validates comment creation request
func (c *CommentRequest) Validate() error {
	if strings.TrimSpace(c.Comment) == "" {
		return errors.New("comment is required")
	}

	if len(c.Comment) < 1 {
		return errors.New("comment cannot be empty")
	}

	if len(c.Comment) > 1000 {
		return errors.New("comment cannot exceed 1000 characters")
	}

	return nil
}

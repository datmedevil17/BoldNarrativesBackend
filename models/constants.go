package models

// Database collection names
const (
	UsersCollection    = "users"
	BlogsCollection    = "blogs"
	CommentsCollection = "comments"
	VotesCollection    = "votes"
	FollowsCollection  = "follows"
)

// API response messages
const (
	MsgUserCreated       = "User created successfully"
	MsgUserSignedIn      = "User signed in successfully"
	MsgUserNotFound      = "User not found"
	MsgInvalidCredentials = "Invalid email or password"
	MsgBlogCreated       = "Blog created successfully"
	MsgBlogUpdated       = "Blog updated successfully"
	MsgBlogDeleted       = "Blog deleted successfully"
	MsgBlogNotFound      = "Blog not found"
	MsgCommentCreated    = "Comment created successfully"
	MsgCommentDeleted    = "Comment deleted successfully"
	MsgCommentNotFound   = "Comment not found"
	MsgFollowSuccess     = "Successfully followed user"
	MsgUnfollowSuccess   = "Successfully unfollowed user"
	MsgVoteAdded         = "Vote added successfully"
	MsgVoteRemoved       = "Vote removed successfully"
	MsgViewIncremented   = "View count incremented"
	MsgUnauthorized      = "Unauthorized access"
	MsgForbidden         = "Forbidden access"
	MsgInternalError     = "Internal server error"
	MsgInvalidRequest    = "Invalid request format"
)

// Blog sorting options
const (
	SortByTime     = "time"
	SortByViews    = "views"
	SortByTrending = "trending"
)

// Default pagination values
const (
	DefaultLimit = 10
	DefaultSkip  = 0
	MaxLimit     = 100
)

// Blog genres (you can expand this list)
var ValidGenres = []string{
	"Technology",
	"Health",
	"Finance",
	"Travel",
	"Food",
	"Lifestyle",
	"Entertainment",
	"Sports",
	"Education",
	"Business",
	"Science",
	"Politics",
	"Art",
	"Music",
	"Other",
}

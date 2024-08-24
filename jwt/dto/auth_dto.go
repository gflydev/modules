package dto

// SignUp struct to describe register a new user.
type SignUp struct {
	Email    string `json:"email" example:"john@jivecode.com" validate:"required,email,lte=255"`
	Password string `json:"password" example:"M1PassW@s" validate:"required,gte=6"`
	Fullname string `json:"fullname" example:"John Doe" validate:"required,lte=255"`
	Phone    string `json:"phone" example:"0989831911" validate:"required,lte=20"`
	Status   string `json:"status" example:"pending"  validate:"omitempty"`
	Avatar   string `json:"avatar" example:"https://i.pravatar.cc/32"  validate:"omitempty"`
}

// SignIn struct to describe sign in user
type SignIn struct {
	Username string `json:"username" validate:"required,email,lte=255"`
	Password string `json:"password" validate:"required,gte=6"`
}

// RefreshToken struct to refresh JWT token.
type RefreshToken struct {
	Token string `json:"token" validate:"required,lte=255"`
}

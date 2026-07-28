package users

type CreateUserRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	ConfigLimit int    `json:"config_limit"`
	GroupID     *uint  `json:"group_id"`
	IsAdmin     bool   `json:"is_admin"`
}

type UpdateUserRequest struct {
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	ConfigLimit *int   `json:"config_limit,omitempty"`
	GroupID     *uint  `json:"group_id,omitempty"`
	IsAdmin     *bool  `json:"is_admin,omitempty"`
}

type UserResponse struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	ConfigLimit int    `json:"config_limit"`
	GroupID     *uint  `json:"group_id"`
	IsAdmin     bool   `json:"is_admin"`
	CreatedAt   string `json:"created_at"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func ToResponse(user User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		ConfigLimit: user.ConfigLimit,
		GroupID:     user.GroupID,
		IsAdmin:     user.IsAdmin,
		CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToResponseList(users []User) []UserResponse {

	result := make([]UserResponse, 0, len(users))

	for _, user := range users {
		result = append(result, ToResponse(user))
	}

	return result
}

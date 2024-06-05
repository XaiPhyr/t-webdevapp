package models

type (
	Authentication struct {
		Username     string `json:"username"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		RememberMe   bool   `json:"remember_me"`
	}

	Login struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		UserType   string `json:"user_type"`
		RememberMe bool   `json:"remember_me"`
	}
)

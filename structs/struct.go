package structs

// Config struct
type Config struct {
	Auth        AuthConfig
	GMail       GMailConfig
	Environment EnvironmentConfig
}

// AuthConfig struct
type AuthConfig struct {
	ActivationSalt            string
	ExpirationOfActivationKey int64
}

// GMailConfig struct
type GMailConfig struct {
	SenderAddress  string
	SenderPassword string
}

// EnvironmentConfig struct
type EnvironmentConfig struct {
	Host string
}

type (
	// User struct
	User struct {
		ID            int64  `json:"id" form:"id" query:"id"`
		Email         string `json:"email" form:"email" query:"email" valid:"required,email"`
		Password      string `json:"-" form:"password" valid:"required,alphanum"`
		Activated     int64  `json:"activated"`
		ActivationKey string `json:"-" query:"activation_key"`
	}
)

package core

type AppConfig struct {
	AppName    string `json:"app_name"`
	APIPrefix  string `json:"api_prefix"`
	APIVersion string `json:"api_version"`
	AppSecret  string `json:"app_secret"`
	DB         struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	TokenConfig struct {
		Duration int64  `json:"duration"` // Token duration in seconds
		Secret   string `json:"secret"`   // Secret key for signing tokens
	}
	AuthConfig struct {
		HeaderPrefix string `json:"header_prefix"` // Prefix for the authorization header
	}
}

func GetAppConfig() *AppConfig {
	return &AppConfig{
		AppName:    "MyApp",
		APIPrefix:  "/api",
		APIVersion: "v1",
		AppSecret:  "supersecretkey",
		DB: struct {
			Host     string `json:"host"`
			Port     int    `json:"port"`
			User     string `json:"user"`
			Password string `json:"password"`
			Name     string `json:"name"`
		}{
			Host:     "localhost",
			Port:     5432,
			User:     "dbuser",
			Password: "dbpassword",
			Name:     "myappdb",
		},
		TokenConfig: struct {
			Duration int64  `json:"duration"` // Token duration in seconds
			Secret   string `json:"secret"`   // Secret key for signing tokens,
		}{
			Duration: 604800, // 7 days in seconds
			Secret:   "mysecretkey",
		},
		AuthConfig: struct {
			HeaderPrefix string `json:"header_prefix"` // Prefix for the authorization header
		}{
			HeaderPrefix: "Bearer ",
		},
	}
}

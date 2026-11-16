package ms_auth

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	Err          error
}

type XboxAuthResult struct {
	Token    string
	UserHash string
}

type XSTSAuthResult struct {
	Token    string
	UserHash string
}

type MinecraftAuthResult struct {
	AccessToken string
	ExpiresIn   int64
}

type MinecraftProfile struct {
	ID   string
	Name string
}

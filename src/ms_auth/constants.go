package ms_auth

const (
	clientID    = "e4789891-9cae-424f-baba-4f8aa9cf90c1"
	redirectURI = "https://syaniclauncher.github.io/auth_callback"
	successURI  = "https://syaniclauncher.github.io/auth_success"

	microsoftAuthorizeURL  = "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize"
	microsoftTokenURL      = "https://login.microsoftonline.com/consumers/oauth2/v2.0/token"
	microsoftDeviceCodeURL = "https://login.microsoftonline.com/consumers/oauth2/v2.0/devicecode"
	xboxAuthURL            = "https://user.auth.xboxlive.com/user/authenticate"
	xstsAuthorizeURL       = "https://xsts.auth.xboxlive.com/xsts/authorize"
	minecraftLoginURL      = "https://api.minecraftservices.com/authentication/login_with_xbox"
	minecraftProfileURL    = "https://api.minecraftservices.com/minecraft/profile"
	microsoftAuthScope     = "XboxLive.signin offline_access"
	xboxSiteName           = "user.auth.xboxlive.com"
	xboxRelyingParty       = "http://auth.xboxlive.com"
	minecraftRelyingParty  = "rp://api.minecraftservices.com/"
)

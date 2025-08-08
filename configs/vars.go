package configs

import "github.com/labstack/echo/v4"

const (
	ProductName    = "NexVul"
	UserCookieName = "3LwqrpZpHXK9z2pvX4"
	DomainName     = ".nexvul.com"
	LogoPath       = "./assets/logo.png"
	HTMLPageURI    = "https://nexvul.com"
	ShotlinkAPI    = "https://shotlink.nexvul.com/get?url="
	FontURL        = "https://github.com/melroy89/Roboto/raw/refs/heads/main/RobotoTTF/Roboto-Regular.ttf"
	FontPath       = "./temp_fonts/Roboto-Regular.ttf"
)

var RequiredHeaders = []string{
	"Content-Security-Policy",
	"X-Frame-Options",
	"Strict-Transport-Security",
	"X-Content-Type-Options",
	"Referrer-Policy",
	"Permissions-Policy",
	"X-XSS-Protection",
	"Expect-CT",
	"Cross-Origin-Resource-Policy",
	"Cross-Origin-Opener-Policy",
	"Cross-Origin-Embedder-Policy",
	"Access-Control-Allow-Origin",
	"Access-Control-Allow-Credentials",
	"Access-Control-Allow-Methods",
	"Access-Control-Allow-Headers",
}

func GetRootURL(c echo.Context) string {
	scheme := c.Scheme()
	host := c.Request().Host
	return scheme + "://" + host
}

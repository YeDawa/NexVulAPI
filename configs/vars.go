package configs

import "github.com/labstack/echo/v4"

const (
	UserCookieName = "3LwqrpZpHXK9z2pvX4"
	DomainName     = ".httpshield.net"
	HTMLPageURI    = "https://httpshield.net/"
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
	"Feature-Policy",
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

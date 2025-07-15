package utils

func HeaderDescription(header string) string {
	switch header {
	case "Content-Security-Policy":
		return "Helps prevent Cross-Site Scripting (XSS) and other code injection attacks by defining allowed content sources."
	case "X-Frame-Options":
		return "Protects against clickjacking by preventing the site from being embedded in an iframe."
	case "Strict-Transport-Security":
		return "Forces browsers to interact with the site only over HTTPS, reducing the risk of man-in-the-middle attacks."
	case "X-Content-Type-Options":
		return "Prevents MIME-sniffing by telling the browser to respect the declared Content-Type."
	case "Referrer-Policy":
		return "Controls how much referrer information is included with requests to other sites, improving privacy."
	case "Permissions-Policy":
		return "Gives fine-grained control over which browser features (camera, mic, etc.) are allowed."
	case "X-XSS-Protection":
		return "Enables basic XSS filtering in some older browsers (deprecated in most modern ones)."
	case "Expect-CT":
		return "Helps enforce Certificate Transparency to detect misissued SSL certificates."
	case "Feature-Policy":
		return "Predecessor to Permissions-Policy; used to limit features like geolocation or fullscreen."
	case "Cross-Origin-Resource-Policy":
		return "Controls which origins can load resources from your site, protecting against cross-origin leaks."
	case "Cross-Origin-Opener-Policy":
		return "Isolates browsing contexts to mitigate cross-origin attacks like Spectre."
	case "Cross-Origin-Embedder-Policy":
		return "Prevents untrusted cross-origin resources from being embedded unless explicitly allowed."
	case "Access-Control-Allow-Origin":
		return "Specifies which origins can access resources via CORS requests."
	case "Access-Control-Allow-Credentials":
		return "Indicates whether credentials (cookies, HTTP auth) are allowed in CORS requests."
	case "Access-Control-Allow-Methods":
		return "Lists the allowed HTTP methods (GET, POST, etc.) for CORS requests."
	case "Access-Control-Allow-Headers":
		return "Specifies allowed request headers for CORS requests."
	case "Access-Control-Expose-Headers":
		return "Specifies which response headers are safe to expose to the browser in CORS."
	case "Access-Control-Max-Age":
		return "Defines how long a preflight request can be cached to reduce CORS overhead."
	default:
		return "This header is recommended to improve the security and behavior of your site."
	}
}

func GenerateRecommendation(header string) string {
	switch header {
	case "X-Content-Type-Options":
		return "Add 'X-Content-Type-Options' with value 'nosniff' to prevent MIME-type sniffing."
	case "Referrer-Policy":
		return "Add 'Referrer-Policy' to control referrer information."
	case "Strict-Transport-Security":
		return "Add 'Strict-Transport-Security' to enforce HTTPS connections."
	case "Content-Security-Policy":
		return "Add 'Content-Security-Policy' to mitigate XSS and code injection attacks."
	case "X-Frame-Options":
		return "Add 'X-Frame-Options' to prevent clickjacking."
	case "Permissions-Policy":
		return "Add 'Permissions-Policy' to control browser feature access."
	case "Cross-Origin-Embedder-Policy":
		return "Add 'Cross-Origin-Embedder-Policy' to isolate cross-origin resources."
	case "Cross-Origin-Opener-Policy":
		return "Add 'Cross-Origin-Opener-Policy' to enhance site isolation."
	case "Cross-Origin-Resource-Policy":
		return "Add 'Cross-Origin-Resource-Policy' to restrict resource sharing."
	case "Access-Control-Allow-Origin":
		return "Add 'Access-Control-Allow-Origin' to define CORS origins."
	case "Access-Control-Allow-Methods":
		return "Add 'Access-Control-Allow-Methods' to define allowed HTTP methods."
	case "Access-Control-Allow-Headers":
		return "Add 'Access-Control-Allow-Headers' to define allowed headers in CORS."
	case "Access-Control-Allow-Credentials":
		return "Add 'Access-Control-Allow-Credentials' to allow cookies and auth in CORS."
	case "Access-Control-Expose-Headers":
		return "Add 'Access-Control-Expose-Headers' to expose additional headers to browsers."
	case "Access-Control-Max-Age":
		return "Add 'Access-Control-Max-Age' to cache CORS preflight response."
	case "Expect-CT":
		return "Add 'Expect-CT' to enable Certificate Transparency enforcement."
	case "Feature-Policy":
		return "Add 'Feature-Policy' (now replaced by Permissions-Policy) to restrict APIs."
	case "X-XSS-Protection":
		return "Add 'X-XSS-Protection' to enable basic XSS filters (note: deprecated in modern browsers)."
	default:
		return "Consider implementing " + header + " for improved security posture."
	}
}
package tasks

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RobotsExposure struct {
	Target  string              `json:"target"`
	Domain  string              `json:"domain"`
	Exposed []SensitiveExposure `json:"exposed_paths"`
}

type SensitiveExposure struct {
	Path       string `json:"path"`
	FullURL    string `json:"full_url"`
	Reason     string `json:"reason"`
	UserAgent  string `json:"user_agent"`
	Accessible bool   `json:"accessible"`
}

var sensitivePaths = []string{
	"/admin/", "/administrator/", "/backup/", "/backups/", "/config/", "/configuration/", "/database/", "/databases/", "/private/", "/priv/",
	"/dev/", "/development/", "/internal/", "/test/", "/testing/", "/old/", "/hidden/", "/users/", "/user/", "/logs/",
	"/log/", "/data/", "/dump/", "/tmp/", "/temp/", "/secret/", "/secrets/", "/.git/", "/.env/", "/.htaccess/",
	"/.htpasswd/", "/api/", "/api/v1/", "/api/v2/", "/staging/", "/sandbox/", "/debug/", "/core/", "/bin/", "/cgi-bin/",
	"/conf/", "/etc/", "/uploads/", "/upload/", "/downloads/", "/download/", "/scripts/", "/source/", "/src/", "/passwords/",
	// WordPress sensitive paths
	"/wp-admin/", "/wp-login.php", "/wp-config.php", "/wp-content/", "/wp-includes/", "/wp-json/", "/wp-cron.php", "/wp-signup.php", "/wp-links-opml.php", "/wp-comments-post.php",
}

func AnalyzeRobotsSensitivePaths(rawURL string) (RobotsExposure, error) {
	target := normalizeURL(rawURL)
	robotsURL := strings.TrimSuffix(target, "/") + "/robots.txt"

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(robotsURL)
	if err != nil {
		return RobotsExposure{}, fmt.Errorf("erro ao acessar robots.txt: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RobotsExposure{}, fmt.Errorf("robots.txt não encontrado (%d)", resp.StatusCode)
	}

	// Lê todo conteúdo do robots.txt para buscar os sensitivePaths
	var robotsLines []string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		robotsLines = append(robotsLines, line)
	}

	domain := extractDomain(target)
	report := RobotsExposure{Target: target, Domain: domain}

	for _, sensitive := range sensitivePaths {
		for _, line := range robotsLines {
			if strings.Contains(line, sensitive) {
				fullURL := joinURL(target, sensitive)
				accessible := isURLAccessible(fullURL)

				report.Exposed = append(report.Exposed, SensitiveExposure{
					Path:       sensitive,
					FullURL:    fullURL,
					Reason:     "path listed in robots.txt",
					UserAgent:  extractUserAgent(line),
					Accessible: accessible,
				})

				break
			}
		}
	}

	return report, nil
}

func extractUserAgent(line string) string {
	if strings.HasPrefix(strings.ToLower(line), "user-agent:") {
		return strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
	}
	
	return "*"
}

func normalizeURL(u string) string {
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		return "https://" + u
	}
	return u
}

func containsPrefix(list []string, path string) bool {
	for _, item := range list {
		if strings.HasPrefix(path, item) || strings.HasPrefix(item, path) {
			return true
		}
	}
	return false
}

func joinURL(base, path string) string {
	u, err := url.Parse(base)
	if err != nil {
		return base + path
	}
	u.Path = strings.TrimSuffix(u.Path, "/") + path
	return u.String()
}

func isURLAccessible(fullURL string) bool {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(fullURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func extractDomain(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

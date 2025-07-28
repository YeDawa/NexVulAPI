package tasks

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type RobotsDirective struct {
	UserAgent string   `json:"user_agent"`
	Allow     []string `json:"allow"`
	Disallow  []string `json:"disallow"`
}

type RobotsData struct {
	Target    string            `json:"target"`
	Sitemaps  []string          `json:"sitemaps"`
	Directives []RobotsDirective `json:"directives"`
}

func ParseRobotsTxt(target string) (RobotsData, error) {
	robotsURL := strings.TrimSuffix(normalizeURL(target), "/") + "/robots.txt"
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(robotsURL)
	if err != nil {
		return RobotsData{}, fmt.Errorf("error accessing robots.txt: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RobotsData{}, fmt.Errorf("robots.txt not found (%d)", resp.StatusCode)
	}

	var (
		currentAgent string
		directives   = make(map[string]*RobotsDirective)
		sitemaps     []string
	)

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(strings.ToLower(line), "user-agent:") {
			currentAgent = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			if _, exists := directives[currentAgent]; !exists {
				directives[currentAgent] = &RobotsDirective{
					UserAgent: currentAgent,
					Allow:     []string{},
					Disallow:  []string{},
				}
			}
		} else if strings.HasPrefix(strings.ToLower(line), "disallow:") {
			if currentAgent != "" {
				path := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
				directives[currentAgent].Disallow = append(directives[currentAgent].Disallow, path)
			}
		} else if strings.HasPrefix(strings.ToLower(line), "allow:") {
			if currentAgent != "" {
				path := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
				directives[currentAgent].Allow = append(directives[currentAgent].Allow, path)
			}
		} else if strings.HasPrefix(strings.ToLower(line), "sitemap:") {
			sitemap := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			sitemaps = append(sitemaps, sitemap)
		}
	}

	var result []RobotsDirective
	for _, v := range directives {
		result = append(result, *v)
	}

	return RobotsData{
		Target:    target,
		Sitemaps:  sitemaps,
		Directives: result,
	}, nil
}

func normalizeURL(u string) string {
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		return "https://" + u
	}
	return u
}
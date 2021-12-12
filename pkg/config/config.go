package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kdevo/config"
	"github.com/kdevo/config/provider"
)

type Config struct {
	RepoOwner string
	RepoName  string

	CategoryName     string
	DiscussionOpener string

	OutputFile string

	SiteRSSURL    string
	SiteMapURL    string
	SiteURLPrefix string

	EventName string
	EventPath string
}

func (c *Config) Validate() error {
	var errors config.Errors
	if c.RepoName == "" {
		errors.Add(config.EmptyErr("RepoName", ""))
	}
	if c.RepoOwner == "" {
		errors.Add(config.EmptyErr("RepoOwner", ""))
	}
	if c.CategoryName == "" {
		errors.Add(config.EmptyErr("CategoryName", ""))
	}
	if c.OutputFile == "" {
		errors.Add(config.EmptyErr("OutputFile", ""))
	}
	if c.SiteMapURL == "" && c.SiteRSSURL == "" {
		errors.Add(config.EmptyErr("SiteMapURL", c.SiteMapURL))
		errors.Add(config.EmptyErr("SiteRSSURL", c.SiteRSSURL))
	}
	if c.SiteMapURL != "" && !strings.HasPrefix(c.SiteMapURL, "http") {
		errors.Add(config.Err("SiteMapURL", c.SiteMapURL, "must be a valid URL (starting with http)"))
	}
	if c.SiteRSSURL != "" && !strings.HasPrefix(c.SiteRSSURL, "http") {
		errors.Add(config.Err("SiteRSSURL", c.SiteRSSURL, "must be a valid URL (starting with http)"))
	}
	return errors.AsError()
}

func (c *Config) Config() (interface{}, error) {
	return c, c.Validate()
}

func (c *Config) Name() string {
	return "Static"
}

func (c *Config) String() string {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("config (unmarshal error)")
	}
	return string(data)
}

func Load() (*Config, error) {
	loader := config.From(provider.Dynamic(func() (interface{}, error) {
		return nil, nil
	})).WithDefaults(provider.Dynamic(
		func() (interface{}, error) {
			var errors config.Errors
			repo := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")
			var repoOwner, repoName string
			if len(repo) == 2 {
				repoOwner = repo[0]
				repoName = repo[1]
			} else {
				errors.Add(config.Err("RepoOwner", repo, fmt.Sprintf("env GITHUB_REPOSITORY uses incorrect format, want {owner}/{repo}")))
				errors.Add(config.Err("RepoName", repo, fmt.Sprintf("env GITHUB_REPOSITORY uses incorrect format, want {owner}/{repo}")))
			}
			return &Config{
				RepoOwner:        repoOwner,
				RepoName:         repoName,
				CategoryName:     os.Getenv("CATEGORY_NAME"),
				DiscussionOpener: os.Getenv("DISCUSSION_OPENER"),
				OutputFile:       os.Getenv("OUTPUT_FILE"),

				SiteRSSURL:    os.Getenv("SITE_RSS_URL"),
				SiteMapURL:    os.Getenv("SITE_MAP_URL"),
				SiteURLPrefix: os.Getenv("SITE_URL_PREFIX"),

				EventName: os.Getenv("GITHUB_EVENT_NAME"),
				EventPath: os.Getenv("GITHUB_EVENT_PATH"),
			}, errors.AsError()
		},
	).WithName("Environment")).
		WithDefaults(&Config{
			CategoryName:     "Blog",
			OutputFile:       "data/discussions.json",
			DiscussionOpener: "Blog post: {{ .URL }}",
			SiteURLPrefix:    "http",
		})
	var cfg Config
	err := loader.Resolve(&cfg)
	return &cfg, err
}

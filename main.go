package main

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2"

	"github.com/hugo-mods/discussions/bridge/pkg/config"
	"github.com/hugo-mods/discussions/bridge/pkg/github"
	"github.com/hugo-mods/discussions/bridge/pkg/model"
	"github.com/hugo-mods/discussions/bridge/pkg/site"
)

func main() {
	cfg, err := config.Load()
	fmt.Println("got config:", cfg)
	if err != nil {
		fatal("configuration error: %s", err)
	}

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("REPO_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), tokenSource)

	client := github.New(httpClient, cfg.RepoOwner, cfg.RepoName)

	categories, err := client.Categories()
	if err != nil {
		fatal("could not retrieve categories: %v", err)
	}
	if len(categories) == 0 {
		fatal("could not find any categories. please ensure that discussions are enabled and there is at least one category.")
	}
	category := categories.ByName(cfg.CategoryName, 1).First()
	if category == nil {
		fatal("could not find discussion with name %q", cfg.CategoryName)
	}
	fmt.Println("got category ID:", category.ID)

	discussions, err := client.Discussions(category.ID)
	if err != nil {
		fatal("could not get discussions for category %q", category.ID)
	}

	webSite, err := site.New(cfg.SiteMapURL, cfg.SiteRSSURL, cfg.DiscussionOpener)
	if err != nil {
		fatal("could not create site: %v", err)
	}
	siteDiscussions := webSite.RelateDiscussions(model.FromGitHubDiscussions(discussions))

	if eventName := cfg.EventName; eventName != "" {
		fmt.Println("triggered by:", eventName)
		fmt.Println("  event path:", cfg.EventPath)
	}
	switch cfg.EventName {
	case "push":
		pages, err := webSite.Pages(cfg.SiteURLPrefix)
		if err != nil {
			fatal("could not get site's pages: %v", err)
		}

		var newPages []site.Page
		for url := range pages {
			if !siteDiscussions.HasPage(url) {
				newPages = append(newPages, pages[url])
			}
		}
		fmt.Printf("got %d pages from site. found %d unsynced discussions.\n", len(pages), len(newPages))
		for _, p := range newPages {
			disc, err := webSite.NewDiscussion(p)
			if err != nil {
				fmt.Printf("could not create discussion: %v", err)
				continue
			}
			if _, err := client.CreateDiscussion(category.ID, disc.Title, disc.Body); err != nil {
				fmt.Printf("could not create discussion: %v", err)
			}
		}
	case "discussion", "discussion_comment":
		if err := siteDiscussions.Save(cfg.OutputFile); err != nil {
			fatal("could not save discussions: %v", err)
		}
		fmt.Printf("wrote %d discussions to %s\n", len(siteDiscussions), cfg.OutputFile)
	default:
		fmt.Printf("unhandled event name %q. doing nothing.\n", cfg.EventName)
	}
}

func fatal(msg string, arg ...interface{}) {
	panic(fmt.Sprintf(msg, arg...))
}

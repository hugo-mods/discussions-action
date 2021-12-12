package site

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/hugo-mods/discussions/bridge/pkg/model"
)

type Site struct {
	SitemapURL     string
	RSSURL         string
	openerTemplate *template.Template
	openerURLRegEx *regexp.Regexp
}

func New(sitemapURL string, rssURL string, opener string) (*Site, error) {
	urlRE := regexp.MustCompile(`{{\s*.URL\s*}}`)
	loc := urlRE.FindStringIndex(opener)
	if loc == nil {
		return nil, fmt.Errorf("could not find {{ .URL }} anywhere in the opener")
	}
	openerRE, err := regexp.Compile(
		regexp.QuoteMeta(opener[:loc[0]]) +
			`(.*)` +
			regexp.QuoteMeta(opener[loc[1]:]))
	if err != nil {
		return nil, err
	}

	template := template.New("opener")
	template, err = template.Parse(opener)
	if err != nil {
		return nil, err
	}
	return &Site{
		SitemapURL:     sitemapURL,
		RSSURL:         rssURL,
		openerTemplate: template,
		openerURLRegEx: openerRE,
	}, nil
}

func (s *Site) RelateDiscussions(ds []model.Discussion) Discussions {
	sds := make(Discussions, len(ds))
	for _, d := range ds {
		subs := s.openerURLRegEx.FindStringSubmatch(d.Message.Body)
		if len(subs) > 1 {
			url := subs[1]
			sds[url] = d
		}
	}
	return sds
}

type Page struct {
	URL         string
	Title       string
	Description string
	UpdatedAt   *time.Time
}

func (s *Site) NewDiscussion(p Page) (*model.Discussion, error) {
	title := p.Title
	if p.Title == "" {
		title = fmt.Sprintf("%s - %s", p.URL, p.UpdatedAt)
	}

	body := &bytes.Buffer{}
	err := s.openerTemplate.Execute(body, p)
	if err != nil {
		return nil, fmt.Errorf("could not exec template: %w", err)
	}
	return &model.Discussion{
		Title: title,
		Message: model.Message{
			Body:     body.String(),
			BodyMIME: "text/markdown",
		},
	}, nil
}

// Pages tries to collect the site's pages by first trying RSS and then Sitemap.
func (s *Site) Pages(urlPrefix string) (map[string]Page, error) {
	var err error
	var pages map[string]Page
	if s.RSSURL != "" {
		pages, err = s.RSS(urlPrefix)
		if err == nil {
			return pages, nil
		}
	}
	if s.SitemapURL != "" {
		pages, err = s.Sitemap(urlPrefix)
		if err == nil {
			return pages, nil
		}
	}
	return nil, fmt.Errorf("could not get pages: %v", err)
}

func (s *Site) Sitemap(urlPrefix string) (map[string]Page, error) {
	resp, err := http.Get(s.SitemapURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response: %w", err)
	}

	type Sitemap struct {
		XMLName xml.Name `xml:"urlset"`
		Pages   []struct {
			URL     string    `xml:"loc"`
			LastMod time.Time `xml:"lastmod"`
		} `xml:"url"`
	}
	sitemap := Sitemap{}
	err = xml.Unmarshal(body, &sitemap)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal sitemap: %v", err)
	}

	result := make(map[string]Page, len(sitemap.Pages)/2)
	for _, p := range sitemap.Pages {
		if strings.HasPrefix(p.URL, urlPrefix) {
			result[p.URL] = Page{URL: p.URL, UpdatedAt: &p.LastMod}
		}
	}
	return result, nil
}

func (s *Site) RSS(urlPrefix string) (map[string]Page, error) {
	resp, err := http.Get(s.RSSURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response: %w", err)
	}

	type Sitemap struct {
		XMLName xml.Name `xml:"rss"`
		Pages   []struct {
			URL         string `xml:"link"`
			LastMod     string `xml:"pubDate"`
			Title       string `xml:"title"`
			Description string `xml:"description"`
		} `xml:"channel>item"`
	}
	sitemap := Sitemap{}
	err = xml.Unmarshal(body, &sitemap)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal RSS feed: %v", err)
	}

	result := make(map[string]Page, len(sitemap.Pages)/2)
	for _, p := range sitemap.Pages {
		if strings.HasPrefix(p.URL, urlPrefix) {
			var updatedAt *time.Time
			for _, layout := range []string{time.RFC822Z, time.RFC822, time.RFC1123Z, time.RFC1123} {
				if upd, err := time.Parse(layout, p.LastMod); err == nil {
					updatedAt = &upd
					break
				}
			}
			result[p.URL] = Page{
				URL:         p.URL,
				Title:       p.Title,
				Description: p.Description,
				UpdatedAt:   updatedAt,
			}
		}
	}
	return result, nil
}

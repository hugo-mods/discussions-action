package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shurcooL/githubv4"
)

type Client struct {
	gql   *githubv4.Client
	owner string
	repo  string

	maxDiscussions int
	maxComments    int
}

func New(client *http.Client, owner string, repo string) *Client {
	gql := githubv4.NewClient(client)
	return &Client{
		gql:   gql,
		owner: owner,
		repo:  repo,

		maxComments:    50,
		maxDiscussions: 100,
	}
}

func (c *Client) WithMaxComments(n int) *Client {
	c.maxComments = n
	return c
}

func (c *Client) WithMaxDiscussions(n int) *Client {
	c.maxDiscussions = n
	return c
}

func (c *Client) Discussions(categoryID string) ([]Discussion, error) {
	var q struct {
		Repository struct {
			Discussions struct {
				Nodes      []Discussion
				TotalCount int
			} `graphql:"discussions(first: $firstDiscussions, categoryId: $categoryID, orderBy: {field: UPDATED_AT, direction: DESC})"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	err := c.gql.Query(context.Background(), &q,
		map[string]interface{}{
			"owner":            githubv4.String(c.owner),
			"name":             githubv4.String(c.repo),
			"categoryID":       githubv4.ID(categoryID),
			"firstDiscussions": githubv4.Int(c.maxDiscussions),
			"firstComments":    githubv4.Int(c.maxComments),
			"firstReactions":   githubv4.Int(50),
		},
	)
	if err != nil {
		return nil, err
	}
	return q.Repository.Discussions.Nodes, nil
}

func (c *Client) Categories() (Categories, error) {
	var q struct {
		Repository struct {
			DiscussionCategories struct {
				Nodes      []Category
				TotalCount int
			} `graphql:"discussionCategories(first: $n)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	err := c.gql.Query(context.Background(), &q,
		map[string]interface{}{
			"owner": githubv4.String(c.owner),
			"name":  githubv4.String(c.repo),
			"n":     githubv4.Int(100),
		},
	)
	if err != nil {
		return nil, err
	}
	return Categories(q.Repository.DiscussionCategories.Nodes), nil
}

func (c *Client) CreateDiscussion(categoryID, title, body string) (string, error) {
	var q struct {
		Repository struct {
			ID string
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	err := c.gql.Query(context.Background(), &q,
		map[string]interface{}{
			"owner": githubv4.String(c.owner),
			"name":  githubv4.String(c.repo),
		},
	)
	if err != nil {
		return "", fmt.Errorf("could not to get repository ID: %v", err)
	}
	repositoryID := q.Repository.ID

	var m struct {
		CreateDiscussion struct {
			Discussion struct {
				ID string
			}
		} `graphql:"createDiscussion(input: $input)"`
	}
	input := githubv4.CreateDiscussionInput{
		RepositoryID: githubv4.ID(repositoryID),
		CategoryID:   githubv4.ID(categoryID),
		Title:        githubv4.String(title),
		Body:         githubv4.String(body),
	}
	err = c.gql.Mutate(context.Background(), &m, input, nil)
	if err != nil {
		return "", fmt.Errorf("could not create discussion: %v", err)
	}
	return m.CreateDiscussion.Discussion.ID, nil
}

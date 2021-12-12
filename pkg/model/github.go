package model

import (
	"strings"

	"github.com/hugo-mods/discussions/bridge/pkg/github"
)

// FromGitHubDiscussions converts the GitHub discussion to an independent Discussion model.
// For compactness, it converts only one level of comments, e.g.:
//  Discussion
//  ├── Comment #1
//  ├── Comment #2
// Note that GitHub also supports replies for comments:
//  Discussion
//  ├── Comment #1
//  ├── Comment #2
//  │   ├── Reply #1
// The exact content will not be included, only the number of replies.
func FromGitHubDiscussions(ghds []github.Discussion) []Discussion {
	ds := make([]Discussion, len(ghds))
	for i := range ghds {
		ds[i] = *FromGitHubDiscussion(&ghds[i])
	}
	return ds

}

func FromGitHubDiscussion(ghd *github.Discussion) *Discussion {
	return &Discussion{
		Title: ghd.Title,
		Message: Message{
			URL:          ghd.URL,
			Author:       FromGitHubAuthor(ghd.Author),
			Body:         ghd.Body,
			BodyMIME:     "text/markdown",
			UpvotesCount: 0,
			Reactions:    FromGitHubReactions(ghd.Reactions.Nodes),
		},
		Comments: FromGitHubComments(ghd.Comments.Nodes),
	}

}

func FromGitHubAuthor(gha github.Author) Author {
	return Author{
		User:       User{Name: gha.Login},
		FullName:   "",
		PictureURL: gha.AvatarURL,
	}
}

func FromGitHubComments(ghcs []github.Comment) []Comment {
	comments := make([]Comment, len(ghcs))
	for i := range ghcs {
		comments[i] = FromGitHubComment(ghcs[i])
	}
	return comments
}

func FromGitHubComment(ghc github.Comment) Comment {
	return Comment{
		Message: Message{
			URL:          ghc.URL,
			Author:       FromGitHubAuthor(ghc.Author),
			Body:         ghc.Body,
			BodyMIME:     "text/markdown",
			UpvotesCount: ghc.UpvoteCount,
			Reactions:    FromGitHubReactions(ghc.Reactions.Nodes),
		},
		CommentsCount: ghc.Replies.TotalCount,
	}
}

func FromGitHubReactions(ghrs []github.Reaction) Reactions {
	reactions := make(Reactions, len(ghrs))
	for _, ghr := range ghrs {
		r := FromGitHubReaction(ghr)
		reactions[r] = append(reactions[r], User{Name: ghr.User.Login})
	}
	return reactions
}

func FromGitHubReaction(ghr github.Reaction) EmojiCode {
	var code EmojiCode
	switch strings.ToLower(ghr.Content) {
	case "thumbs_up":
		code = ThumbsUp
	case "thumbs_down":
		code = ThumbsDown
	case "smile":
		code = Smile
	case "hooray":
		code = Party
	case "confused":
		code = Confused
	case "heart":
		code = Heart
	case "rocket":
		code = Rocket
	case "eyes":
		code = Eyes
	default:
		code = ThoughtBalloon
	}
	return code
}

package model

// Message is anything that is written by an author and can have replies.
type Message struct {
	// URL where the comment is located at.
	URL string `json:"url"`
	// Author is the user who created the comment.
	Author Author `json:"author"`
	// Body is the discussion's main text.
	Body string `json:"body"`
	// BodyMIME is the MIME type of the Body, e.g. "text/markdown" or "text/plain".
	BodyMIME string `json:"bodyMimeType"`
	// UpvotesCount describes how many times the discussion has been found useful.
	UpvotesCount int `json:"upvotesCount"`
	// Reactions are used to express a feeling by an emoji.
	Reactions Reactions `json:"reactions,omitempty"`
}

// Discussions can have comments that are arbitrary nested.
type Discussion struct {
	Message
	// Title is the subject which briefly describes what the discussion is about.
	Title string `json:"title"`
	// Comments are the comments of the discussion.
	Comments []Comment `json:"comments"`
}

type Comment struct {
	Message
	// Comments can be commented, too. Optional, can be "collapsed". In this case, only CommentsCount is given.
	Comments []Comment `json:"comments,omitempty"`
	// CommentsCount describes how many comment replies there are and is always given.
	CommentsCount int `json:"commentsCount,omitempty"`
}

type Author struct {
	User
	// FullName is the author's real name, typically consisting of a prename and a surname.
	FullName string
	// PictureURL is the author's profile picture URL.
	PictureURL string
}

type User struct {
	// Name is the user's unique name.
	Name string
}

type EmojiCode string

const (
	ThumbsUp   EmojiCode = ":+1:"
	ThumbsDown EmojiCode = ":-1:"
	Smile      EmojiCode = ":smile:"
	Party      EmojiCode = ":tada:"
	Confused   EmojiCode = ":confused:"
	Heart      EmojiCode = ":heart:"
	Rocket     EmojiCode = ":rocket:"
	Eyes       EmojiCode = ":eyes:"

	Thinking       EmojiCode = ":thinking:"
	ThoughtBalloon EmojiCode = ":thought_balloon:"
)

// Reactions map an emoji to users (the users who have selected this emoji).
type Reactions map[EmojiCode][]User

package github

type Discussion struct {
	URL         string
	Title       string
	Body        string
	Author      Author
	Locked      bool
	UpvoteCount int
	Comments    struct {
		Nodes      []Comment
		TotalCount int
	} `graphql:"comments(first: $firstComments)"`
	Reactions struct {
		Nodes      []Reaction
		TotalCount int
	} `graphql:"reactions(first: $firstReactions)"`
}

type Author struct {
	Login     string
	URL       string
	AvatarURL string `graphql:"avatarUrl(size: 64)"`
}

type Comment struct {
	URL               string
	Author            Author
	AuthorAssociation string
	Body              string
	UpvoteCount       int
	Reactions         struct {
		Nodes      []Reaction
		TotalCount int
	} `graphql:"reactions(first: $firstReactions)"`
	Replies struct {
		TotalCount int
	}
}

type Reaction struct {
	User    Author
	Content string
}

type Discussions []Discussion

func (ds Discussions) ByFilter(filter func(c Discussion) bool, n int) Discussions {
	res := make(Discussions, 0, n)
	for i := 0; i < len(ds) && len(res) < n; i++ {
		if !filter(ds[i]) {
			continue
		}
		res = append(res, ds[i])
	}
	return res
}

type Category struct {
	ID    string
	Emoji string
	Name  string
	// Description string
}

type Categories []Category

// ByName filters categories by name and stops after the nth occurrence.
func (cs Categories) ByName(name string, n int) Categories {
	return cs.ByFilter(func(c Category) bool { return c.Name == name }, n)
}

func (cs Categories) ByFilter(filter func(c Category) bool, n int) Categories {
	res := make(Categories, 0, n)
	for i := 0; i < len(cs) && len(res) < n; i++ {
		if !filter(cs[i]) {
			continue
		}
		res = append(res, cs[i])
	}
	return res
}

func (cs Categories) First() *Category {
	if len(cs) == 0 {
		return nil
	} else {
		return &cs[0]
	}
}

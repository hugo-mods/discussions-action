package site_test

import (
	"reflect"
	"testing"

	"github.com/hugo-mods/discussions/bridge/pkg/model"
	"github.com/hugo-mods/discussions/bridge/pkg/site"
)

func TestInferDiscussion(t *testing.T) {
	testCases := []struct {
		opener      string
		discussions []model.Discussion
		want        site.Discussions
	}{
		{
			opener: "**Blog Post**: {{ .URL }}",
			discussions: []model.Discussion{
				{Message: model.Message{Body: "**Blog Post**: https://hugo-mods.github.io/blog/icons/"}},
				{Message: model.Message{Body: "Hi there, here's my post: https://hugo-mods.github.io/blog/post/"}},
				{Message: model.Message{Body: "**Blog Post**: https://hugo-mods.github.io/blog/test/", Author: model.Author{FullName: "Test"}}},
			},
			want: site.Discussions{
				"https://hugo-mods.github.io/blog/icons/": {Message: model.Message{Body: "**Blog Post**: https://hugo-mods.github.io/blog/icons/"}},
				"https://hugo-mods.github.io/blog/test/":  {Message: model.Message{Body: "**Blog Post**: https://hugo-mods.github.io/blog/test/", Author: model.Author{FullName: "Test"}}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.opener, func(t *testing.T) {
			s, err := site.New("", "", tc.opener)
			if err != nil {
				t.Fatal(err)
			}
			got := s.RelateDiscussions(tc.discussions)
			if ok := reflect.DeepEqual(tc.want, got); !ok {
				t.Errorf("unexpected result:\n  want=%v\n   got=%v", tc.want, got)
			}
		})
	}
}

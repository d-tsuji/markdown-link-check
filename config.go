package mlc

import (
	"context"
	"regexp"
	"strings"

	"github.com/google/go-github/v32/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

type config struct {
	ctx      context.Context
	client   *github.Client
	owner    string
	repo     string
	branch   string
	token    string
	allMode  bool
	filePath string
}

var pat = regexp.MustCompile(`http.://raw\.githubusercontent\.com/(.+?)/(.+?)/(.+?)/(.*)`)

func NewConfig(c *cli.Context) *config {
	if c.Bool("all") {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: c.String("token")},
		)
		tc := oauth2.NewClient(c.Context, ts)
		client := github.NewClient(tc)
		return &config{
			ctx:     c.Context,
			client:  client,
			owner:   c.String("user"),
			repo:    c.String("repo"),
			branch:  c.String("branch"),
			token:   c.String("token"),
			allMode: c.Bool("all"),
		}
	} else {
		g := pat.FindSubmatch([]byte(c.Args().First()))
		return &config{
			ctx:      c.Context,
			client:   github.NewClient(nil),
			owner:    string(g[1]),
			repo:     string(g[2]),
			branch:   string(g[3]),
			filePath: string(g[4]),
		}
	}
}

func (c config) fetchFiles(ctx context.Context, path string) ([]*content, error) {
	var contents []*content
	_, dirContents, _, err := c.client.Repositories.GetContents(ctx, c.owner, c.repo, path, &github.RepositoryContentGetOptions{Ref: c.branch})
	if err != nil {
		return nil, err
	}
	for _, v := range dirContents {
		if *v.Type == "dir" {
			fs, err := c.fetchFiles(ctx, *v.Path)
			if err != nil {
				return nil, err
			}
			contents = append(contents, fs...)
			continue
		}
		if strings.HasSuffix(*v.Name, ".md") {
			contents = append(contents, &content{filePath: *v.Path})
		}
	}
	return contents, nil
}

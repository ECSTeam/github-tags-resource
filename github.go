package resource

import (
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

type GitHub interface {
	ListTags() ([]*github.RepositoryTag, error)
	GetTag(id string) (*github.RepositoryTag, error)
	GetTarballLink(tag string) (*url.URL, error)
	GetZipballLink(tag string) (*url.URL, error)
}

type GitHubClient struct {
	client *github.Client

	user       string
	repository string
}

func NewGitHubClient(source Source) (*GitHubClient, error) {
	var client *github.Client

	if source.AccessToken == "" {
		client = github.NewClient(nil)
	} else {
		var err error
		client, err = oauthClient(source)
		if err != nil {
			return nil, err
		}
	}

	if source.GitHubAPIURL != "" {
		var err error
		client.BaseURL, err = url.Parse(source.GitHubAPIURL)
		if err != nil {
			return nil, err
		}

		client.UploadURL, err = url.Parse(source.GitHubAPIURL)
		if err != nil {
			return nil, err
		}
	}

	if source.GitHubUploadsURL != "" {
		var err error
		client.UploadURL, err = url.Parse(source.GitHubUploadsURL)
		if err != nil {
			return nil, err
		}
	}

	return &GitHubClient{
		client:     client,
		user:       source.User,
		repository: source.Repository,
	}, nil
}

func (g *GitHubClient) ListTags() ([]*github.RepositoryTag, error) {
	tags, res, err := g.client.Repositories.ListTags(g.user, g.repository, nil)
	if err != nil {
		return []*github.RepositoryTag{}, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (g *GitHubClient) GetTag(id string) (*github.RepositoryTag, error) {
	tags, err := g.ListTags()
	if err != nil {
		return &github.RepositoryTag{}, err
	}

	for _, tag := range tags {
		if *tag.Name == id {
			return tag, nil
		}
	}

	return nil, nil
}

func (g *GitHubClient) GetTarballLink(tag string) (*url.URL, error) {
	opt := &github.RepositoryContentGetOptions{Ref: tag}
	u, res, err := g.client.Repositories.GetArchiveLink(g.user, g.repository, github.Tarball, opt)
	if err != nil {
		return nil, err
	}
	res.Body.Close()
	return u, nil
}

func (g *GitHubClient) GetZipballLink(tag string) (*url.URL, error) {
	opt := &github.RepositoryContentGetOptions{Ref: tag}
	u, res, err := g.client.Repositories.GetArchiveLink(g.user, g.repository, github.Zipball, opt)
	if err != nil {
		return nil, err
	}
	res.Body.Close()
	return u, nil
}

func oauthClient(source Source) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: source.AccessToken,
	})

	oauthClient := oauth2.NewClient(oauth2.NoContext, ts)

	githubHTTPClient := &http.Client{
		Transport: oauthClient.Transport,
	}

	return github.NewClient(githubHTTPClient), nil
}

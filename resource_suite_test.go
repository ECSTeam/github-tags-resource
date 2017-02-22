package resource_test

import (
	"testing"

	"github.com/google/go-github/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGithubReleaseResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Github Release Resource Suite")
}

func newRepositoryTag(id int, version string) *github.RepositoryTag {
	return &github.RepositoryTag{
		Name: github.String(version),
	}
}

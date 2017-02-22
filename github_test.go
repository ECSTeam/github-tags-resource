package resource_test

import (
	"net/http"

	. "github.com/ecsteam/github-tags-resource"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("GitHub Client", func() {
	var server *ghttp.Server
	var client *GitHubClient
	var source Source

	BeforeEach(func() {
		server = ghttp.NewServer()
	})

	JustBeforeEach(func() {
		source.GitHubAPIURL = server.URL()

		var err error
		client, err = NewGitHubClient(source)
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		server.Close()
	})

	Context("with bad URLs", func() {
		BeforeEach(func() {
			source.AccessToken = "hello?"
		})

		It("returns an error if the API URL is bad", func() {
			source.GitHubAPIURL = ":"

			_, err := NewGitHubClient(source)
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if the API URL is bad", func() {
			source.GitHubUploadsURL = ":"

			_, err := NewGitHubClient(source)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("with an OAuth Token", func() {
		BeforeEach(func() {
			source = Source{
				User:        "concourse",
				Repository:  "concourse",
				AccessToken: "abc123",
			}

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/repos/concourse/concourse/tags"),
					ghttp.RespondWith(200, "[]"),
					ghttp.VerifyHeaderKV("Authorization", "Bearer abc123"),
				),
			)
		})

		It("sends one", func() {
			_, err := client.ListTags()
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("without an OAuth Token", func() {
		BeforeEach(func() {
			source = Source{
				User:       "concourse",
				Repository: "concourse",
			}

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/repos/concourse/concourse/tags"),
					ghttp.RespondWith(200, "[]"),
					ghttp.VerifyHeader(http.Header{"Authorization": nil}),
				),
			)
		})

		It("sends one", func() {
			_, err := client.ListTags()
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("GetRelease", func() {
		BeforeEach(func() {
			source = Source{
				User:       "concourse",
				Repository: "concourse",
			}
		})
		Context("When GitHub's rate limit has been exceeded", func() {
			BeforeEach(func() {
				rateLimitResponse := `{
          "message": "API rate limit exceeded for 127.0.0.1. (But here's the good news: Authenticated requests get a higher rate limit. Check out the documentation for more details.)",
          "documentation_url": "https://developer.github.com/v3/#rate-limiting"
        }`

				rateLimitHeaders := http.Header(map[string][]string{
					"X-RateLimit-Limit":     {"60"},
					"X-RateLimit-Remaining": {"0"},
					"X-RateLimit-Reset":     {"1377013266"},
				})

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/repos/concourse/concourse/tags"),
						ghttp.RespondWith(403, rateLimitResponse, rateLimitHeaders),
					),
				)
			})

			It("Returns an appropriate error", func() {
				_, err := client.GetTag("20")
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("API rate limit exceeded for 127.0.0.1. (But here's the good news: Authenticated requests get a higher rate limit. Check out the documentation for more details.)"))
			})
		})
	})

	Describe("GetReleaseByTag", func() {
		BeforeEach(func() {
			source = Source{
				User:       "concourse",
				Repository: "concourse",
			}
		})
		Context("When GitHub's rate limit has been exceeded", func() {
			BeforeEach(func() {
				rateLimitResponse := `{
          "message": "API rate limit exceeded for 127.0.0.1. (But here's the good news: Authenticated requests get a higher rate limit. Check out the documentation for more details.)",
          "documentation_url": "https://developer.github.com/v3/#rate-limiting"
        }`

				rateLimitHeaders := http.Header(map[string][]string{
					"X-RateLimit-Limit":     {"60"},
					"X-RateLimit-Remaining": {"0"},
					"X-RateLimit-Reset":     {"1377013266"},
				})

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/repos/concourse/concourse/releases/tags/some-tag"),
						ghttp.RespondWith(403, rateLimitResponse, rateLimitHeaders),
					),
				)
			})
		})
	})
})

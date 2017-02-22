package resource_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/google/go-github/github"

	"github.com/ecsteam/github-tags-resource"
	"github.com/ecsteam/github-tags-resource/fakes"
)

var _ = Describe("Check Command", func() {
	var (
		githubClient *fakes.FakeGitHub
		command      *resource.CheckCommand

		returnedTags []*github.RepositoryTag
	)

	BeforeEach(func() {
		githubClient = &fakes.FakeGitHub{}
		command = resource.NewCheckCommand(githubClient)

		returnedTags = []*github.RepositoryTag{}
	})

	JustBeforeEach(func() {
		githubClient.ListTagsReturns(returnedTags, nil)
	})

	Context("when this is the first time that the resource has been run", func() {
		Context("when there are no releases", func() {
			BeforeEach(func() {
				returnedTags = []*github.RepositoryTag{}
			})

			It("returns no versions", func() {
				versions, err := command.Run(resource.CheckRequest{})
				Ω(err).ShouldNot(HaveOccurred())
				Ω(versions).Should(BeEmpty())
			})
		})

		Context("when there are releases", func() {
			BeforeEach(func() {
				returnedTags = []*github.RepositoryTag{
					newRepositoryTag(1, "v0.4.0"),
					newRepositoryTag(2, "0.1.3"),
					newRepositoryTag(3, "v0.1.2"),
				}
			})

			It("outputs the most recent version only", func() {
				command := resource.NewCheckCommand(githubClient)

				response, err := command.Run(resource.CheckRequest{})
				Ω(err).ShouldNot(HaveOccurred())

				Ω(response).Should(HaveLen(1))
				Ω(response[0]).Should(Equal(resource.Version{
					Tag: "v0.4.0",
				}))
			})
		})
	})

	Context("when there are prior versions", func() {
		Context("when there are no releases", func() {
			BeforeEach(func() {
				returnedTags = []*github.RepositoryTag{}
			})

			It("returns no versions", func() {
				versions, err := command.Run(resource.CheckRequest{})
				Ω(err).ShouldNot(HaveOccurred())
				Ω(versions).Should(BeEmpty())
			})
		})

		Context("when there are releases", func() {
			Context("and the releases do not contain a draft release", func() {
				BeforeEach(func() {
					returnedTags = []*github.RepositoryTag{
						newRepositoryTag(1, "v0.1.4"),
						newRepositoryTag(2, "0.4.0"),
						newRepositoryTag(3, "v0.1.3"),
						newRepositoryTag(4, "0.1.2"),
					}
				})

				It("returns an empty list if the lastet version has been checked", func() {
					command := resource.NewCheckCommand(githubClient)

					response, err := command.Run(resource.CheckRequest{
						Version: resource.Version{
							Tag: "0.4.0",
						},
					})
					Ω(err).ShouldNot(HaveOccurred())

					Ω(response).Should(BeEmpty())
				})

				It("returns all of the versions that are newer", func() {
					command := resource.NewCheckCommand(githubClient)

					response, err := command.Run(resource.CheckRequest{
						Version: resource.Version{
							Tag: "v0.1.3",
						},
					})
					Ω(err).ShouldNot(HaveOccurred())

					Ω(response).Should(Equal([]resource.Version{
						{Tag: "v0.1.3"},
						{Tag: "v0.1.4"},
						{Tag: "0.4.0"},
					}))
				})

				It("returns the latest version if the current version is not found", func() {
					command := resource.NewCheckCommand(githubClient)

					response, err := command.Run(resource.CheckRequest{
						Version: resource.Version{
							Tag: "v3.4.5",
						},
					})
					Ω(err).ShouldNot(HaveOccurred())

					Ω(response).Should(Equal([]resource.Version{
						{Tag: "0.4.0"},
					}))
				})

				Context("when there are not-quite-semver versions", func() {
					BeforeEach(func() {
						returnedTags = append(returnedTags, newRepositoryTag(5, "v1"))
						returnedTags = append(returnedTags, newRepositoryTag(6, "v0"))
					})

					It("combines them with the semver versions in a reasonable order", func() {
						command := resource.NewCheckCommand(githubClient)

						response, err := command.Run(resource.CheckRequest{
							Version: resource.Version{
								Tag: "v0.1.3",
							},
						})
						Ω(err).ShouldNot(HaveOccurred())

						Ω(response).Should(Equal([]resource.Version{
							{Tag: "v0.1.3"},
							{Tag: "v0.1.4"},
							{Tag: "0.4.0"},
							{Tag: "v1"},
						}))
					})
				})
			})
		})
	})
})

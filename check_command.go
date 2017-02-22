package resource

import (
	"sort"

	"github.com/google/go-github/github"

	"github.com/cppforlife/go-semi-semantic/version"
)

type CheckCommand struct {
	github GitHub
}

func NewCheckCommand(github GitHub) *CheckCommand {
	return &CheckCommand{
		github: github,
	}
}

func (c *CheckCommand) Run(request CheckRequest) ([]Version, error) {
	tags, err := c.github.ListTags()
	if err != nil {
		return []Version{}, err
	}

	if len(tags) == 0 {
		return []Version{}, nil
	}

	var filteredTags []*github.RepositoryTag

	for _, release := range tags {

		// Should we skip this release
		//   a- prerelease condition dont match our source config
		//   b- release condition match  prerealse in github since github has true/false to describe release/prerelase
		// if request.Source.PreTag != *release.Prerelease && request.Source.Tag == *release.Prerelease {
		// 	continue
		// }
		//
		// if release.TagName == nil {
		// 	continue
		// }
		// if _, err := version.NewVersionFromString(determineVersionFromTag(*release.TagName)); err != nil {
		// 	continue
		// }

		filteredTags = append(filteredTags, release)
	}

	sort.Sort(byVersion(filteredTags))

	if len(filteredTags) == 0 {
		return []Version{}, nil
	}
	latestTag := filteredTags[len(filteredTags)-1]

	if (request.Version == Version{}) {
		return []Version{
			versionFromTag(latestTag),
		}, nil
	}

	if *latestTag.Name == request.Version.Tag {
		return []Version{}, nil
	}

	upToLatest := false
	reversedVersions := []Version{}

	for _, release := range filteredTags {
		if !upToLatest {
			version := *release.Name
			upToLatest = request.Version.Tag == version
		}

		if upToLatest {
			reversedVersions = append(reversedVersions, versionFromTag(release))
		}
	}

	if !upToLatest {
		// current version was removed; start over from latest
		reversedVersions = append(
			reversedVersions,
			versionFromTag(filteredTags[len(filteredTags)-1]),
		)
	}

	return reversedVersions, nil
}

type byVersion []*github.RepositoryTag

func (r byVersion) Len() int {
	return len(r)
}

func (r byVersion) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r byVersion) Less(i, j int) bool {
	first, err := version.NewVersionFromString(determineVersionFromTag(*r[i].Name))
	if err != nil {
		return true
	}

	second, err := version.NewVersionFromString(determineVersionFromTag(*r[j].Name))
	if err != nil {
		return false
	}

	return first.IsLt(second)
}

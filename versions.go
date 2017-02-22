package resource

import (
	"regexp"

	"github.com/google/go-github/github"
)

// determineVersionFromTag converts git tags v1.2.3 into semver 1.2.3 values
func determineVersionFromTag(tag string) string {
	re := regexp.MustCompile("v?([^v].*)")
	matches := re.FindStringSubmatch(tag)
	if len(matches) > 0 {
		return matches[1]
	} else {
		return ""
	}
}

func versionFromTag(release *github.RepositoryTag) Version {
	return Version{Tag: *release.Name}
}

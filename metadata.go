package resource

import "github.com/google/go-github/github"

func metadataFromTag(release *github.RepositoryTag) []MetadataPair {
	metadata := []MetadataPair{}

	if release.Name != nil {
		nameMeta := MetadataPair{
			Name:  "name",
			Value: *release.Name,
		}

		metadata = append(metadata, nameMeta)
	}

	if release.Name != nil {
		metadata = append(metadata, MetadataPair{
			Name:  "tag",
			Value: *release.Name,
		})
	}

	return metadata
}

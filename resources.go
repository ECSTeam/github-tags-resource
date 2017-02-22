package resource

type Source struct {
	User       string `json:"user"`
	Repository string `json:"repository"`

	GitHubAPIURL     string `json:"github_api_url"`
	GitHubUploadsURL string `json:"github_uploads_url"`
	AccessToken      string `json:"access_token"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

func NewCheckRequest() CheckRequest {
	res := CheckRequest{}
	return res
}

func NewInRequest() InRequest {
	res := InRequest{}
	return res
}

type InRequest struct {
	Source  Source   `json:"source"`
	Version *Version `json:"version"`
	Params  InParams `json:"params"`
}

type InParams struct {
	IncludeSourceTarball bool `json:"include_source_tarball"`
	IncludeSourceZip     bool `json:"include_source_zip"`
}

type InResponse struct {
	Version  Version        `json:"version"`
	Metadata []MetadataPair `json:"metadata"`
}

type Version struct {
	Tag string `json:"tag,omitempty"`
	ID  string `json:"id,omitempty"`
}

type MetadataPair struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	URL      string `json:"url"`
	Markdown bool   `json:"markdown"`
}

package resource

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
)

type InCommand struct {
	github GitHub
	writer io.Writer
}

func NewInCommand(github GitHub, writer io.Writer) *InCommand {
	return &InCommand{
		github: github,
		writer: writer,
	}
}

func (c *InCommand) Run(destDir string, request InRequest) (InResponse, error) {
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return InResponse{}, err
	}

	var foundTag *github.RepositoryTag
	foundTag, err = c.github.GetTag(request.Version.ID)

	if err != nil {
		return InResponse{}, err
	}

	if foundTag == nil {
		return InResponse{}, errors.New("no tags")
	}

	if foundTag.Name != nil && *foundTag.Name != "" {
		tagPath := filepath.Join(destDir, "tag")
		err = ioutil.WriteFile(tagPath, []byte(*foundTag.Name), 0644)
		if err != nil {
			return InResponse{}, err
		}

		version := determineVersionFromTag(*foundTag.Name)
		versionPath := filepath.Join(destDir, "version")
		err = ioutil.WriteFile(versionPath, []byte(version), 0644)
		if err != nil {
			return InResponse{}, err
		}

	}

	if request.Params.IncludeSourceTarball {
		u, err := c.github.GetTarballLink(request.Version.Tag)
		if err != nil {
			return InResponse{}, err
		}
		fmt.Fprintln(c.writer, "downloading source tarball to source.tar.gz")
		if err := c.downloadFile(u.String(), filepath.Join(destDir, "source.tar.gz")); err != nil {
			return InResponse{}, err
		}
	}

	if request.Params.IncludeSourceZip {
		u, err := c.github.GetZipballLink(request.Version.Tag)
		if err != nil {
			return InResponse{}, err
		}
		fmt.Fprintln(c.writer, "downloading source zip to source.zip")
		if err := c.downloadFile(u.String(), filepath.Join(destDir, "source.zip")); err != nil {
			return InResponse{}, err
		}
	}

	return InResponse{
		Version:  versionFromTag(foundTag),
		Metadata: metadataFromTag(foundTag),
	}, nil
}

func (c *InCommand) downloadFile(url, destPath string) error {
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file `%s`: HTTP status %d", filepath.Base(destPath), resp.StatusCode)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

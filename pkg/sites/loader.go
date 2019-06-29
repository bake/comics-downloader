package sites

import (
	"errors"
	"fmt"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	log "github.com/sirupsen/logrus"
)

func initializeCollection(issues []string, url, format, source string, siteSource BaseSite) ([]*core.Comic, error) {
	var collection []*core.Comic
	var err error

	if len(issues) == 0 {
		return collection, errors.New("No issues found.")
	}

	for _, url := range issues {
		name, issueNumber := GetInfo(siteSource, url)
		name = util.Parse(name)
		issueNumber = util.Parse(issueNumber)

		dir, _ := util.PathSetup(source, name)
		fileName := util.GenerateFileName(dir, name, issueNumber, format)

		if util.FileDoesNotExist(fileName) {
			comic := &core.Comic{
				Name:        name,
				IssueNumber: issueNumber,
				URLSource:   url,
				Source:      source,
				Format:      format,
			}
			if err = Initialize(siteSource, comic); err != nil {
				return collection, err
			}
			collection = append(collection, comic)
		} else {
			log.Info(fmt.Sprintf("%s/%s.%s Already exist", name, issueNumber, format))
		}
	}

	return collection, nil
}

// LoadComicFromSource will return a `comic` instance initialized based on the source
func LoadComicFromSource(source, url, country, format string, all, last bool) ([]*core.Comic, error) {
	var siteSource BaseSite
	var collection []*core.Comic
	var issues []string
	var err error

	options := map[string]string{"country": country}

	switch source {
	case "www.comicextra.com":
		siteSource = &Comicextra{}
	case "mangarock.com":
		siteSource = NewMangarock(options)
	case "www.mangareader.net":
		siteSource = &Mangareader{}
	case "www.mangatown.com":
		siteSource = &Mangatown{}
	default:
		err = fmt.Errorf("It was not possible to determine the source")
		return collection, err
	}

	issues, err = RetrieveIssueLinks(siteSource, url, all, last)
	if err != nil {
		return collection, err
	}

	return initializeCollection(issues, url, format, source, siteSource)
}

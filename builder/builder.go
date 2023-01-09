// Package builder implements the logic to get things from Discogs
// and covert them over to the format we'll use for drtp.
package builder

import (
	"encoding/json"
	"fmt"

	drtppb "github.com/brotherlogic/damprecordsthepast/proto"
)

type retriever interface {
	get(url string) ([]byte, error)
}

var (
	urlBase = "https://api.discogs.com/"
)

type Bridge struct {
	r retriever
}

func (b *Bridge) GetReleases(artist int32) ([]*drtppb.Release, error) {
	releases, _, err := b.pullReleases(artist, 1)

	var pr []*drtppb.Release
	for _, release := range releases {
		pr = append(pr,
			&drtppb.Release{
				Id: int32(release.Id),
			})
	}

	return pr, err
}

type ReleaseReturn struct {
	Pagination *Pagination
	Releases   []*ReleasePageData
}

type Pagination struct {
	PerPage int `json:"per_page"`
	Items   int
	Pages   int
}

type ReleasePageData struct {
	Id int
}

func (b *Bridge) pullReleases(artist int32, pageNumber int) ([]*ReleasePageData, *Pagination, error) {
	url := fmt.Sprintf("%vartists/%v/releases", urlBase, artist)

	res, err := b.r.get(url)
	if err != nil {
		return nil, nil, err
	}

	ret := &ReleaseReturn{}
	err = json.Unmarshal(res, ret)
	if err != nil {
		return nil, nil, err
	}

	return ret.Releases, ret.Pagination, nil
}

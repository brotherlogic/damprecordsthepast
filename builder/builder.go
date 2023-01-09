// Package builder implements the logic to get things from Discogs
// and covert them over to the format we'll use for drtp.
package builder

import (
	"encoding/json"
	"fmt"
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

func (b *Bridge) GetReleases(artist int32) {
	b.pullReleases(artist, 1)
}

type ReleaseReturn struct {
	Pagination *Pagination
	Releases   []*ReleasePageData
}

type Pagination struct {
	perPage int
	items   int
	pages   int
}

type ReleasePageData struct {
	id int
}

func (b *Bridge) pullReleases(artist int32, pageNumber int) ([]*ReleasePageData, *Pagination, error) {
	url := fmt.Sprintf("%vartists/%v/releases", urlBase, artist)

	res, err := b.r.get(url)
	if err != nil {
		return nil, nil, err
	}

	var ret *ReleaseReturn
	err = json.Unmarshal(res, ret)
	if err != nil {
		return nil, nil, err
	}

	return ret.Releases, ret.Pagination, nil
}

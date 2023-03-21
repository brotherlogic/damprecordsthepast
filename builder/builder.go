// Package builder implements the logic to get things from Discogs
// and covert them over to the format we'll use for drtp.
package builder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	drtppb "github.com/brotherlogic/damprecordsthepast/proto"
)

type retriever interface {
	get(url string) ([]byte, error)
}

var (
	urlBase = "https://api.discogs.com/"
)

func GetBridge() *Bridge {
	return &Bridge{r: &ProdRetriever{}}
}

type Bridge struct {
	r retriever
}

type ProdRetriever struct{}

func (p *ProdRetriever) get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	log.Printf("GET %v -> %v", url, resp.StatusCode)

	if resp.StatusCode != 200 {
		if resp.StatusCode == 429 {
			time.Sleep(time.Second * 10)
			return p.get(url)
		}
		return []byte{}, fmt.Errorf("non-200 response: %v -> %v", resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (b *Bridge) getSubReleases(artist int32) ([]*drtppb.Release, error) {
	releases, pagination, err := b.pullSubReleases(artist, 1)
	if err != nil {
		return nil, err
	}

	var pr []*drtppb.Release
	for _, release := range releases {
		rr, err := b.pullRelease(int32(release.Id))
		if err != nil {
			return nil, err
		}
		pr = append(pr, rr)
	}
	for pageNumber := 2; pageNumber <= pagination.Pages; pageNumber++ {
		releases, _, err := b.pullSubReleases(artist, pageNumber)
		if err != nil {
			return nil, err
		}

		for _, release := range releases {
			rr, err := b.pullRelease(int32(release.Id))
			if err != nil {
				return nil, err
			}
			pr = append(pr, rr)
		}
	}

	return pr, err
}

func (b *Bridge) GetReleases(artist int32) ([]*drtppb.Release, error) {
	releases, pagination, err := b.pullReleases(artist, 1)
	if err != nil {
		return nil, err
	}

	var pr []*drtppb.Release
	for _, release := range releases {
		if release.MainRelease > 0 {
			subreleases, err := b.getSubReleases(int32(release.Id))
			if err != nil {
				return nil, err
			}
			pr = append(pr, subreleases...)
		} else {
			rr, err := b.pullRelease(int32(release.Id))
			if err != nil {
				return nil, err
			}
			pr = append(pr, rr)
		}
	}

	for pageNumber := 2; pageNumber <= pagination.Pages; pageNumber++ {
		releases, _, err := b.pullReleases(artist, pageNumber)
		if err != nil {
			return nil, err
		}

		for _, release := range releases {
			if release.MainRelease > 0 {
				subreleases, err := b.getSubReleases(int32(release.Id))
				if err != nil {
					return nil, err
				}
				pr = append(pr, subreleases...)
			} else {
				rr, err := b.pullRelease(int32(release.Id))
				if err != nil {
					return nil, err
				}
				pr = append(pr, rr)
			}
		}
	}

	return pr, err
}

type ReleaseReturn struct {
	Pagination *Pagination
	Releases   []*ReleasePageData
}

type SubReleaseReturn struct {
	Pagination *Pagination
	Versions   []*ReleasePageData
}

type Pagination struct {
	PerPage int `json:"per_page"`
	Items   int
	Pages   int
}

type ReleasePageData struct {
	Id          int `json:"id"`
	MainRelease int `json:"main_release"`
	Title       string
	Tracklist   []Track
}

type Track struct {
	Position string
	Title    string
	Type     string
	Artists  []Artist
}

type Artist struct {
	Id int `json:"id"`
}

func (b *Bridge) pullReleases(artist int32, pageNumber int) ([]*ReleasePageData, *Pagination, error) {
	url := fmt.Sprintf("%vartists/%v/releases?page=%v", urlBase, artist, pageNumber)

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

func (b *Bridge) pullRelease(release int32) (*drtppb.Release, error) {
	url := fmt.Sprintf("%vreleases/%v", urlBase, release)

	res, err := b.r.get(url)
	if err != nil {
		return nil, err
	}

	ret := &ReleasePageData{}
	err = json.Unmarshal(res, ret)
	if err != nil {
		return nil, err
	}

	rel := &drtppb.Release{
		Id:    int32(ret.Id),
		Title: ret.Title,
	}

	for _, track := range ret.Tracklist {
		artists := []*drtppb.Artist{}
		for _, artist := range track.Artists {
			artists = append(artists, &drtppb.Artist{Id: int32(artist.Id)})
		}
		rel.Tracks = append(rel.Tracks, &drtppb.Track{Title: track.Title, Artists: artists})
	}

	return rel, nil
}

func (b *Bridge) pullSubReleases(release int32, pageNumber int) ([]*ReleasePageData, *Pagination, error) {
	url := fmt.Sprintf("%vmasters/%v/versions?page=%v", urlBase, release, pageNumber)

	res, err := b.r.get(url)
	if err != nil {
		return nil, nil, err
	}

	ret := &SubReleaseReturn{}
	err = json.Unmarshal(res, ret)
	if err != nil {
		return nil, nil, err
	}

	return ret.Versions, ret.Pagination, nil
}

type UserCollectionPageData struct {
	Pagination *Pagination
	Releases   []*ReleasePageData
}

func (b *Bridge) pullUserCollection(username string, pageNumber int) ([]*ReleasePageData, *Pagination, error) {
	url := fmt.Sprintf("%vusers/%v/collection/folders/0/releases?page=%v&per_page=100", urlBase, username, pageNumber)

	res, err := b.r.get(url)
	if err != nil {
		return nil, nil, err
	}

	ret := &UserCollectionPageData{}
	err = json.Unmarshal(res, ret)
	if err != nil {
		return nil, nil, err
	}

	return ret.Releases, ret.Pagination, nil
}

func (b *Bridge) GetUserCollection(user string) ([]*drtppb.Release, error) {
	releases, pagination, err := b.pullUserCollection(user, 1)
	if err != nil {
		return nil, err
	}

	var pr []*drtppb.Release
	for _, release := range releases {
		pr = append(pr,
			&drtppb.Release{
				Id: int32(release.Id),
			})
	}

	for pageNumber := 2; pageNumber <= pagination.Pages; pageNumber++ {
		releases, _, err := b.pullUserCollection(user, pageNumber)
		if err != nil {
			return nil, err
		}

		for _, release := range releases {
			pr = append(pr,
				&drtppb.Release{
					Id: int32(release.Id),
				})
		}
	}

	return pr, err
}

func (b *Bridge) BuildMatcher(f string) *drtppb.Matcher {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatalf("Bad file read: %v (%v)", err, f)
	}

	lines := strings.Split(string(data), "\n")

	if lines[0] == "M" { // This is a file of master releases
		m := &drtppb.Matcher{Name: lines[1], SimpleName: lines[2]}
		for _, line := range lines[3:] {
			num, err := strconv.ParseInt(line, 10, 32)
			if err != nil {
				log.Fatalf("Bad num: %v", line)
			}
			releases, err := b.getSubReleases(int32(num))
			if err != nil {
				log.Fatalf("Bad read: %v", err)
			}

			match := &drtppb.Match{Release: releases}
			m.Matches = append(m.Matches, match)
		}

		return m
	}

	return nil
}

package builder

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

func GetTestBridge() *Bridge {
	b := &Bridge{}
	b.r = &testRetriever{}
	return b
}

type testRetriever struct{}

func clean(s string) string {
	return strings.Replace(
		strings.Replace(
			strings.Replace(
				s, "?", "_", -1),
			":", "_", -1),
		"/", "_", -1)
}

func (t *testRetriever) get(url string) ([]byte, error) {
	log.Printf("Getting %v", url)
	filename := fmt.Sprintf("test/%v", clean(url))
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		if strings.Contains(url, "master") {
			fmt.Printf("curl -s \"%v\" -o \"%v\"\nsleep 10\n", url, clean(url))
			str := `{"pagination": {
				"per_page": 50,
				"items": 4,
				"page": 1,
				"urls": {},
				"pages": 1
			  }
			  }`
			return []byte(str), nil
		}
	}

	return b, err
}

func TestGetReleases(t *testing.T) {
	b := GetTestBridge()

	releases, err := b.GetReleases(2228)
	if err != nil {
		t.Fatalf("Error in get releases: %v", err)
	}

	if len(releases) == 0 {
		t.Errorf("No releases returned: %v, %v", releases, err)
	}

	if len(releases) <= 50 {
		t.Errorf("Pagination has failed: %v", len(releases))
	}

	foundNIN := false
	foundDragnet := false
	for _, release := range releases {
		if release.Id == 5241 {
			foundNIN = true
		}

		if release.Id == 530182 && release.Title != "The Less You Look, The More You Find" {
			t.Errorf("Bad Release Title: %+v", release)
		}

		if release.Id == 371281 {
			foundDragnet = true
			if len(release.Tracks) != 11 {
				t.Errorf("Tracks not found: %+v", release)
			} else if release.Tracks[0].Title != "Psykick Dancehall" {
				t.Errorf("Wrong track title: %+v", release)
			}
		}
	}

	if foundNIN {
		t.Errorf("We've found Nine Inch Nails in the mix")
	}

	if !foundDragnet {
		t.Errorf("We did not find Dragnet")
	}
}

package builder

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func GetTestBridge() *Bridge {
	b := &Bridge{}
	b.r = &testRetriever{}
	return b
}

type testRetriever struct{}

func (t *testRetriever) get(url string) ([]byte, error) {
	filename := fmt.Sprintf("test/%v", strings.Replace(
		strings.Replace(
			strings.Replace(
				url, "?", "_", -1),
			":", "_", -1),
		"/", "_", -1))
	return ioutil.ReadFile(filename)
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
}

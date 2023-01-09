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
	filename := fmt.Sprintf("test/%v", strings.Replace(url,
		"/", "_", -1))
	return ioutil.ReadFile(filename)
}

func TestGetReleases(t *testing.T) {
	b := GetTestBridge()

	b.GetReleases(2228)
}

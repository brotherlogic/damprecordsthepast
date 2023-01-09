package builder

import "testing"

func TestGetReleases(t *testing.T) {
	b := &Bridge{}

	b.GetReleases(2228)
}

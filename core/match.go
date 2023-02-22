package core

import pb "github.com/brotherlogic/damprecordsthepast/proto"

func computeMatch(user *pb.User, matcher *pb.Matcher) float64 {
	matched := float64(0)
	count := float64(0)

	for _, entry := matcher.GetReleaseId() {
		found := false
		for _, release := range user.GetOwnedReleases() {
			if release == entry {
				found = true
				break
			}
		}

		if found {
			matched++
		}
		count++
	}

	return matched / count
}
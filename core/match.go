package core

import (
	pb "github.com/brotherlogic/damprecordsthepast/proto"
)

func ComputeMatch(user *pb.User, matcher *pb.Matcher) (float64, []*pb.Release) {
	matched := float64(0)
	count := float64(0)

	var unmatched []*pb.Release

	for _, entry := range matcher.GetMatches() {
		found := false
		for _, id := range entry.GetRelease() {
			for _, release := range user.GetOwnedReleases() {
				if release == id.GetId() {
					found = true
					break
				}
			}
		}

		if found {
			matched++
		} else {
			unmatched = append(unmatched, entry.GetRelease()...)
		}
		count++
	}

	return 100 * (matched / count), unmatched
}

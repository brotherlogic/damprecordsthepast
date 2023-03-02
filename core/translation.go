package core

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	pb "github.com/brotherlogic/damprecordsthepast/proto"
)

func Marshalmatcher(in *pb.Matcher) *pb.StoredMatcher {
	sm := &pb.StoredMatcher{Name: in.GetName(), SimpleName: in.GetSimpleName()}

	var strs []string
	for _, match := range in.GetMatches() {
		sort.SliceStable(match.ReleaseId, func(a, b int) bool {
			return match.ReleaseId[a] < match.ReleaseId[b]
		})

		str := fmt.Sprintf("%v", match.ReleaseId[0])
		for i := 1; i < len(match.ReleaseId); i++ {
			str += fmt.Sprintf(",%v", match.ReleaseId[i])
		}

		strs = append(strs, str)

	}

	sm.Matches = strings.Join(strs, ";")

	return sm
}

func UnmarshalMatcher(in *pb.StoredMatcher) *pb.Matcher {
	m := &pb.Matcher{Name: in.GetName(), SimpleName: in.GetSimpleName(), Matches: make([]*pb.Match, 0)}

	for _, match := range strings.Split(in.GetMatches(), ";") {
		var releases []int32
		for _, elem := range strings.Split(match, ",") {
			num, err := strconv.ParseInt(elem, 10, 32)
			if err != nil {
				panic(err)
			}

			releases = append(releases, int32(num))
		}

		m.Matches = append(m.Matches, &pb.Match{ReleaseId: releases})
	}

	return m
}

func convertToString(nums []int32) string {
	sort.SliceStable(nums, func(a, b int) bool {
		return nums[a] < nums[b]
	})

	str := fmt.Sprintf("%v", nums[0])
	for i := 1; i < len(nums); i++ {
		str += fmt.Sprintf(",%v", nums[i])
	}
	return str
}

func convertToNums(nums string) []int32 {
	var rnums []int32
	for _, elem := range strings.Split(nums, ",") {
		num, err := strconv.ParseInt(elem, 10, 32)
		if err != nil {
			panic(err)
		}

		rnums = append(rnums, int32(num))
	}

	return rnums
}

func MarshalUser(in *pb.User) *pb.StoredUser {
	return &pb.StoredUser{
		Name:          in.GetName(),
		Token:         in.GetToken(),
		ImageUrl:      in.GetImageUrl(),
		UserId:        in.GetUserId(),
		OwnedReleases: convertToString(in.GetOwnedReleases()),
	}
}

func UnmarshalUser(in *pb.StoredUser) *pb.User {
	return &pb.User{
		Name:          in.GetName(),
		Token:         in.GetToken(),
		ImageUrl:      in.GetImageUrl(),
		UserId:        in.GetUserId(),
		OwnedReleases: convertToNums(in.GetOwnedReleases()),
	}
}

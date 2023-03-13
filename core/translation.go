package core

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	pb "github.com/brotherlogic/damprecordsthepast/proto"
)

func convertForStorage(in string) string {
	return strings.ReplaceAll(strings.ReplaceAll(in, ";", "ASEMICOLON"), ",", "ACOMMA")
}
func convertFromStorage(in string) string {
	return strings.ReplaceAll(strings.ReplaceAll(in, "ASEMICOLON", ";"), "ACOMMA", ",")
}

func Marshalmatcher(in *pb.Matcher) *pb.StoredMatcher {
	sm := &pb.StoredMatcher{Name: in.GetName(), SimpleName: in.GetSimpleName()}

	var strs []string
	for _, match := range in.GetMatches() {
		sort.SliceStable(match.Release, func(a, b int) bool {
			return match.Release[a].GetId() < match.Release[b].GetId()
		})

		str := fmt.Sprintf("%v|%v", match.Release[0].GetId(), convertForStorage(match.Release[0].GetTitle()))
		for i := 1; i < len(match.Release); i++ {
			str += fmt.Sprintf(",%v|%v", match.Release[i].GetId(), convertForStorage(match.Release[i].GetTitle()))
		}

		strs = append(strs, str)

	}

	sm.Matches = strings.Join(strs, ";")

	return sm
}

func UnmarshalMatcher(in *pb.StoredMatcher) *pb.Matcher {
	m := &pb.Matcher{Name: in.GetName(), SimpleName: in.GetSimpleName(), Matches: make([]*pb.Match, 0)}

	for _, match := range strings.Split(in.GetMatches(), ";") {
		var releases []*pb.Release
		for _, elem := range strings.Split(match, ",") {
			splits := strings.Split(elem, "|")
			num, err := strconv.ParseInt(splits[0], 10, 32)
			if err != nil {
				panic(err)
			}

			releases = append(releases, &pb.Release{Id: int32(num), Title: convertFromStorage(splits[1])})
		}

		m.Matches = append(m.Matches, &pb.Match{Release: releases})
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

package core

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	pb "github.com/brotherlogic/damprecordsthepast/proto"
)

func Marshalmatcher(in *pb.Matcher) *pb.StoredMatcher {
	sm := &pb.StoredMatcher{Name: in.GetName()}

	sort.SliceStable(in.ReleaseId, func(a, b int) bool {
		return in.ReleaseId[a] < in.ReleaseId[b]
	})

	str := fmt.Sprintf("%v", in.ReleaseId[0])
	for i := 1; i < len(in.ReleaseId); i++ {
		str += fmt.Sprintf(",%v", in.ReleaseId[i])
	}
	sm.ReleaseIds = str

	return sm
}

func UnmarshalMatcher(in *pb.StoredMatcher) *pb.Matcher {
	m := &pb.Matcher{Name: in.GetName(), ReleaseId: make([]int32, 0)}

	for _, elem := range strings.Split(in.GetReleaseIds(), ",") {
		num, err := strconv.ParseInt(elem, 10, 32)
		if err != nil {
			panic(err)
		}

		m.ReleaseId = append(m.ReleaseId, int32(num))
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

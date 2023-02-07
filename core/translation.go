package core

import pb "github.com/brotherlogic/damprecordsthepast/proto"
import "sort"
import "fmt"
import "strings"
import "strconv"

func Marshal(in *pb.Matcher) *pb.StoredMatcher {
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

func Unmarshal(in *pb.StoredMatcher) *pb.Matcher {
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

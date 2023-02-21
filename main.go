package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pb "github.com/brotherlogic/damprecordsthepast/proto"

	"github.com/brotherlogic/damprecordsthepast/builder"
	"github.com/brotherlogic/damprecordsthepast/core"
	"github.com/brotherlogic/damprecordsthepast/remote"
	"github.com/brotherlogic/damprecordsthepast/webbuilder"
)

func main() {
	switch os.Args[1] {
	case "sync":
		bridge := builder.GetBridge()
		releases, err := bridge.GetReleases(2228)
		if err != nil {
			log.Fatalf("Unable to get releases: %v", err)
		}

		for _, release := range releases {
			fmt.Printf("%v\n", release.Id)
		}
	case "write":
		remote.Connect()
	case "build":
		err := webbuilder.BuildMatchPage([]*pb.User{{Name: "brotherlogic"}}, &pb.Matcher{Name: "Complete"})
		if err != nil {
			log.Fatalf("Unable to build website: %v", err)
		}
	case "store_full_match":
		bridge := builder.GetBridge()
		releases, err := bridge.GetReleases(78465) // Working on Swell Maps
		if err != nil {
			log.Fatalf("Unable to get releases: %v", err)
		}
		log.Printf("Got %v releases", len(releases))

		match := &pb.Matcher{
			Name: "full",
		}
		for _, release := range releases {
			match.ReleaseId = append(match.ReleaseId, release.GetId())
		}

		ctx := context.Background()
		remote := remote.Connect()
		err = remote.WriteMatcher(ctx, core.Marshalmatcher(match))
		fmt.Printf("Stored: %v\n", err)
	}
}

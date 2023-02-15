package main

import (
	"fmt"
	"log"
	"os"

	pb "github.com/brotherlogic/damprecordsthepast/proto"

	"github.com/brotherlogic/damprecordsthepast/builder"
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
	}
}

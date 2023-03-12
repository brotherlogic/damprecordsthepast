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
		ctx := context.Background()
		remote := remote.Connect()
		users, err := remote.GetUsers(ctx)
		if err != nil {
			log.Fatalf("Cannot get users: %v", err)
		}
		matchers, err := remote.GetMatchers(ctx)
		if err != nil {
			log.Fatalf("Cannot read matchers: %v", err)
		}
		for _, matcher := range matchers {
			err = webbuilder.BuildMatchPage(users, matcher)
			if err != nil {
				log.Fatalf("Unable to build website: %v", err)
			}
		}

		err = webbuilder.BuildIndexPage(users, matchers)
		if err != nil {
			log.Fatalf("Unable to build index: %v", err)
		}

		for _, user := range users {
			for _, matcher := range matchers {
				err = webbuilder.BuildUserMatchPage(user, matcher)
				if err != nil {
					log.Fatalf("Unable to build user match page for %v and %v", user.GetName(), matcher.GetName())
				}
			}
		}
	case "user":
		bridge := builder.GetBridge()
		releases, err := bridge.GetUserCollection("brotherlogic")
		if err != nil {
			log.Fatalf("Unable to pull collection: %v", err)
		}
		fmt.Printf("Found %v releases\n", len(releases))
	case "file":
		bridge := builder.GetBridge()
		matcher := bridge.BuildMatcher(os.Args[2])
		ctx := context.Background()
		remote := remote.Connect()
		err := remote.WriteMatcher(ctx, core.Marshalmatcher(matcher))
		fmt.Printf("Stored: %v\n", err)
	case "store_full_user":
		bridge := builder.GetBridge()
		releases, err := bridge.GetUserCollection("brotherlogic")
		if err != nil {
			log.Fatalf("Unable to pull collection: %v", err)
		}

		user := &pb.User{
			Name: "brotherlogic",
		}
		for _, release := range releases {
			user.OwnedReleases = append(user.OwnedReleases, release.GetId())
		}

		ctx := context.Background()
		remote := remote.Connect()
		err = remote.WriteUser(ctx, core.MarshalUser(user))
		fmt.Printf("Stored: %v\n", err)
	case "store_full_match":
		bridge := builder.GetBridge()
		releases, err := bridge.GetReleases(2228) // Working on Swell Maps
		if err != nil {
			log.Fatalf("Unable to get releases: %v", err)
		}
		log.Printf("Got %v releases", len(releases))

		match := &pb.Matcher{
			Name:       "All Releases",
			SimpleName: "full",
		}
		for _, release := range releases {
			match.Matches = append(match.Matches, &pb.Match{Release: []*pb.Release{release}})
		}

		ctx := context.Background()
		remote := remote.Connect()
		err = remote.WriteMatcher(ctx, core.Marshalmatcher(match))
		fmt.Printf("Stored: %v\n", err)
	default:
		fmt.Printf("Unknown command: %v\n", os.Args[1])
	}
}

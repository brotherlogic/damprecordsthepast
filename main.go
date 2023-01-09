package main

import (
	"fmt"
	"log"
	"os"

	"github.com/brotherlogic/damprecordsthepast/builder"
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
	}
}

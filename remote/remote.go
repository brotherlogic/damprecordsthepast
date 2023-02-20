package remote

import (
	"context"
	"fmt"
	"log"

	pb "github.com/brotherlogic/damprecordsthepast/proto"
	"google.golang.org/api/option"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

type Remote struct {
	store *firestore.Client
}

func (r *Remote) Close() {
	r.store.Close()
}

func Connect() *Remote {
	// Use the application default credentials
	ctx := context.Background()
	// Fetch the service account key JSON file contents
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	conf := &firebase.Config{ProjectID: "damprecordsthepast"}
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return &Remote{store: client}
}

func (r *Remote) WriteMatcher(ctx context.Context, matcher *pb.StoredMatcher) error {
	_, err := r.store.Collection("matchers").Doc(fmt.Sprintf("%v", matcher.GetName())).Set(ctx, matcher)
	return err
}

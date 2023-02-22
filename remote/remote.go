package remote

import (
	"context"
	"fmt"
	"log"

	"github.com/brotherlogic/damprecordsthepast/core"
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

func (r *Remote) GetUsers(ctx context.Context) ([]*pb.User, error) {
	res, err := r.store.Collection("users").Documents(ctx).GetAll()

	if err != nil {
		return nil, err
	}

	var users []*pb.User
	for _, user := range res {
		suser := &pb.StoredUser{}
		user.DataTo(suser)
		users = append(users, core.UnmarshalUser(suser))
	}

	return users, nil
}

func (r *Remote) GetMatcher(ctx context.Context, name string) (*pb.Matcher, error) {
	res, err := r.store.Collection("matchers").Doc(name).Get(ctx)
	if err != nil {
		return nil, err
	}

	matcher := &pb.StoredMatcher{}
	res.DataTo(matcher)
	return core.UnmarshalMatcher(matcher), nil
}

func (r *Remote) WriteUser(ctx context.Context, user *pb.StoredUser) error {
	_, err := r.store.Collection("users").Doc(fmt.Sprintf("%v", user.GetName())).Set(ctx, user)
	return err
}

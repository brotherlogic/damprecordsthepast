package remote

import (
	"context"
	"log"

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
	conf := &firebase.Config{ProjectID: "damprecordsthepast"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return &Remote{store: client}
}

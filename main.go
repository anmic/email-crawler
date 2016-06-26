package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"gopkg.in/pg.v4"

	"github.com/anmic/email-crawler/config"
	"github.com/anmic/email-crawler/models"
	"github.com/anmic/email-crawler/routes"
)

func newGmailService(ctx context.Context, config *oauth2.Config, db *pg.DB) (*gmail.Service, error) {
	var token models.Token
	err := db.Model(&token).Order("id DESC").Limit(1).Select()
	if err == pg.ErrNoRows {
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		fmt.Printf("Go to the following link in your browser then type the "+
			"authorization code: \n%v\n", authURL)

		done := make(chan bool)
		go func() {
			err := routes.GmailCallbackListener(config, db, done)
			if err != nil {
				log.Fatalf("ListenAndServe failed: %s", err)
			}
		}()
		<-done

		err = db.Model(&token).Order("id DESC").Limit(1).Select()
	}
	if err != nil {
		return nil, err
	}

	client := config.Client(ctx, token.Token())

	srv, err := gmail.New(client)
	return srv, err
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	db := pg.Connect(cfg.Postgres)
	err = models.CreateSchema(db)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	srv, err := newGmailService(ctx, cfg.Gmail, db)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	user := "anmic.testdev@gmail.com"
	msgsData, err := srv.Users.Messages.List(user).Do()
	if err != nil {
		panic(err)
	}

	for i, _ := range msgsData.Messages {
		msg, err := srv.Users.Messages.Get(user, msgsData.Messages[i].Id).Do()
		if err != nil {
			panic(err)
		}

		for _, header := range msg.Payload.Headers {
			if header.Name == "Return-Path" {

				lead := models.Lead{
					Email:     header.Value,
					CreatedAt: time.Now(),
				}

				_, err = db.Model(&lead).OnConflict("(email) DO NOTHING").Create()
				if err != nil {
					panic(err)
				}
				break
			}
		}
	}

	var leads []models.Lead
	err = db.Model(&leads).Select()
	if err != nil {
		panic(err)
	}
	log.Printf("leads : %+v \n", leads)
}

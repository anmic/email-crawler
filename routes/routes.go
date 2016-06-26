package routes

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"gopkg.in/pg.v4"

	"github.com/anmic/email-crawler/models"
)

func GmailCallbackListener(conf *oauth2.Config, db *pg.DB, done chan bool) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		fmt.Printf("\n code : %+v \n\n", code)
		if code == "" {
			return
		}

		token, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			panic(err)
		}

		tkn := models.Token{
			AccessToken:  token.AccessToken,
			TokenType:    token.TokenType,
			RefreshToken: token.RefreshToken,
			Expiry:       token.Expiry,
		}

		err = db.Create(&tkn)
		if err != nil {
			panic(err)
		}

		fmt.Printf("\n token : %+v \n\n", token)
		done <- true
	})

	err := http.ListenAndServe(":8080", nil)
	return err
}

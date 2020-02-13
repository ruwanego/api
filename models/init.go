package models

import (
	"context"

	"github.com/hiconvo/api/log"
	"github.com/hiconvo/api/utils/secrets"
)

var supportUser *User
var welcomeMessage string = `Hello, 👋

Welcome to Convo! Convo has two main features, **events** and **Convos**.

Convo events make it easy to plan events with real people. Invite your guests by name or email and they can RSVP in one click without having to create accounts of their own.

Convo also allows you to share content with people directly via *Convos*. A Convo is like a Facebook post except that it's private and only visible to the people you choose.

Read more about Convo and why I built it on [the blog](https://blog.convo.events/hello-world).

If you have any suggestions or feedback, please respond to this Convo directly and I'll get back to you.

Thanks,

Alex`

func init() {
	supportPw := secrets.Get("SUPPORT_PASSWORD", "support")
	ctx := context.Background()

	u, found, err := GetUserByEmail(ctx, "support@convo.events")
	if err != nil {
		panic(err)
	}

	if !found {
		u, err = NewUserWithPassword("support@convo.events", "Convo Support", "", supportPw)
		if err != nil {
			panic(err)
		}

		err = u.Commit(ctx)
		if err != nil {
			panic(err)
		}

		log.Print("models.init: Created new support user")
	}

	supportUser = &u
}

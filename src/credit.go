package main

// badge [name, merit-badge-id]
// score card [user(email), badge(badge), points(int), points_required(int), given(bool)]

/*
PLAN:
 when a user adds an experiment credit is notified with the user and the activity (added an experiment),
 credit checks if the is already a score card for that user+badge, If there is and they have already acheved it,
 the request is taken no further, if they have not yet acheved the badge one point is added. If this brings them up to
 the required amount of points for the badge, a new badge is requested from merit and offered to the user.

 The server will only be accessable inside the network.
*/

/*
IDEA:
 the user could be notified every time they earn a point (when points_required > 1), like with steme 'cheves'
*/

import (
	"github.com/codegangsta/martini"
	// "github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"labix.org/v2/mgo"
)

const (
	dbName = "credit"
)

type Badge struct {
	Name  string
	Badge int
}

type User struct {
	Name  string
	Email string
}

type Card struct {
	User     int
	Badge    int
	Points   int
	Required int
	Given    bool
}

// type Wish struct {
// 	Name        string `form:"name"`
// 	Description string `form:"description"`
// }

// DB Returns a martini.Handler
func DB() martini.Handler {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB(dbName))
		defer s.Close()
		c.Next()
	}
}

func GetAllCards(db *mgo.Database) []Card {
	var cards []Card
	db.C("cards").Find(nil).All(&cards)
	return cards
}

// GetCardsByUser returns all cards for a given user
func GetCardsByUser(db *mgo.Database, user string) []Card {
	var cards []Card
	db.C("cards").Find(nil).All(&cards)
	return cards
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(DB())

	//list cards
	m.Get("/cards", func(r render.Render, db *mgo.Database) {
		r.HTML(200, "list", GetAllCards(db))
	})

	//list cards by user(email)
	m.Get("/cards/email/:email", func(params martini.Params, r render.Render, db *mgo.Database) {
		// card := Card{1, 1, 1, 1, false}
		// db.C("cards").Insert(card)
		r.HTML(200, "list", GetCardsByUser(db, params["email"]))
	})

	// m.Post("/vards", binding.Form(Wish{}), func(wish Wish, r render.Render, db *mgo.Database) {
	// 	db.C("wishes").Insert(wish)
	// 	r.HTML(200, "list", GetAll(db))
	// })

	m.Run()
}

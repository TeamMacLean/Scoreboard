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
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"os"
	"strconv"
)

const (
	dbName     = "credit"
	cardsTable = "cards"
)

var (
	badges []Badge
	config Config
)

type Config struct {
	Authorization string
	BaseUrl       string
	/*
		Username string
		Password string
		"Basic "+base64(email + ":" + password)
	*/
}

type Badge struct {
	ID       int
	Name     string
	Badge    int
	Evidence string
}

type Card struct {
	Email    string
	Badge    int
	Points   int
	Required int
	Given    bool
}

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

func GenerateBadge() {
	/*
		required: recipient, evidence, badgeId
	*/
}

func UpdateCard(db *mgo.Database, email string, badge int) Card {
	card := CardByEmailAndBadge(db, email, badge)

	if len(card.Email) < 1 {
		err := db.C(cardsTable).Insert(&Card{email, badge, 0, 100, false})
		if err != nil {
			panic(err)
		}
	}

	if card.Given {
	} else {
		query := bson.M{"email": email, "badge": badge}
		change := bson.M{"$set": bson.M{"points": card.Points + 1}}
		err := db.C(cardsTable).Update(query, change)
		if err != nil {
			panic(err)
		}
		card = CardByEmailAndBadge(db, email, badge)
		if card.Points == card.Required {
			query := bson.M{"email": email, "badge": badge}
			change := bson.M{"$set": bson.M{"given": true}}
			err := db.C(cardsTable).Update(query, change)
			if err != nil {
				panic(err)
			}
		}
	}

	return CardByEmailAndBadge(db, email, badge)
}

func CardByEmailAndBadge(db *mgo.Database, email string, badge int) Card {
	card := Card{}
	db.C(cardsTable).Find(bson.M{"email": email, "badge": badge}).One(&card)
	return card
}

func GetAllCards(db *mgo.Database) []Card {
	var cards []Card
	err := db.C(cardsTable).Find(nil).All(&cards)
	if err != nil {
		panic(err)
	}
	return cards
}

func GetAllBadges() []Badge {
	return badges
}

func GetCardsByEmail(db *mgo.Database, email string) []Card {
	var cards []Card
	err := db.C(cardsTable).Find(bson.M{"Email": email}).All(&cards)
	if err != nil {
		panic(err)
	}
	return cards
}

func Auth(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("X-API-KEY") != "secret123" {
		res.WriteHeader(http.StatusUnauthorized)
	}
}

func LoadConfig() {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &config)
}

func LoadBadges() {
	file, err := ioutil.ReadFile("./badges.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &badges)
}

func main() {

	LoadConfig()
	LoadBadges()

	m := martini.Classic()
	m.Use(render.Renderer())
	// m.Use(Auth)
	m.Use(DB())

	m.Get("/", func() string {
		return "Merry Christmas!"
	})

	m.Get("/badges", func(r render.Render) {
		r.JSON(200, GetAllBadges())
	})

	m.Get("/cards", func(db *mgo.Database, r render.Render) {
		r.JSON(200, GetAllCards(db))
	})

	m.Get("/cards/update/:email/:badge", func(params martini.Params, db *mgo.Database, r render.Render) {

		email := params["email"]
		badge := params["badge"]
		i, err := strconv.Atoi(badge)
		if err != nil {
			panic(err)
		}
		out := UpdateCard(db, email, i)
		r.JSON(200, out)
	})
	m.Run()
}
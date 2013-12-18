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
	ID             int
	Name           string
	Badge          int
	Evidence       string
	PointsRequired int
}

type Card struct {
	Email  string
	Badge  int
	Points int
	// PointsRequired int
	Given bool
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

func GenerateBadge(recipient string, evidence string, badgeID string) string {

	fullurl := config.BaseUrl + "/api"

	client := &http.Client{}
	req, err := http.NewRequest("POST", fullurl, nil)
	if err != nil {
		return "error 1"
	}
	req.Header.Add("Authorization", config.Authorization)
	req.Header.Add("recipient", recipient)
	req.Header.Add("evidence", evidence)
	req.Header.Add("badgeId", badgeID)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	// if resp.StatusCode == 200 { // OK
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "error 2"
	}
	return string(bs)
	// }
	// return "error 3"
}

// func Auth(res http.ResponseWriter, req *http.Request) {
// 	if req.Header.Get("X-API-KEY") != "secret123" {
// 		res.WriteHeader(http.StatusUnauthorized)
// 	}
// }

func AccessControlAllowOrigin(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*") //tighten this up!
}

func UpdateCard(db *mgo.Database, email string, badge int) Card {
	card := CardByEmailAndBadge(db, email, badge)

	var currentBadge Badge
	var realBadge bool = false

	for i := range badges {
		if badges[i].ID == badge {
			realBadge = true
			currentBadge = badges[i]
		}
	}
	if realBadge {
	} else {
		print("no such badge")
		return Card{}
	}

	if len(card.Email) < 1 {
		err := db.C(cardsTable).Insert(&Card{email, badge, 0, false})
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
		card.Points += 1
		if card.Points == currentBadge.PointsRequired {
			query := bson.M{"email": email, "badge": badge}
			change := bson.M{"$set": bson.M{"given": true}}
			card.Given = true
			err := db.C(cardsTable).Update(query, change)
			if err != nil {
				panic(err)
			}
		}
	}

	return card
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

func GetBadgeByID(badgeID int) Badge {
	var currentBadge Badge
	for i := range badges {
		if badges[i].ID == badgeID {
			currentBadge = badges[i]
		}
	}
	return currentBadge
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
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
}

func LoadBadges() {
	file, err := ioutil.ReadFile("./badges.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	err = json.Unmarshal(file, &badges)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	if len(badges) < 1 {
		fmt.Println("could not create badges array")
		os.Exit(1)
	}
}

func main() {

	LoadConfig()
	LoadBadges()

	m := martini.Classic()
	m.Use(render.Renderer())
	// m.Use(Auth)
	m.Use(AccessControlAllowOrigin)
	m.Use(DB())

	m.Get("/", func() string {
		return "Merry Christmas!"
	})

	m.Get("/badges", func(r render.Render) {
		r.JSON(200, GetAllBadges())
	})
	m.Get("/badges/:id", func(params martini.Params, r render.Render) {
		badgeID := params["id"]
		i, err := strconv.Atoi(badgeID)
		if err != nil {
			panic(err)
		}
		r.JSON(200, GetBadgeByID(i))
	})

	m.Get("/cards", func(db *mgo.Database, r render.Render) {
		r.JSON(200, GetAllCards(db))
	})

	m.Get("/gen", func() string { return GenerateBadge() })

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
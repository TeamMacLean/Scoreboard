package main

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
	Given  bool
	Assert string
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

func GenerateBadge(recipient string, evidence string, badgeId int) string {

	badge := strconv.Itoa(badgeId)

	fmt.Println("creating badge for " + recipient + " : " + evidence + " : " + badge)

	fullurl := config.BaseUrl + "/api"

	client := &http.Client{}
	req, err := http.NewRequest("POST", fullurl, nil)
	if err != nil {
		fmt.Println("error 1")
		return ""
	}
	req.Header.Add("Authorization", config.Authorization)
	req.Header.Add("recipient", recipient)
	req.Header.Add("evidence", evidence)
	req.Header.Add("badgeId", badge)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode == 200 { // OK
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error 2")
			return ""
		}
		fmt.Println("returning %v ", bs)
		return string(bs)
	}
	fmt.Println("error 3")
	return ""
}

// this is a little hack to allow ajax access to the service
func AccessControlAllowOrigin(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*") //tighten this up!
	//possibly add an option in the config for url(s) accessing this.
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
		err := db.C(cardsTable).Insert(&Card{email, badge, 0, false, ""}) //gen new card
		if err != nil {
			panic(err)
		}
	}

	if card.Given {
		if len(card.Assert) > 0 {
			query := bson.M{"email": email, "badge": badge}
			change := bson.M{"$set": bson.M{"assert": ""}}
			err := db.C(cardsTable).Update(query, change)
			if err != nil {
				panic(err)
			}
		}

	} else {

		query := bson.M{"email": email, "badge": badge}
		change := bson.M{"$set": bson.M{"points": card.Points + 1}}
		err := db.C(cardsTable).Update(query, change)
		if err != nil {
			panic(err)
		}
		card.Points += 1
		if card.Points == currentBadge.PointsRequired {

			fullBadge := GetBadgeByID(badge)

			url := GenerateBadge(email, "http://www.google.com", fullBadge.Badge) //TODO fix the evidence feild

			query := bson.M{"email": email, "badge": badge}
			change := bson.M{"$set": bson.M{"given": true, "assert": url}}
			card.Assert = url
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

// i cant even remember is this is for anything
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

// unused by default
func Auth(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("X-API-KEY") != "secret123" {
		res.WriteHeader(http.StatusUnauthorized)
	}
}

// load the config.json file, if it does not exits or is not readable the app will exit
func LoadConfig() {

	var inputFile = "./config.json"

	if len(os.Args) == 3 {
		inputFile = os.Args[1]
		fmt.Printf("using %v for config file", inputFile)
	}

	file, err := ioutil.ReadFile(inputFile)

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

	var inputFile = "./badges.json"

	if len(os.Args) == 3 {
		inputFile = os.Args[2]
		fmt.Printf("using %v for badges file", inputFile)
	}

	file, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	err = json.Unmarshal(file, &badges)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	for i := range badges {
		fmt.Printf("loaded badge: %+v \n", badges[i])
	}

	if len(badges) < 1 {
		fmt.Println("could not create badges array")
		os.Exit(1)
	}
}

func main() {

	fmt.Println("Loading config...")
	LoadConfig()

	fmt.Println("Loading badge...")
	LoadBadges()

	m := martini.Classic()
	m.Use(render.Renderer())
	// m.Use(Auth)
	m.Use(AccessControlAllowOrigin)
	m.Use(DB())

	m.Get("/", func() string {
		return "Scoreboard is running!"
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

	fmt.Println("Starting Credit...")
	m.Run()
}

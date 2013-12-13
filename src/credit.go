package main

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"labix.org/v2/mgo"
)

const (
	dbName = "credit"
)

type Wish struct {
	Name        string `form:"name"`
	Description string `form:"description"`
}

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

// GetAll returns all Wishes in the database
func GetAll(db *mgo.Database) []Wish {
	var wishlist []Wish
	db.C("wishes").Find(nil).All(&wishlist)
	return wishlist
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(DB())

	m.Get("/", func(r render.Render, db *mgo.Database) {
		r.HTML(200, "index", GetAll(db))
	})

	m.Get("/badges", func(r render.Render, db *mgo.Database) {
		r.HTML(200, "list", GetAll(db))
	})

	m.Post("/badges", binding.Form(Wish{}), func(wish Wish, r render.Render, db *mgo.Database) {
		db.C("wishes").Insert(wish)
		r.HTML(200, "list", GetAll(db))
	})

	m.Run()
}

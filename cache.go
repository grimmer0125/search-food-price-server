package main

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	mongoURL   = "localhost:27017"
	collection = "QueryResults"
)

type QueryResult struct {
	// ID        bson.ObjectId `bson:"_id,omitempty"`
	Store           string    `bson:"store"`
	QueryKey        string    `bson:"queryKey"`
	Title           string    `bson:"title"`
	PreviewImageURL string    `bson:"previewImageURL"`
	Price           string    `bson:"price"`
	ProductID       float64   `bson:"productID"`
	CreatedAt       time.Time `bson:"createdAt"`
}

func setupMongo() {
	session, err := mgo.Dial(mongoURL)
	if err != nil {
		fmt.Printf("fail to connect to mongo")
		return
	}
	defer session.Close()
	c := session.DB("test").C(collection)

	// Ref: mongodb's index part on golang
	// https://gist.github.com/border/3489566
	// https://github.com/go-mgo/mgo/issues/480
	// https://github.com/go-mgo/mgo/blob/9a2573d4ae52a2bf9f5b7900a50e2f8bcceeb774/session_test.go#L3366
	// https://stackoverflow.com/questions/36720669/mgo-ttl-indexes-creation-to-selectively-delete-documents
	index := mgo.Index{
		Key: []string{"createdAt"},
		// Unique:      true,
		// DropDups:    true,
		// Background:  true,
		ExpireAfter: 60 * 24 * time.Minute,
	}

	err = c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func mongoRead(queryKey string) (result *QueryResult) {

	session, err := mgo.Dial(mongoURL)
	if err != nil {
		fmt.Printf("fail to connect to mongo")
		return
	}
	defer session.Close()
	c := session.DB("test").C(collection)

	var result2 QueryResult //	result1 := bson.M
	err = c.Find(bson.M{"queryKey": queryKey}).One(&result2)
	if err != nil {
		fmt.Printf("not found")
		return
	}

	return &result2
}

// https://hk.saowen.com/a/1eb464b5aef2ccface49aa405889b8b8853f77a2dd183bb6d7d6ecc04344f615
func mongoInsert(result *QueryResult) {
	session, err := mgo.Dial(mongoURL)
	if err != nil {
		fmt.Printf("fail to connect to mongo")
		return
	}

	defer session.Close()
	c := session.DB("test").C(collection)

	//https://stackoverflow.com/questions/43278696/golang-mgo-insert-or-update-not-working-as-expected
	_, err := c.Upsert(
		bson.M{"queryKey": result.QueryKey}, result,
	)

	if err != nil {
		panic("upsert error")
	}

	// err = c.Insert(result)
	// if err != nil {
	// 	panic("insert error")
	// }
}

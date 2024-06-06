package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"test.com/helper"
	"test.com/model"
)

type MongoDB struct {
	DatabaseName string
}

func NewMongo(dbname string) *MongoDB {
	return &MongoDB{
		DatabaseName: dbname,
	}
}

func (m MongoDB) Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		return nil, err
	}

	return client, nil
}

func InsertData(client *mongo.Client, file string, ctx context.Context) {
	collection := client.Database("IMDB").Collection("movies")

	datas := helper.ReadCSV(file, ',')

	id := 0
	for _, data := range datas[1:] {
		res, err := collection.InsertOne(ctx, model.IMDBMovie{
			Title:            data[1],
			ReleaseDate:      data[2],
			Genres:           data[3],
			OriginalLanguage: data[4],
			Overview:         data[5],
			VoteCount:        data[6],
			VoteAverage:      data[7],
		})

		if err != nil {
			log.Println(err.Error())
		}

		id++
		log.Printf("%v Data inserted with ID %v", id, res.InsertedID)
	}
}

// Select get all data from collection and store the data into destination

func (m *MongoDB) Select(ctx context.Context, client *mongo.Client, collection string, dest interface{}) {

	coll := client.Database(m.DatabaseName).Collection(collection)

	cursor, err := coll.Find(ctx, bson.D{})

	if err != nil {
		panic(err.Error())
	}

	err = cursor.All(ctx, dest)

	if err != nil {
		panic(err.Error())
	}
}

var skip int64 = 0
var limit int64 = 3000

func (m *MongoDB) SelectWithLimit(ctx context.Context, cl *mongo.Client, data chan []model.IMDBMovie) {
	coll := cl.Database("IMDB").Collection("movies")

	for {
		cursor, err := coll.Find(ctx, bson.D{}, &options.FindOptions{Skip: &skip, Limit: &limit})
		if err != nil {
			break
		}

		var value []model.IMDBMovie
		err = cursor.All(ctx, &value)

		if err != nil {
			break
		}

		data <- value

		skip += limit

		if int64(len(value)) < limit {
			break
		}
	}

	close(data)

}

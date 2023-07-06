package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Job struct {
	Title     string
	Company   string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	//check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to DB")

	// coll := client.Database("jobs").Collection("jobs")
	// newJobs := []interface{}{
	// 	Job{Title: "Senior Frontend Developer", Company: "Blockpour", Status: "Applied", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	// 	Job{Title: "Go Developer", Company: "Kidsloop", Status: "Interviewing", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	// 	Job{Title: "AWS Engineer", Company: "On the beach", Status: "Rejected", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	// }

	// result, err := coll.InsertMany(context.TODO(), newJobs)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Print(result)
	DB := "jobs"
	COLL := "jobs"

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Get("/get-jobs", func(c *fiber.Ctx) error {
		coll := client.Database(DB).Collection(COLL)
		cursor, err := coll.Find(context.TODO(), bson.D{})
		if err != nil {
			panic(err)
		}
		var results []Job
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}

		for _, result := range results {
			cursor.Decode(&result)
			if err != nil {
				panic(err)
			}
		}
		return c.JSON(results)
	})

	app.Post("/add-job", func(c *fiber.Ctx) error {
		payload := struct {
			Title   string `json:"title"`
			Company string `json:"company"`
			Status  string `json:"status"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		coll := client.Database(DB).Collection(COLL)
		newJob := Job{Title: payload.Title, Company: payload.Company, Status: payload.Status, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		result, err := coll.InsertOne(context.TODO(), newJob)
		if err != nil {
			panic(err)
		}
		return c.JSON(result)
	})

	app.Listen(":5000")
}

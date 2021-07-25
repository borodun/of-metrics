package collector

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	collection *mongo.Collection
)

type FunctionUsage struct {
	ID       string  `bson:"_id"`
	CpuUsage float64 `bson:"cpu_used"`
}

func loginMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("mongo_uri")
	if len(mongoURI) == 0 {
		log.Fatalln("Mongo URI isn't provided")
	}

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("faas-measurements").Collection("functions")
}

func getInfo() []FunctionUsage {
	groupStage := bson.D{{"$group", bson.D{{"_id", "$function_name"}, {"cpu_used", bson.D{{"$sum", "$cpu_used"}}}}}}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{groupStage})
	if err != nil {
		panic(err)
	}

	var usages []FunctionUsage
	if err = cursor.All(ctx, &usages); err != nil {
		panic(err)
	}

	return usages
}

func monitorFunctions() {
	go func() {
		for {
			usage := getInfo()
			for _, el := range usage {
				opts.WithLabelValues(el.ID).Set(el.CpuUsage)
			}

			time.Sleep(time.Second * 5)
		}
	}()
}

var (
	opts = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "of_cpu_usage",
		Help: "The total cpu usage",
	}, []string{
		"function_name",
	})
)

func Run() {
	loginMongo()
	monitorFunctions()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
)

type Hoge struct {
	CreatedAt time.Time `json:"createdAt"`
}

func main() {
	log.Println("Start Main")

	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: "sinmetal-go",
	})
	if err != nil {
		panic(err)
	}
	trace.RegisterExporter(exporter)

	ctx := context.Background()
	ctx, span := trace.StartSpan(ctx, "/main")
	defer span.End()

	//sc, err := NewSpannerClient(ctx, "projects/gcpug-public-spanner/instances/merpay-sponsored-instance/databases/sinmetal")
	//if err != nil {
	//	panic(err)
	//}

	http.HandleFunc("/spanner", SpannerSimpleQueryHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

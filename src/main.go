package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"vdbbench/src/weaviate"
)

type VDB interface {
	Prepare(ctx context.Context) error
	Import(ctx context.Context, input *json.Decoder) error
}

func main() {
	vdb, err := weaviate.New("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	if err := vdb.Prepare(ctx); err != nil {
		log.Fatal(err)
	}

	// Time(ctx, "import:", func() {
	// f, err := os.Open("./data/sources.json")
	// if err != nil {
	// log.Fatalf("open error: %v", err)
	// }
	// defer f.Close()

	// buff := bufio.NewReader(f)
	// decoder := json.NewDecoder(buff)

	// if err := vdb.Import(ctx, decoder); err != nil {
	// log.Fatal(err)
	// }
	// })

	Time(ctx, "query:", func() {
		f, err := os.Open("./data/queries.json")
		if err != nil {
			log.Fatalf("open error: %v", err)
		}
		defer f.Close()

		buff := bufio.NewReader(f)
		decoder := json.NewDecoder(buff)

		if err := vdb.Query(ctx, decoder); err != nil {
			log.Fatal(err)
		}
	})
}

func Time(ctx context.Context, info string, fn func()) {
	start := time.Now()
	defer func() {
		end := time.Now()
		dur := end.Sub(start)
		fmt.Println(info, "duration:", dur)
	}()

	fn()
}

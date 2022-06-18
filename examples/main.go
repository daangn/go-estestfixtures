package main

import (
	"context"
	"fmt"
	"time"

	"github.com/daangn/go-estestfixtures"
	"github.com/daangn/go-estestfixtures/escontainer"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	containerName := "examples-main-elasticsearch"

	fmt.Println("Step1 ------ ", "Run", containerName, "container")
	dsn, err := escontainer.NewDSNFromElasticsearchContainer(ctx, containerName)
	if err != nil {
		panic(err)
	}

	fmt.Println("Step2 ------ ", "Create loader to load dump files to DSN => ", dsn)
	loader, err := estestfixtures.NewLoader(
		ctx,
		dsn,
		estestfixtures.WithLimit(100),
		estestfixtures.WithTargetNames(
			"faqs_20220529",
			"posts",
		),
		estestfixtures.WithDirectory("../fixturedata"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Step3 ------ ", "Loading dump files DSN => ", dsn)
	if err = loader.Load(ctx); err != nil {
		panic(err)
	}

	fmt.Println("Step4 ------ ", "Create loader to dump files from DSN => ", dsn)
	loader, err = estestfixtures.NewLoader(
		ctx,
		"",
		estestfixtures.WithDirectory("../testdata/elasticsearch"),
		estestfixtures.WithLimit(100),
		estestfixtures.WithTargetNames(
			"faqs_20220529",
			"posts",
		),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Step5 ------ ", "Dump files from DSN => ", dsn)
	if err = loader.Dump(ctx); err != nil {
		panic(err)
	}
}

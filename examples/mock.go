package main

import (
	"context"
	"time"

	"github.com/daangn/go-estestfixtures"
	"github.com/daangn/go-estestfixtures/escontainer"
)

type MockESClient struct {
	Client   *ESClient
	Fixtures *estestfixtures.Loader
}

func NewMockESClient(ctx context.Context, fixturePath string) (*MockESClient, error) {
	dsn, err := escontainer.NewDSNFromElasticsearchContainer(ctx, "examples-sample-test")
	if err != nil {
		return nil, err
	}

	esClient, err := NewESClient(dsn, 20*time.Second, 20*time.Second)
	if err != nil {
		return nil, err
	}

	loader, err := estestfixtures.NewLoader(
		ctx,
		dsn,
		estestfixtures.WithDirectory(fixturePath),
		estestfixtures.WithLimit(100),
		estestfixtures.WithTargetNames(
			"faqs_20220529",
			"posts",
		),
	)
	if err != nil {
		return nil, err
	}
	if err = loader.Load(ctx); err != nil {
		return nil, err
	}

	return &MockESClient{
		Client:   esClient,
		Fixtures: loader,
	}, nil
}

func (m *MockESClient) Reset(ctx context.Context) error {
	return m.Fixtures.ClearElasticsearch(ctx)
}

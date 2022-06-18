package main

import (
	"context"
	"time"

	"github.com/olivere/elastic/v7"
)

type ESClient struct {
	es  *elastic.Client
	dsn string
}

func NewESClient(dsn string, initTimeout, maxTimeout time.Duration) (*ESClient, error) {
	es, err := elastic.NewClient(
		elastic.SetURL(dsn),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
		elastic.SetGzip(true),
		elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewExponentialBackoff(initTimeout, maxTimeout))),
	)
	if err != nil {
		return nil, err
	}

	return &ESClient{
		es:  es,
		dsn: dsn,
	}, nil
}

func (c *ESClient) Foo(ctx context.Context) error {
	_, _, err := c.es.Ping(c.dsn).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

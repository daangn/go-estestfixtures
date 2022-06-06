package esfixture

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"net"
	"net/http"
	"time"
)

func RunESContainer(containerName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	cont, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         containerName,
			Image:        "docker.elastic.co/elasticsearch/elasticsearch:7.15.1",
			ExposedPorts: []string{"9200/tcp"},
			Env: map[string]string{
				"discovery.type": "single-node",
			},
		},
		Started: true,
	})
	if err != nil {
		return "", err
	}
	port, err := cont.MappedPort(ctx, "9200/tcp")
	if err != nil {
		return "", err
	}

	if err = cont.Start(ctx); err != nil {
		return "", err
	}

	dsn := net.JoinHostPort("0.0.0.0", port.Port())
	for {
		<-time.After(2 * time.Second)
		if s, subErr := cont.State(ctx); subErr != nil {
		} else {
			if s.Running {
				_, subErr = http.Get(fmt.Sprintf("%s://%s/_cluster/health", "http", dsn))
				if subErr == nil {
					<-time.After(2 * time.Second)
					break
				}
			}
		}
	}

	return dsn, nil
}

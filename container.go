package estestfixtures

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

func runESContainer(containerName string) (string, error) {
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
	successCount := 0
	for {
		<-time.After(2 * time.Second)
		if s, subErr := cont.State(ctx); subErr != nil {
		} else {
			if s.ExitCode != 0 {
				return "", fmt.Errorf("[%s] container's exit code is %d", containerName, s.ExitCode)
			}
			if s.Running {
				_, subErr = http.Get(fmt.Sprintf("%s://%s/_cluster/health", "http", dsn))
				if subErr == nil {
					successCount = successCount + 1
					if successCount > 2 {
						break
					}
				}
			}
		}
	}

	return dsn, nil
}

package node

import (
	"fmt"
	"os"
	"strings"

	"github.com/cartesi/rollups-node/internal/services"
)

var ValidatorServices = []services.Service{
	*StateServer,
	*AdvanceRunner,
	*AuthorityClaimer,
	*Dispatcher,
	*GraphQLServer,
	*Indexer,
	*InspectServer,
}

var (
	AdvanceRunner = services.NewService(
		"advance-runner",
		"cartesi-rollups-advance-runner",
		healthcheckPort("advance-runner"),
	)
	AuthorityClaimer = services.NewService(
		"authority-claimer",
		"cartesi-rollups-authority-claimer",
		healthcheckPort("authority-claimer"),
	)
	Dispatcher = services.NewService(
		"dispatcher",
		"cartesi-rollups-dispatcher",
		healthcheckPort("dispatcher"),
	)
	GraphQLServer = services.NewService(
		"graphql-server",
		"cartesi-rollups-graphql-server",
		healthcheckPort("graphql-server"),
	)
	Indexer = services.NewService(
		"indexer",
		"cartesi-rollups-indexer",
		healthcheckPort("indexer"),
	)
	InspectServer = services.NewService(
		"inspect-server",
		"cartesi-rollups-inspect-server",
		healthcheckPort("inspect-server"),
	)
	StateServer = services.NewService(
		"state-server",
		"cartesi-rollups-state-server",
		stateServerHealthcheckPort(),
	)
)

func healthcheckPort(serviceName string) string {
	env := healthcheckEnv(serviceName)
	if port, ok := os.LookupEnv(env); ok {
		if serviceName == "state-server" {
			split := strings.Split(port, ":")
			return split[1]
		}
		return port
	} else {
		return "8080"
	}
}

func stateServerHealthcheckPort() string {
	if address, ok := os.LookupEnv("SS_SERVER_ADDRESS"); ok {
		split := strings.Split(address, ":")
		if len(split) > 1 {
			return split[1]
		}
	}
	return "50051"
}

func healthcheckEnv(serviceName string) string {
	suffix := "_HEALTHCHECK_PORT"
	if serviceName == "dispatcher" || serviceName == "authority-claimer" {
		suffix = "_HTTP_SERVER_PORT"
	}
	normalizedName := strings.Replace(serviceName, "-", "_", -1)
	return fmt.Sprintf("%s%s", strings.ToUpper(normalizedName), suffix)
}

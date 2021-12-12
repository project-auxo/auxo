package broker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"

	brokerConfig "github.com/project-auxo/auxo/olympus/internal/config"
	"github.com/project-auxo/auxo/olympus/logging"
	pb "github.com/project-auxo/auxo/oracle/proto"
)

const timeOut = time.Minute

var (
	log    = logging.Base()
	once   sync.Once
	client pb.OracleBackendServiceClient // Singleton
)

func GetOracleClient(cfg *brokerConfig.Config) (pb.OracleBackendServiceClient, context.Context) {
	once.Do(func() {
		conn, err := grpc.Dial(
			fmt.Sprintf(
				"%s:%d", cfg.Broker.BackendClient.Hostname, cfg.Broker.BackendClient.Port), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("dial: %v", err)
		}
		client = pb.NewOracleBackendServiceClient(conn)
	})
	ctx, cancel := context.WithTimeout(context.TODO(), timeOut)
	defer cancel()

	return client, ctx
}

package olympus

import (
	"fmt"
	"sync"

	"google.golang.org/grpc"

	hestiaCfg "github.com/project-auxo/auxo/hestia/internal/config"
	"github.com/project-auxo/auxo/olympus/logging"
	pb "github.com/project-auxo/auxo/olympus/proto/olympus"
)

var (
	log    = logging.Base()
	once   sync.Once
	client pb.OlympusFrontendServiceClient // Singleton
)

func GetClient(cfg *hestiaCfg.Config) pb.OlympusFrontendServiceClient {
	once.Do(func() {
		conn, err := grpc.Dial(
			fmt.Sprintf("%s:%d", cfg.Hestia.FrontendClient.Hostname,
				cfg.Hestia.FrontendClient.Port), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("dial: %v", err)
		}
		client = pb.NewOlympusFrontendServiceClient(conn)
	})
	return client
}

package olympus

import (
	"fmt"
	"sync"

	"google.golang.org/grpc"

	hestiaCfg "github.com/project-auxo/auxo/hestia/internal/config"
	"github.com/project-auxo/auxo/olympus/logging"
	pb "github.com/project-auxo/auxo/olympus/proto/hestia"
)

var (
	log    = logging.Base()
	client pb.HestiaFrontendServiceClient
	once   sync.Once
)

func initClient(cfg *hestiaCfg.Config) {
	once.Do(func() {
		conn, err := grpc.Dial(
			fmt.Sprintf("%s:%d", cfg.Hestia.FrontendClient.Hostname,
				cfg.Hestia.FrontendClient.Port), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("dial: %v", err)
		}
		client = pb.NewHestiaFrontendServiceClient(conn)
	})
}

func GetClient(cfg *hestiaCfg.Config) pb.HestiaFrontendServiceClient {
	initClient(cfg)
	return client
}

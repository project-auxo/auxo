package olympus

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/project-auxo/auxo/olympus/logging"
	pb "github.com/project-auxo/auxo/olympus/proto/hestia"
)

// Get the number of agents that are operating on Olympus
func GetNumberOfAgents(client pb.HestiaFrontendServiceClient) gin.HandlerFunc {
	log = logging.Base()
	return func(gctx *gin.Context) {
		ctx, cancel := context.WithTimeout(
			context.Background(), time.Duration(10)*time.Second)
		defer cancel()
		getNumberOfAgentsRep, err := client.GetNumberOfAgents(ctx, &pb.GetNumberOfAgentsReq{})
		if err != nil {
			gctx.String(
				http.StatusInternalServerError, "%v.GetNumberOfAgents(_) = _, %v", client, err)
		}
		gctx.JSON(http.StatusOK, getNumberOfAgentsRep.Number)
	}
}

package util

import (
	"fmt"

	discpb "github.com/project-auxo/auxo/olympus/proto/discovery"
	"google.golang.org/protobuf/proto"
)

func UnmarshalDiscoveryMessage(msg []byte) (msgProto *discpb.DiscoveryMessage, err error) {
	msgProto = &discpb.DiscoveryMessage{}
	if err = proto.Unmarshal(msg, msgProto); err != nil {
		err = fmt.Errorf("could not unmarshal the received msg: %v", err)
	}
	return
}

package grpc

import (
	"exam/internal/generated/grpc/service"
	"exam/internal/server/device"
	"exam/internal/server/image"
	"exam/internal/server/kt"
	"exam/internal/server/mri"
	"exam/internal/server/node"
	"exam/internal/server/node_segment"
	"exam/internal/server/segment"
	"exam/internal/services"
)

// из за эмбедина приходится делать приписку перед Handler
type Handler struct {
	device.DeviceHandler
	mri.MriHandler
	kt.KtHandler
	image.ImageHandler
	node_segment.NodeSegmentHandler
	node.NodeHandler
	segment.SegmentHandler

	service.UnsafeExamSrvServer
}

func New(
	services *services.Services,
) *Handler {
	deviceHandler := device.New(services)
	mriHandler := mri.New(services)
	ktHandler := kt.New(services)
	imageHandler := image.New(services)
	nodeSegmentHandler := node_segment.New(services)
	nodeHandler := node.New(services)
	segmentHandler := segment.New(services)

	return &Handler{
		DeviceHandler:      deviceHandler,
		MriHandler:         mriHandler,
		KtHandler:          ktHandler,
		ImageHandler:       imageHandler,
		NodeSegmentHandler: nodeSegmentHandler,
		NodeHandler:        nodeHandler,
		SegmentHandler:     segmentHandler,
	}
}

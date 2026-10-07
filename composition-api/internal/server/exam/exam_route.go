package exam

import (
	"composition-api/internal/server/exam/device"
	"composition-api/internal/server/exam/image"
	"composition-api/internal/server/exam/kt"
	mri "composition-api/internal/server/exam/mri"
	"composition-api/internal/server/exam/node"
	"composition-api/internal/server/exam/node_segment"
	"composition-api/internal/server/exam/segment"
	services "composition-api/internal/services"
)

type ExamRoute interface {
	segment.SegmentHandler
	node_segment.NodeSegmentHandler
	device.DeviceHandler
	image.ImageHandler
	node.NodeHandler
	mri.MriHandler
	kt.KTHandler
}

type examRoute struct {
	segment.SegmentHandler
	node_segment.NodeSegmentHandler
	device.DeviceHandler
	image.ImageHandler
	node.NodeHandler
	mri.MriHandler
	kt.KTHandler
}

func NewExamRoute(services *services.Services) ExamRoute {
	segmentHandler := segment.NewHandler(services)
	nodeSegmentHandler := node_segment.NewHandler(services)
	deviceHandler := device.NewHandler(services)
	imageHandler := image.NewHandler(services)
	nodeHandler := node.NewHandler(services)
	mriHandler := mri.NewHandler(services)
	ktHandler := kt.NewHandler(services)

	return &examRoute{
		SegmentHandler:     segmentHandler,
		NodeSegmentHandler: nodeSegmentHandler,
		DeviceHandler:      deviceHandler,
		ImageHandler:       imageHandler,
		NodeHandler:        nodeHandler,
		MriHandler:         mriHandler,
		KTHandler:          ktHandler,
	}
}

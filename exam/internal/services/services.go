package services

import (
	"exam/internal/repository"
	"exam/internal/services/device"
	"exam/internal/services/image"
	"exam/internal/services/kt"
	"exam/internal/services/mri"
	"exam/internal/services/node"
	"exam/internal/services/node_segment"
	"exam/internal/services/segment"
	"exam/internal/services/splitter"

	dbus "exam/internal/dbus/producers"
)

type Services struct {
	Device      device.Service
	Mri         mri.Service
	Image       image.Service
	Node        node.Service
	Segment     segment.Service
	NodeSegment node_segment.Service
	Splitter    splitter.Service
	Kt          kt.Service
}

func New(
	dao repository.DAO,
	dbus dbus.Producer,
) *Services {
	device := device.New(dao)
	mri := mri.New(dao)
	image := image.New(dao, dbus)
	node := node.New(dao)
	segment := segment.New(dao)
	nodeSegment := node_segment.New(dao)
	splitter := splitter.New()
	kt := kt.New(dao, dbus)

	return &Services{
		Device:      device,
		Mri:         mri,
		Image:       image,
		Node:        node,
		Segment:     segment,
		NodeSegment: nodeSegment,
		Splitter:    splitter,
		Kt:          kt,
	}
}

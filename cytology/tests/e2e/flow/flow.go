package flow

import (
	"context"

	pb "cytology/internal/generated/grpc/service"

	"github.com/google/uuid"
)

type FlowData struct {
	CytologyImageID       uuid.UUID
	CopiedCytologyImageID uuid.UUID
	ExternalID            uuid.UUID
	DoctorID              uuid.UUID
	PatientID             uuid.UUID
	OriginalImageID       uuid.UUID
	SegmentationGroupID   int32
	SegmentationID        int32
}

type Deps struct {
	Adapter pb.CytologySrvClient
}

type Flow interface {
	Do(ctx context.Context) (FlowData, error)
}

type flowfunc func(ctx context.Context, data FlowData) (FlowData, error)

type flowelem struct {
	flowfunc flowfunc
	next     *flowelem
}

func (f *flowelem) do(ctx context.Context, data FlowData) (FlowData, error) {
	flowRes, err := f.flowfunc(ctx, data)
	if err != nil {
		return FlowData{}, err
	}

	if f.next != nil {
		return f.next.do(ctx, flowRes)
	}
	return flowRes, nil
}

type flowfuncDepsInjector func(deps *Deps) flowfunc

type _flow struct {
	head *flowelem
}

func (f *_flow) Do(ctx context.Context) (FlowData, error) {
	return f.head.do(ctx, FlowData{})
}

func New(deps *Deps, flows ...flowfuncDepsInjector) Flow {
	if len(flows) == 0 {
		panic("flows is empty")
	}

	flowHead := &flowelem{}

	prevFlow := flowHead
	for _, fl := range flows {
		elem := &flowelem{flowfunc: fl(deps)}
		prevFlow.next = elem
		prevFlow = elem
	}

	return &_flow{head: flowHead.next}
}

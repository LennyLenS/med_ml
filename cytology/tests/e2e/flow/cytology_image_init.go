package flow

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	pb "cytology/internal/generated/grpc/service"
)

var CytologyImageInit flowfuncDepsInjector = func(deps *Deps) flowfunc {
	return func(ctx context.Context, data FlowData) (FlowData, error) {
		externalID := uuid.New()
		doctorID := uuid.New()
		patientID := uuid.New()

		resp, err := deps.Adapter.CreateCytologyImage(ctx, &pb.CreateCytologyImageIn{
			ExternalId:       externalID.String(),
			DoctorId:         doctorID.String(),
			PatientId:        patientID.String(),
			DiagnosticNumber: 1,
		})
		if err != nil {
			return FlowData{}, fmt.Errorf("create cytology image: %w", err)
		}

		data.CytologyImageID = uuid.MustParse(resp.Id)
		data.ExternalID = externalID
		data.DoctorID = doctorID
		data.PatientID = patientID

		return data, nil
	}
}

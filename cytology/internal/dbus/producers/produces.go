package dbus

import (
	"context"

	cytologyanalysisrequestedpb "cytology/internal/generated/dbus/produce/cytologyanalysisrequested"

	dbuslib "github.com/WantBeASleep/med_ml_lib/dbus"
)

type Producer interface {
	SendCytologyAnalysisRequested(ctx context.Context, msg *cytologyanalysisrequestedpb.AnalysisRequested) error
}

type producer struct {
	producerCytologyAnalysisRequested dbuslib.Producer[*cytologyanalysisrequestedpb.AnalysisRequested]
}

func New(
	producerCytologyAnalysisRequested dbuslib.Producer[*cytologyanalysisrequestedpb.AnalysisRequested],
) Producer {
	return &producer{
		producerCytologyAnalysisRequested: producerCytologyAnalysisRequested,
	}
}

func (a *producer) SendCytologyAnalysisRequested(ctx context.Context, msg *cytologyanalysisrequestedpb.AnalysisRequested) error {
	return a.producerCytologyAnalysisRequested.Send(ctx, msg)
}

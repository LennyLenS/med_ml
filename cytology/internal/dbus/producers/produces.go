package dbus

import (
	"context"

	cytologysplittedpb "cytology/internal/generated/dbus/produce/cytologysplitted"

	dbuslib "github.com/WantBeASleep/med_ml_lib/dbus"
)

type Producer interface {
	SendCytologySplitted(ctx context.Context, msg *cytologysplittedpb.CytologySplitted) error
}

type producer struct {
	producerCytologySplitted dbuslib.Producer[*cytologysplittedpb.CytologySplitted]
}

func New(
	producerCytologySplitted dbuslib.Producer[*cytologysplittedpb.CytologySplitted],
) Producer {
	return &producer{
		producerCytologySplitted: producerCytologySplitted,
	}
}

func (a *producer) SendCytologySplitted(ctx context.Context, msg *cytologysplittedpb.CytologySplitted) error {
	return a.producerCytologySplitted.Send(ctx, msg)
}

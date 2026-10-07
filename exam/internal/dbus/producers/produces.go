// время покажет, но пока выглядит будто бесполезный пакет, раз воткнул сюда дженерики, можно уж было
// вообще везде либу использовать просто пихая интерфейс

package dbus

import (
	"context"

	ktpreparedpb "exam/internal/generated/dbus/produce/ktprepared"
	mricompletepb "exam/internal/generated/dbus/produce/mricomplete"
	mrisplittedpb "exam/internal/generated/dbus/produce/mrisplitted"

	dbuslib "github.com/WantBeASleep/med_ml_lib/dbus"
)

type Producer interface {
	SendMriSplitted(ctx context.Context, msg *mrisplittedpb.MriSplitted) error
	SendMriComplete(ctx context.Context, msg *mricompletepb.MriComplete) error
	SendKtPrepared(ctx context.Context, msg *ktpreparedpb.KtPrepared) error
}

type producer struct {
	producerMriSplitted dbuslib.Producer[*mrisplittedpb.MriSplitted]
	producerMriComplete dbuslib.Producer[*mricompletepb.MriComplete]
	producerKtPrepared  dbuslib.Producer[*ktpreparedpb.KtPrepared]
}

func New(
	producerMriSplitted dbuslib.Producer[*mrisplittedpb.MriSplitted],
	producerMriComplete dbuslib.Producer[*mricompletepb.MriComplete],
	producerKtPrepared dbuslib.Producer[*ktpreparedpb.KtPrepared],
) Producer {
	return &producer{
		producerMriSplitted: producerMriSplitted,
		producerMriComplete: producerMriComplete,
		producerKtPrepared:  producerKtPrepared,
	}
}

func (a *producer) SendMriSplitted(ctx context.Context, msg *mrisplittedpb.MriSplitted) error {
	return a.producerMriSplitted.Send(ctx, msg)
}

func (a *producer) SendMriComplete(ctx context.Context, msg *mricompletepb.MriComplete) error {
	return a.producerMriComplete.Send(ctx, msg)
}

func (a *producer) SendKtPrepared(ctx context.Context, msg *ktpreparedpb.KtPrepared) error {
	return a.producerKtPrepared.Send(ctx, msg)
}

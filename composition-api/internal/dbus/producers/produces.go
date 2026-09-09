// время покажет, но пока выглядит будто бесполезный пакет, раз воткнул сюда дженерики, можно уж было
// вообще везде либу использовать просто пихая интерфейс

package producers

import (
	"context"

	"github.com/IBM/sarama"
	dbuslib "github.com/WantBeASleep/med_ml_lib/dbus"

	ktuploadpb "composition-api/internal/generated/dbus/produce/ktupload"
	mriuploadpb "composition-api/internal/generated/dbus/produce/mriupload"
	uziuploadpb "composition-api/internal/generated/dbus/produce/uziupload"
)

type Producer interface {
	SendUziUpload(ctx context.Context, msg *uziuploadpb.UziUpload) error
	SendMriUpload(ctx context.Context, msg *mriuploadpb.MriUpload) error
	SendKtUpload(ctx context.Context, msg *ktuploadpb.KtUpload) error
}

type producer struct {
	producerUziUpload dbuslib.Producer[*uziuploadpb.UziUpload]
	producerMriUpload dbuslib.Producer[*mriuploadpb.MriUpload]
	producerKtUpload  dbuslib.Producer[*ktuploadpb.KtUpload]
}

func New(
	client sarama.SyncProducer,
) Producer {
	producerUziUpload := dbuslib.NewProducer[*uziuploadpb.UziUpload](
		client,
		"uziupload",
	)
	producerMriUpload := dbuslib.NewProducer[*mriuploadpb.MriUpload](client, "mriupload")
	producerKtUpload := dbuslib.NewProducer[*ktuploadpb.KtUpload](client, "ktupload")

	return &producer{
		producerUziUpload: producerUziUpload,
		producerMriUpload: producerMriUpload,
		producerKtUpload:  producerKtUpload,
	}
}

func (a *producer) SendMriUpload(ctx context.Context, msg *mriuploadpb.MriUpload) error {
	return a.producerMriUpload.Send(ctx, msg)
}

func (a *producer) SendKtUpload(ctx context.Context, msg *ktuploadpb.KtUpload) error {
	return a.producerKtUpload.Send(ctx, msg)
}

func (a *producer) SendUziUpload(ctx context.Context, msg *uziuploadpb.UziUpload) error {
	return a.producerUziUpload.Send(ctx, msg)
}

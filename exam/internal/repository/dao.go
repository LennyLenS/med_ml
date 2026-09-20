package repository

import (
	"context"

	daolib "github.com/WantBeASleep/med_ml_lib/dao"
	"github.com/jmoiron/sqlx"
	minio "github.com/minio/minio-go/v7"

	"exam/internal/repository/device"
	"exam/internal/repository/echographic"
	"exam/internal/repository/image"
	"exam/internal/repository/kt"
	"exam/internal/repository/mri"
	"exam/internal/repository/node"
	"exam/internal/repository/segment"
)

type DAO interface {
	daolib.DAO
	NewFileRepo() FileRepo
	NewDeviceQuery(ctx context.Context) device.Repository
	NewMriQuery(ctx context.Context) mri.Repository
	NewKtQuery(ctx context.Context) kt.Repository
	NewImageQuery(ctx context.Context) image.Repository
	NewSegmentQuery(ctx context.Context) segment.Repository
	NewNodeQuery(ctx context.Context) node.Repository
	NewEchographicQuery(ctx context.Context) echographic.Repository
}

type dao struct {
	daolib.DAO

	s3       *minio.Client
	s3bucket string
}

func NewRepository(psql *sqlx.DB, s3 *minio.Client, s3bucket string) DAO {
	return &dao{
		DAO:      daolib.NewDao(psql),
		s3:       s3,
		s3bucket: s3bucket,
	}
}

// SS3
func (d *dao) NewFileRepo() FileRepo {
	return &fileRepo{
		s3:     d.s3,
		bucket: d.s3bucket,
	}
}

// POSTNIGRES
func (d *dao) NewDeviceQuery(ctx context.Context) device.Repository {
	deviceQuery := device.NewR()
	d.NewRepo(ctx, deviceQuery)

	return deviceQuery
}

func (d *dao) NewMriQuery(ctx context.Context) mri.Repository {
	mriQuery := mri.NewR()
	d.NewRepo(ctx, mriQuery)

	return mriQuery
}

func (d *dao) NewKtQuery(ctx context.Context) kt.Repository {
	ktQuery := kt.NewR()
	d.NewRepo(ctx, ktQuery)

	return ktQuery
}

func (d *dao) NewImageQuery(ctx context.Context) image.Repository {
	imageQuery := image.NewR()
	d.NewRepo(ctx, imageQuery)

	return imageQuery
}

func (d *dao) NewSegmentQuery(ctx context.Context) segment.Repository {
	segment := segment.NewRepo()
	d.NewRepo(ctx, segment)

	return segment
}

func (d *dao) NewNodeQuery(ctx context.Context) node.Repository {
	node := node.NewRepo()
	d.NewRepo(ctx, node)

	return node
}

func (d *dao) NewEchographicQuery(ctx context.Context) echographic.Repository {
	echographic := echographic.NewRepo()
	d.NewRepo(ctx, echographic)

	return echographic
}

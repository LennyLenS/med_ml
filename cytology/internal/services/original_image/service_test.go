package original_image

import (
	"context"
	"errors"
	"testing"
	"time"

	"cytology/internal/repository"
	cytologyrepo "cytology/internal/repository/cytology_image"
	cytologyentity "cytology/internal/repository/cytology_image/entity"
	repoentity "cytology/internal/repository/entity"
	originalrepo "cytology/internal/repository/original_image"
	originalentity "cytology/internal/repository/original_image/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type originalImageDAO struct {
	repository.DAO
	originals originalrepo.Repository
	cytology  cytologyrepo.Repository
}

func (d *originalImageDAO) NewOriginalImageQuery(context.Context) originalrepo.Repository {
	return d.originals
}

func (d *originalImageDAO) NewCytologyImageQuery(context.Context) cytologyrepo.Repository {
	return d.cytology
}

type originalImageQuery struct {
	originalrepo.Repository
	images map[uuid.UUID][]originalentity.OriginalImage
	err    error
}

func (q *originalImageQuery) GetOriginalImagesByCytologyID(id uuid.UUID) ([]originalentity.OriginalImage, error) {
	if q.err != nil {
		return nil, q.err
	}
	if images, ok := q.images[id]; ok {
		return images, nil
	}
	return nil, repoentity.ErrNotFound
}

type cytologyImageQuery struct {
	cytologyrepo.Repository
	images map[uuid.UUID]cytologyentity.CytologyImage
	err    error
}

func (q *cytologyImageQuery) GetCytologyImageByID(id uuid.UUID) (cytologyentity.CytologyImage, error) {
	if q.err != nil {
		return cytologyentity.CytologyImage{}, q.err
	}
	if image, ok := q.images[id]; ok {
		return image, nil
	}
	return cytologyentity.CytologyImage{}, repoentity.ErrNotFound
}

func TestGetOriginalImagesByCytologyID_Versions(t *testing.T) {
	rootID, prevID, currentID := uuid.New(), uuid.New(), uuid.New()
	original := originalentity.OriginalImage{
		Id: uuid.New(), CytologyID: rootID, ImagePath: "original/image/path",
		CreateDate: time.Now(), ViewedFlag: true,
	}
	ownOriginal := original
	ownOriginal.Id = uuid.New()
	ownOriginal.CytologyID = currentID
	queryErr := errors.New("database unavailable")

	for _, tc := range []struct {
		name        string
		originals   map[uuid.UUID][]originalentity.OriginalImage
		originalErr error
		cytologyErr error
		cycle       bool
		want        []originalentity.OriginalImage
		wantErr     error
	}{
		{name: "inherits through multiple versions", originals: map[uuid.UUID][]originalentity.OriginalImage{rootID: {original}}, want: []originalentity.OriginalImage{original}},
		{name: "prefers own original", originals: map[uuid.UUID][]originalentity.OriginalImage{rootID: {original}, currentID: {ownOriginal}}, want: []originalentity.OriginalImage{ownOriginal}},
		{name: "inherits nearest original", originals: map[uuid.UUID][]originalentity.OriginalImage{rootID: {original}, prevID: {ownOriginal}}, want: []originalentity.OriginalImage{ownOriginal}},
		{name: "inherits after empty result", originals: map[uuid.UUID][]originalentity.OriginalImage{currentID: {}, rootID: {original}}, want: []originalentity.OriginalImage{original}},
		{name: "no original in history", wantErr: repoentity.ErrNotFound},
		{name: "original query error", originalErr: queryErr, wantErr: queryErr},
		{name: "history query error", cytologyErr: queryErr, wantErr: queryErr},
		{name: "history cycle", cycle: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			images := map[uuid.UUID]cytologyentity.CytologyImage{
				rootID:    {Id: rootID},
				prevID:    {Id: prevID, PrevID: uuid.NullUUID{UUID: rootID, Valid: true}},
				currentID: {Id: currentID, PrevID: uuid.NullUUID{UUID: prevID, Valid: true}},
			}
			if tc.cycle {
				images[rootID] = cytologyentity.CytologyImage{Id: rootID, PrevID: uuid.NullUUID{UUID: currentID, Valid: true}}
			}
			svc := New(&originalImageDAO{
				originals: &originalImageQuery{images: tc.originals, err: tc.originalErr},
				cytology:  &cytologyImageQuery{images: images, err: tc.cytologyErr},
			})
			got, err := svc.GetOriginalImagesByCytologyID(context.Background(), currentID)
			if tc.cycle {
				require.ErrorContains(t, err, "cycle in cytology image history")
			} else if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, originalentity.OriginalImage{}.SliceToDomain(tc.want), got)
			}
		})
	}
}

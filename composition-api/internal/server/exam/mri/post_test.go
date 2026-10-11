package mri

import (
	"bytes"
	"context"
	"net/textproto"
	"testing"

	api "composition-api/internal/generated/http/api"
	"composition-api/internal/server/security"
	"composition-api/internal/services"
	mriSrv "composition-api/internal/services/mri"

	"github.com/google/uuid"
	ht "github.com/ogen-go/ogen/http"
	"github.com/stretchr/testify/require"
)

type uploadService struct {
	mriSrv.Service
	input *mriSrv.CreateMriArg
}

func (s *uploadService) Create(_ context.Context, input mriSrv.CreateMriArg) (uuid.UUID, error) {
	s.input = &input
	return uuid.New(), nil
}

func TestMriPostFormats(t *testing.T) {
	for _, tc := range []struct{ name, format, want string }{
		{"scan.dcm", "application/dicom", "application/dicom"},
		{"scan.DCM", "application/octet-stream", "application/dicom"},
		{"scan.dicom", "", "application/dicom"},
		{"scan", "application/dicom", "application/dicom"},
		{"scan.tiff", "image/tiff", "image/tiff"},
		{"scan.bin", "application/octet-stream", ""},
		{"scan.dcm", "image/png", ""},
	} {
		t.Run(tc.name+tc.format, func(t *testing.T) {
			srv := &uploadService{}
			h := NewHandler(&services.Services{MriService: srv})
			ctx := security.WithToken(context.Background(), security.Token{Id: uuid.New()})
			body := bytes.NewReader([]byte("decoded asynchronously"))
			header := make(textproto.MIMEHeader)
			header.Set("Content-Type", tc.format)
			response, err := h.MriPost(ctx, &api.MriPostReq{File: ht.MultipartFile{Name: tc.name, File: body, Header: header}})
			if tc.want == "" {
				require.Error(t, err)
				require.Nil(t, srv.input)
				return
			}
			require.NoError(t, err)
			require.IsType(t, &api.SimpleUuid{}, response)
			require.Equal(t, tc.want, srv.input.File.Header.Get("Content-Type"))
			position, err := body.Seek(0, 1)
			require.NoError(t, err)
			require.Zero(t, position)
		})
	}
}

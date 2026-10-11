package upload

import (
	"mime"
	"net/textproto"
	"path/filepath"
	"strings"

	ht "github.com/ogen-go/ogen/http"
)

func NormalizeFormat(file *ht.MultipartFile) string {
	format := file.Header.Get("Content-Type")
	if parsed, _, err := mime.ParseMediaType(format); err == nil {
		format = parsed
	}
	ext := strings.ToLower(filepath.Ext(file.Name))
	if (ext == ".dcm" || ext == ".dicom") && (format == "" || format == "application/octet-stream") {
		format = "application/dicom"
	}
	if file.Header == nil {
		file.Header = make(textproto.MIMEHeader)
	}
	file.Header.Set("Content-Type", format)
	return format
}

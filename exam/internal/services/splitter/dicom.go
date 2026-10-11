package splitter

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/amrshadid/go-dicom/dataset"
	"github.com/amrshadid/go-dicom/filebase"
	"github.com/amrshadid/go-dicom/filereader"
	"github.com/amrshadid/go-dicom/tag"

	"exam/internal/domain"
)

const maxDicomSamples = 64 * 1024 * 1024

func (s *service) DicomToTiff(file domain.File) (domain.File, error) {
	frames, err := decodeDicomFrames(file, true)
	if err != nil {
		return domain.File{}, err
	}
	data, err := encodeTiffFrames(frames)
	if err != nil {
		return domain.File{}, fmt.Errorf("encode dicom as tiff: %w", err)
	}
	return domain.File{Format: Tiff, Size: int64(len(data)), Buf: bytes.NewReader(data)}, nil
}

func decodeDicomFrames(file domain.File, grayscale16 bool) ([]image.Image, error) {
	reader, ok := file.Buf.(io.ReadSeeker)
	if !ok {
		return nil, fmt.Errorf("DICOM input must be seekable")
	}
	if file.Size > 512*1024*1024 {
		return nil, fmt.Errorf("DICOM exceeds 512 MiB")
	}
	parsed, err := filereader.ReadDICOMFile(filebase.NewFileReader(reader))
	if err != nil {
		return nil, fmt.Errorf("read dicom: %w", err)
	}
	ds := parsed.GetDataset()
	info, err := ds.GetPixelDataInfo()
	if err != nil {
		return nil, fmt.Errorf("read pixel metadata: %w", err)
	}
	if err := validateDicomPixels(info); err != nil {
		return nil, err
	}
	for _, t := range []tag.Tag{tag.New(0x0028, 0x3000), tag.New(0x0028, 0x3010), tag.New(0x5200, 0x9229), tag.New(0x5200, 0x9230)} {
		if _, ok := ds.Get(t); ok {
			return nil, fmt.Errorf("DICOM LUT or enhanced functional group sequences are not supported")
		}
	}
	window, err := dicomWindow(ds, info)
	if err != nil {
		return nil, err
	}
	pixels, err := ds.PixelArrayInterpreted()
	if err != nil {
		return nil, fmt.Errorf("decode dicom pixels: %w", err)
	}
	var frames []image.Image
	switch values := pixels.(type) {
	case [][][]uint8:
		frames, err = renderDicomFrames(values, info, window, grayscale16)
	case [][][]uint16:
		frames, err = renderDicomFrames(values, info, window, grayscale16)
	case [][][]int8:
		frames, err = renderDicomFrames(values, info, window, grayscale16)
	case [][][]int16:
		frames, err = renderDicomFrames(values, info, window, grayscale16)
	default:
		return nil, fmt.Errorf("unsupported DICOM pixel array %T", pixels)
	}
	if err != nil {
		return nil, err
	}
	return frames, nil
}

func validateDicomPixels(info *dataset.PixelDataInfo) error {
	if info.Rows <= 0 || info.Columns <= 0 || info.NumberOfFrames <= 0 || (info.SamplesPerPixel != 1 && info.SamplesPerPixel != 3) {
		return fmt.Errorf("invalid DICOM frame dimensions or samples per pixel")
	}
	n := 1
	for _, size := range []int{info.Rows, info.Columns, info.NumberOfFrames, info.SamplesPerPixel} {
		if size > maxDicomSamples/n {
			return fmt.Errorf("DICOM exceeds %d decoded samples", maxDicomSamples)
		}
		n *= size
	}
	if (info.BitsAllocated != 1 && info.BitsAllocated != 8 && info.BitsAllocated != 16) || info.BitsStored <= 0 || info.BitsStored > info.BitsAllocated || info.HighBit != info.BitsStored-1 || (info.PixelRepresentation != 0 && info.PixelRepresentation != 1) {
		return fmt.Errorf("unsupported DICOM pixel representation")
	}
	switch strings.TrimSpace(info.PhotometricInterpretation) {
	case "MONOCHROME1", "MONOCHROME2":
		if info.SamplesPerPixel == 1 {
			return nil
		}
	case "RGB", "YBR_FULL", "YBR_FULL_422":
		if info.SamplesPerPixel == 3 && info.BitsAllocated == 8 && info.PixelRepresentation == 0 && (info.PlanarConfiguration == 0 || info.PlanarConfiguration == 1) {
			return nil
		}
	}
	return fmt.Errorf("unsupported DICOM photometric interpretation %q", info.PhotometricInterpretation)
}

type pixelWindow struct {
	slope, intercept, low, high float64
}

func dicomWindow(ds *dataset.Dataset, info *dataset.PixelDataInfo) (pixelWindow, error) {
	w := pixelWindow{slope: 1, high: float64(uint64(1)<<info.BitsStored) - 1}
	if info.PixelRepresentation == 1 {
		w.low = -float64(uint64(1) << (info.BitsStored - 1))
		w.high = -w.low - 1
	}
	var err error
	if w.slope, err = dicomFloat(ds, "RescaleSlope", 1); err != nil {
		return w, err
	}
	if w.intercept, err = dicomFloat(ds, "RescaleIntercept", 0); err != nil {
		return w, err
	}
	if w.slope == 0 {
		return w, fmt.Errorf("DICOM rescale slope must not be zero")
	}
	w.low, w.high = w.low*w.slope+w.intercept, w.high*w.slope+w.intercept
	if w.low > w.high {
		w.low, w.high = w.high, w.low
	}
	center, width := ds.GetStringByKeyword("WindowCenter"), ds.GetStringByKeyword("WindowWidth")
	if center != "" || width != "" {
		if center == "" || width == "" {
			return w, fmt.Errorf("DICOM window center and width must both be present")
		}
		function := strings.TrimSpace(ds.GetStringByKeyword("VOILUTFunction"))
		if function != "" && function != "LINEAR" {
			return w, fmt.Errorf("unsupported DICOM VOI LUT function %q", function)
		}
		c, err := dicomFloat(ds, "WindowCenter", 0)
		if err != nil {
			return w, err
		}
		v, err := dicomFloat(ds, "WindowWidth", 0)
		if err != nil || v < 1 {
			return w, fmt.Errorf("invalid DICOM window width %q", width)
		}
		w.low, w.high = c-0.5-(v-1)/2, c-0.5+(v-1)/2
	}
	return w, nil
}

func dicomFloat(ds *dataset.Dataset, keyword string, fallback float64) (float64, error) {
	value := strings.TrimSpace(strings.Split(ds.GetStringByKeyword(keyword), "\\")[0])
	if value == "" {
		return fallback, nil
	}
	f, err := strconv.ParseFloat(value, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, fmt.Errorf("invalid DICOM %s %q", keyword, value)
	}
	return f, nil
}

func (w pixelWindow) grayValue(sample, maximum float64) float64 {
	v := sample*w.slope + w.intercept
	if v <= w.low {
		return 0
	}
	if v >= w.high {
		return maximum
	}
	return math.Round((v - w.low) * maximum / (w.high - w.low))
}

func (w pixelWindow) gray(sample float64) uint8 {
	return uint8(w.grayValue(sample, 255))
}

func renderDicomFrames[T ~uint8 | ~uint16 | ~int8 | ~int16](values [][][]T, info *dataset.PixelDataInfo, window pixelWindow, grayscale16 bool) ([]image.Image, error) {
	if len(values) != info.NumberOfFrames {
		return nil, fmt.Errorf("DICOM frame count mismatch")
	}
	frames := make([]image.Image, len(values))
	photo := strings.TrimSpace(info.PhotometricInterpretation)
	for f, rows := range values {
		if len(rows) != info.Rows {
			return nil, fmt.Errorf("DICOM row count mismatch in frame %d", f)
		}
		img := image.NewNRGBA(image.Rect(0, 0, info.Columns, info.Rows))
		var gray *image.Gray16
		if grayscale16 {
			gray = image.NewGray16(img.Bounds())
		}
		for y, row := range rows {
			if len(row) != info.Columns*info.SamplesPerPixel {
				return nil, fmt.Errorf("DICOM column count mismatch in frame %d", f)
			}
			for x := 0; x < info.Columns; x++ {
				if grayscale16 && info.SamplesPerPixel == 1 {
					value := uint16(window.grayValue(float64(row[x]), 65535))
					if photo == "MONOCHROME1" {
						value = 65535 - value
					}
					gray.SetGray16(x, y, color.Gray16{Y: value})
					continue
				}
				var r, g, b uint8
				if info.SamplesPerPixel == 1 {
					r = window.gray(float64(row[x]))
					if photo == "MONOCHROME1" {
						r = 255 - r
					}
					g, b = r, r
				} else {
					r, g, b = uint8(row[x*3]), uint8(row[x*3+1]), uint8(row[x*3+2])
					if photo != "RGB" {
						r, g, b = color.YCbCrToRGB(r, g, b)
					}
				}
				pixel := color.NRGBA{R: r, G: g, B: b, A: 255}
				if grayscale16 {
					gray.SetGray16(x, y, color.Gray16Model.Convert(pixel).(color.Gray16))
				} else {
					img.SetNRGBA(x, y, pixel)
				}
			}
		}
		frames[f] = img
		if grayscale16 {
			frames[f] = gray
		}
	}
	return frames, nil
}

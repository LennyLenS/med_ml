package splitter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/amrshadid/go-dicom/compress"
	"github.com/chai2010/tiff"
	"github.com/stretchr/testify/require"

	"exam/internal/domain"
)

func dicomFixture(rows, columns, frames, samples, bits, stored, representation, planar int, photo string, pixels []byte) []byte {
	b := bytes.NewBuffer(make([]byte, 128))
	b.WriteString("DICM")
	writeElement := func(group, element uint16, vr string, value []byte) {
		if len(value)%2 != 0 {
			value = append(append([]byte(nil), value...), ' ')
			if vr == "UI" || vr == "OB" {
				value[len(value)-1] = 0
			}
		}
		binary.Write(b, binary.LittleEndian, group)
		binary.Write(b, binary.LittleEndian, element)
		b.WriteString(vr)
		if vr == "OW" || vr == "OB" {
			binary.Write(b, binary.LittleEndian, uint16(0))
			binary.Write(b, binary.LittleEndian, uint32(len(value)))
		} else {
			binary.Write(b, binary.LittleEndian, uint16(len(value)))
		}
		b.Write(value)
	}
	writeElement(0x0002, 0x0010, "UI", []byte("1.2.840.10008.1.2.1"))
	writeElement(0x0028, 0x0004, "CS", []byte(photo))
	writeElement(0x0028, 0x0008, "IS", []byte(fmt.Sprint(frames)))
	for _, item := range [][2]int{{0x0002, samples}, {0x0006, planar}, {0x0010, rows}, {0x0011, columns}, {0x0100, bits}, {0x0101, stored}, {0x0102, stored - 1}, {0x0103, representation}} {
		value := make([]byte, 2)
		binary.LittleEndian.PutUint16(value, uint16(item[1]))
		writeElement(0x0028, uint16(item[0]), "US", value)
	}
	writeElement(0x7fe0, 0x0010, "OW", pixels)
	return b.Bytes()
}

func TestDicomToTiffAndPng(t *testing.T) {
	tests := []struct {
		name                                                  string
		frames, samples, bits, stored, representation, planar int
		photo                                                 string
		pixels                                                []byte
		want                                                  [][]color.NRGBA
	}{
		{"single grayscale", 1, 1, 8, 8, 0, 0, "MONOCHROME2", []byte{0, 255}, [][]color.NRGBA{{{A: 255}, {R: 255, G: 255, B: 255, A: 255}}}},
		{"multiple grayscale", 2, 1, 8, 8, 0, 0, "MONOCHROME2", []byte{0, 255, 64, 128}, [][]color.NRGBA{{{A: 255}, {R: 255, G: 255, B: 255, A: 255}}, {{R: 64, G: 64, B: 64, A: 255}, {R: 128, G: 128, B: 128, A: 255}}}},
		{"inverted grayscale", 1, 1, 8, 8, 0, 0, "MONOCHROME1", []byte{0, 255}, [][]color.NRGBA{{{R: 255, G: 255, B: 255, A: 255}, {A: 255}}}},
		{"12 bit grayscale", 1, 1, 16, 12, 0, 0, "MONOCHROME2", []byte{0, 0, 255, 15}, [][]color.NRGBA{{{A: 255}, {R: 255, G: 255, B: 255, A: 255}}}},
		{"signed grayscale", 1, 1, 16, 12, 1, 0, "MONOCHROME2", []byte{0, 8, 255, 7}, [][]color.NRGBA{{{A: 255}, {R: 255, G: 255, B: 255, A: 255}}}},
		{"rgb", 1, 3, 8, 8, 0, 0, "RGB", []byte{255, 0, 0, 0, 255, 0}, [][]color.NRGBA{{{R: 255, A: 255}, {G: 255, A: 255}}}},
		{"planar rgb", 1, 3, 8, 8, 0, 1, "RGB", []byte{255, 0, 0, 255, 0, 0}, [][]color.NRGBA{{{R: 255, A: 255}, {G: 255, A: 255}}}},
		{"ybr", 1, 3, 8, 8, 0, 0, "YBR_FULL", []byte{0, 128, 128, 255, 128, 128}, [][]color.NRGBA{{{A: 255}, {R: 255, G: 255, B: 255, A: 255}}}},
		{"ybr422", 1, 3, 8, 8, 0, 0, "YBR_FULL_422", []byte{0, 255, 128, 128}, [][]color.NRGBA{{{A: 255}, {R: 255, G: 255, B: 255, A: 255}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := dicomFixture(1, 2, tt.frames, tt.samples, tt.bits, tt.stored, tt.representation, tt.planar, tt.photo, tt.pixels)
			converted, err := New().DicomToTiff(domain.File{Format: Dicom, Size: int64(len(data)), Buf: bytes.NewReader(data)})
			require.NoError(t, err)
			require.Equal(t, Tiff, converted.Format)
			encoded, err := io.ReadAll(converted.Buf)
			require.NoError(t, err)
			require.Equal(t, int64(len(encoded)), converted.Size)
			pages, pageErrors, err := tiff.DecodeAll(bytes.NewReader(encoded))
			require.NoError(t, err)
			require.Len(t, pages, tt.frames)
			for i, page := range pages {
				require.Len(t, page, 1)
				require.NoError(t, pageErrors[i][0])
				for x, want := range tt.want[i] {
					require.Equal(t, color.NRGBAModel.Convert(color.Gray16Model.Convert(want)), color.NRGBAModel.Convert(page[0].At(x, 0)))
				}
			}
			converted.Buf = bytes.NewReader(encoded)
			pngs, err := New().SplitToPng(converted)
			require.NoError(t, err)
			require.Len(t, pngs, tt.frames)
			for i, file := range pngs {
				img, err := png.Decode(file.Buf)
				require.NoError(t, err)
				for x, want := range tt.want[i] {
					require.Equal(t, color.NRGBAModel.Convert(color.Gray16Model.Convert(want)), color.NRGBAModel.Convert(img.At(x, 0)))
				}
			}
		})
	}
}

func TestDicomToTiffRejectsInvalidInput(t *testing.T) {
	for name, data := range map[string][]byte{
		"broken":               []byte("not a dicom"),
		"truncated pixels":     dicomFixture(2, 2, 1, 1, 8, 8, 0, 0, "MONOCHROME2", []byte{1, 2}),
		"unsupported color":    dicomFixture(1, 2, 1, 1, 8, 8, 0, 0, "PALETTE COLOR", []byte{1, 2}),
		"excessive dimensions": dicomFixture(65535, 65535, 1, 1, 8, 8, 0, 0, "MONOCHROME2", []byte{1, 2}),
	} {
		t.Run(name, func(t *testing.T) {
			_, err := New().DicomToTiff(domain.File{Format: Dicom, Size: int64(len(data)), Buf: bytes.NewReader(data)})
			require.Error(t, err)
		})
	}
}

func TestDicomToTiffCompressedFrames(t *testing.T) {
	for _, syntax := range []string{"1.2.840.10008.1.2.5", "1.2.840.10008.1.2.4.50"} {
		t.Run(syntax, func(t *testing.T) {
			base := dicomFixture(1, 2, 2, 1, 8, 8, 0, 0, "MONOCHROME2", []byte{64, 64, 128, 128})
			var output bytes.Buffer
			output.Write(base[:132])
			uid := append([]byte(syntax), 0)
			if len(uid)%2 != 0 {
				uid = append(uid, 0)
			}
			output.Write([]byte{2, 0, 16, 0, 'U', 'I'})
			binary.Write(&output, binary.LittleEndian, uint16(len(uid)))
			output.Write(uid)
			oldMetaSize := 8 + int(binary.LittleEndian.Uint16(base[138:140]))
			pixelTag := bytes.Index(base, []byte{0xe0, 0x7f, 0x10, 0, 'O', 'W'})
			output.Write(base[132+oldMetaSize : pixelTag])
			output.Write([]byte{0xe0, 0x7f, 0x10, 0, 'O', 'B', 0, 0, 255, 255, 255, 255})
			writeItem := func(element uint16, data []byte) {
				if len(data)%2 != 0 {
					data = append(data, 0)
				}
				binary.Write(&output, binary.LittleEndian, uint16(0xfffe))
				binary.Write(&output, binary.LittleEndian, element)
				binary.Write(&output, binary.LittleEndian, uint32(len(data)))
				output.Write(data)
			}
			writeItem(0xe000, nil)
			for _, value := range []uint8{64, 128} {
				var encoded []byte
				if syntax == "1.2.840.10008.1.2.5" {
					var err error
					encoded, err = compress.NewRLECompressor().CompressFrame([]byte{value, value}, 1, 8)
					require.NoError(t, err)
				} else {
					img := image.NewGray(image.Rect(0, 0, 2, 1))
					img.Pix[0], img.Pix[1] = value, value
					var buffer bytes.Buffer
					require.NoError(t, jpeg.Encode(&buffer, img, &jpeg.Options{Quality: 100}))
					encoded = buffer.Bytes()
				}
				writeItem(0xe000, encoded)
			}
			writeItem(0xe0dd, nil)
			converted, err := New().DicomToTiff(domain.File{Format: Dicom, Size: int64(output.Len()), Buf: bytes.NewReader(output.Bytes())})
			require.NoError(t, err)
			pages, _, err := tiff.DecodeAll(converted.Buf)
			require.NoError(t, err)
			require.Len(t, pages, 2)
			for i, want := range []uint8{64, 128} {
				pixel := color.NRGBAModel.Convert(pages[i][0].At(0, 0)).(color.NRGBA)
				require.Equal(t, want, pixel.R)
			}
		})
	}
}

func TestDicomRejectsEnhancedFunctionalGroups(t *testing.T) {
	for _, element := range []uint16{0x9229, 0x9230} {
		data := dicomFixture(1, 2, 1, 1, 8, 8, 0, 0, "MONOCHROME2", []byte{0, 255})
		sequence := make([]byte, 12)
		binary.LittleEndian.PutUint16(sequence, 0x5200)
		binary.LittleEndian.PutUint16(sequence[2:], element)
		copy(sequence[4:], "SQ")
		data = append(data, sequence...)
		_, err := New().DicomToTiff(domain.File{Format: Dicom, Size: int64(len(data)), Buf: bytes.NewReader(data)})
		require.ErrorContains(t, err, "functional group")
	}
}

func TestDicomTiffPreservesGrayscalePrecision(t *testing.T) {
	data := dicomFixture(1, 2, 1, 1, 16, 12, 0, 0, "MONOCHROME2", []byte{1, 8, 2, 8})
	file, err := New().DicomToTiff(domain.File{Format: Dicom, Size: int64(len(data)), Buf: bytes.NewReader(data)})
	require.NoError(t, err)
	pages, _, err := tiff.DecodeAll(file.Buf)
	require.NoError(t, err)
	first := color.Gray16Model.Convert(pages[0][0].At(0, 0)).(color.Gray16)
	second := color.Gray16Model.Convert(pages[0][0].At(1, 0)).(color.Gray16)
	require.Greater(t, second.Y, first.Y)
	require.InDelta(t, 32800, first.Y, 20)
}

func TestPixelWindow(t *testing.T) {
	window := pixelWindow{slope: 2, intercept: -100, low: 0, high: 200}
	require.Equal(t, uint8(0), window.gray(50))
	require.Equal(t, uint8(128), window.gray(100))
	require.Equal(t, uint8(255), window.gray(150))
	threshold := pixelWindow{slope: 1, low: 5, high: 5}
	require.Equal(t, uint8(0), threshold.gray(5))
	require.Equal(t, uint8(255), threshold.gray(6))
}

func TestDicomSamples(t *testing.T) {
	dir := os.Getenv("DICOM_TEST_DIR")
	if dir == "" {
		t.Skip("DICOM_TEST_DIR is not set")
	}
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	count := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".dcm") {
			continue
		}
		count++
		t.Run(entry.Name(), func(t *testing.T) {
			file, err := os.Open(filepath.Join(dir, entry.Name()))
			require.NoError(t, err)
			defer file.Close()
			stat, err := file.Stat()
			require.NoError(t, err)
			converted, err := New().DicomToTiff(domain.File{Format: Dicom, Size: stat.Size(), Buf: file})
			require.NoError(t, err)
			pages, pageErrors, err := tiff.DecodeAll(converted.Buf)
			require.NoError(t, err)
			require.NotEmpty(t, pages)
			for i, page := range pages {
				require.Len(t, page, 1)
				require.NoError(t, pageErrors[i][0])
				require.False(t, page[0].Bounds().Empty())
			}
			_, err = converted.Buf.(io.Seeker).Seek(0, io.SeekStart)
			require.NoError(t, err)
			pngs, err := New().SplitToPng(converted)
			require.NoError(t, err)
			require.Len(t, pngs, len(pages))
			t.Logf("frames=%d width=%d height=%d tiff_bytes=%d", len(pages), pages[0][0].Bounds().Dx(), pages[0][0].Bounds().Dy(), converted.Size)
		})
	}
	require.Positive(t, count)
}

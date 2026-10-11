package splitter

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"math"
)

func encodeTiffFrames(frames []image.Image) ([]byte, error) {
	if len(frames) == 0 {
		return nil, fmt.Errorf("no TIFF frames")
	}
	const directorySize = 2 + 10*12 + 4
	size := uint64(8)
	for _, frame := range frames {
		if frame == nil || frame.Bounds().Empty() {
			return nil, fmt.Errorf("empty TIFF frame")
		}
		pageSize := uint64(directorySize) + uint64(frame.Bounds().Dx())*uint64(frame.Bounds().Dy())*2
		size += (pageSize + 1) &^ 1
	}
	if size > math.MaxUint32 {
		return nil, fmt.Errorf("TIFF exceeds 32-bit offsets")
	}
	data := make([]byte, int(size))
	copy(data, "II")
	order := binary.LittleEndian
	order.PutUint16(data[2:], 42)
	order.PutUint32(data[4:], 8)
	offset := uint32(8)
	for i, frame := range frames {
		width, height := uint32(frame.Bounds().Dx()), uint32(frame.Bounds().Dy())
		pixelsOffset := offset + directorySize
		pixelBytes := width * height * 2
		entries := [][4]uint32{
			{256, 4, 1, width},
			{257, 4, 1, height},
			{258, 3, 1, 16},
			{259, 3, 1, 1},
			{262, 3, 1, 1},
			{273, 4, 1, pixelsOffset},
			{277, 3, 1, 1},
			{278, 4, 1, height},
			{279, 4, 1, pixelBytes},
			{284, 3, 1, 1},
		}
		order.PutUint16(data[offset:], uint16(len(entries)))
		for j, entry := range entries {
			pos := int(offset) + 2 + j*12
			order.PutUint16(data[pos:], uint16(entry[0]))
			order.PutUint16(data[pos+2:], uint16(entry[1]))
			order.PutUint32(data[pos+4:], entry[2])
			order.PutUint32(data[pos+8:], entry[3])
		}
		pos := int(pixelsOffset)
		for y := frame.Bounds().Min.Y; y < frame.Bounds().Max.Y; y++ {
			for x := frame.Bounds().Min.X; x < frame.Bounds().Max.X; x++ {
				pixel := color.Gray16Model.Convert(frame.At(x, y)).(color.Gray16)
				order.PutUint16(data[pos:], pixel.Y)
				pos += 2
			}
		}
		next := (pixelsOffset + pixelBytes + 1) &^ 1
		if i+1 < len(frames) {
			order.PutUint32(data[offset+directorySize-4:], next)
		}
		offset = next
	}
	return data, nil
}

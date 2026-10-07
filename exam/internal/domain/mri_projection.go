package domain

import "fmt"

type MriProjection string

const (
	MriProjectionLong  MriProjection = "long"
	MriProjectionCross MriProjection = "cross"
)

func (s MriProjection) String() string {
	return string(s)
}

func (s MriProjection) Parse(projection string) (MriProjection, error) {
	switch projection {
	case "long":
		return MriProjectionLong, nil
	case "cross":
		return MriProjectionCross, nil
	default:
		return "", fmt.Errorf("invalid projection: %s", projection)
	}
}

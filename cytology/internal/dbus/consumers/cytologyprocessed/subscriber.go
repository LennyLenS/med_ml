package cytologyprocessed

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/WantBeASleep/med_ml_lib/dbus"
	"github.com/google/uuid"

	"cytology/internal/domain"
	pb "cytology/internal/generated/dbus/consume/cytologyprocessed"
	"cytology/internal/services"
	"cytology/internal/services/segmentation"
	"cytology/internal/services/segmentation_group"
)

// Маппинг имен классов из GeoJSON в SegType
var classNameToSegType = map[string]domain.SegType{
	"Нормальная клетка":            domain.SegTypeNIL,
	"Клетка Гюртле":                domain.SegTypeNIR,
	"Макрофаг":                     domain.SegTypeNIM,
	"Скопление упорядоченное":      domain.SegTypeCNO,
	"Скопление неупорядоченное":    domain.SegTypeCGE,
	"Скопление микрофолликулярное": domain.SegTypeC2N,
	"Скопление папиллярное":        domain.SegTypeCPS,
	"Скопление фолликулярное":      domain.SegTypeCFC,
	"Скопление лимфоидное":         domain.SegTypeCLY,
	"Метастаз отсутствует":         domain.SegTypeSOS,
	"Метастаз сомнительный":        domain.SegTypeSDS,
	"Метастаз вероятный":           domain.SegTypeSMS,
	"Метастаз определенный":        domain.SegTypeSTS,
	"Метастаз папиллярный":         domain.SegTypeSPS,
	"Метастаз отсутствует (нет)":   domain.SegTypeSNM,
	"Метастаз тиреоидный":          domain.SegTypeSTM,
}

// Маппинг SegType в GroupType
var segTypeToGroupType = map[domain.SegType]domain.GroupType{
	domain.SegTypeNIL: domain.GroupTypeCE,
	domain.SegTypeNIR: domain.GroupTypeCE,
	domain.SegTypeNIM: domain.GroupTypeCE,
	domain.SegTypeCNO: domain.GroupTypeCL,
	domain.SegTypeCGE: domain.GroupTypeCL,
	domain.SegTypeC2N: domain.GroupTypeCL,
	domain.SegTypeCPS: domain.GroupTypeCL,
	domain.SegTypeCFC: domain.GroupTypeCL,
	domain.SegTypeCLY: domain.GroupTypeCL,
	domain.SegTypeSOS: domain.GroupTypeME,
	domain.SegTypeSDS: domain.GroupTypeME,
	domain.SegTypeSMS: domain.GroupTypeME,
	domain.SegTypeSTS: domain.GroupTypeME,
	domain.SegTypeSPS: domain.GroupTypeME,
	domain.SegTypeSNM: domain.GroupTypeME,
	domain.SegTypeSTM: domain.GroupTypeME,
}

type subscriber struct {
	services *services.Services
}

func New(
	services *services.Services,
) dbus.Consumer[*pb.CytologyProcessed] {
	return &subscriber{
		services: services,
	}
}

func getSegTypeFromClassName(className string) (domain.SegType, domain.GroupType, bool) {
	// Пробуем точное совпадение
	if segType, ok := classNameToSegType[className]; ok {
		groupType := segTypeToGroupType[segType]
		return segType, groupType, true
	}

	// Пробуем найти по частичному совпадению (без учета регистра)
	classNameLower := strings.ToLower(className)
	for name, segType := range classNameToSegType {
		if strings.Contains(strings.ToLower(name), classNameLower) || strings.Contains(classNameLower, strings.ToLower(name)) {
			groupType := segTypeToGroupType[segType]
			return segType, groupType, true
		}
	}

	return "", "", false
}

func (h *subscriber) Consume(ctx context.Context, message *pb.CytologyProcessed) error {
	// Валидация UUID
	cytologyID, err := uuid.Parse(message.CytologyId)
	if err != nil {
		return fmt.Errorf("cytology id is not uuid: %s", message.CytologyId)
	}

	if _, err := uuid.Parse(message.OriginalImageId); err != nil {
		return fmt.Errorf("original image id is not uuid: %s", message.OriginalImageId)
	}

	// Обрабатываем каждую feature из protobuf структуры
	featureCollection := message.GeojsonFeatures
	if featureCollection == nil {
		return errors.New("geojson_features is nil")
	}

	for _, feature := range featureCollection.Features {
		if feature == nil {
			continue
		}

		// Получаем имя класса из properties
		properties := feature.Properties
		if properties == nil {
			continue
		}

		classification := properties.Classification
		if classification == nil {
			continue
		}

		className := classification.Name
		if className == "" {
			continue
		}

		// Получаем тип сегмента по имени класса
		segType, groupType, found := getSegTypeFromClassName(className)
		if !found {
			// Пропускаем неизвестные классы
			continue
		}

		// Создаем группу сегментов
		groupArg := segmentation_group.CreateSegmentationGroupArg{
			CytologyID: cytologyID,
			SegType:    segType,
			GroupType:  groupType,
			IsAI:       true,
			Details:    nil,
		}

		groupID, err := h.services.SegmentationGroup.CreateSegmentationGroup(ctx, groupArg)
		if err != nil {
			return fmt.Errorf("create segmentation group: %w", err)
		}

		// Получаем geometry
		geometry := feature.Geometry
		if geometry == nil {
			continue
		}

		var points []domain.SegmentationPoint

		// Обрабатываем координаты в зависимости от типа geometry
		switch g := geometry.GeometryType.(type) {
		case *pb.Geometry_Point:
			if g.Point != nil {
				points = append(points, domain.SegmentationPoint{
					X:   int(g.Point.X),
					Y:   int(g.Point.Y),
					UID: 0, // Будет установлен при создании
				})
			}
		case *pb.Geometry_Polygon:
			if g.Polygon != nil {
				// Обрабатываем первое кольцо полигона (внешний контур)
				if len(g.Polygon.Rings) > 0 {
					ring := g.Polygon.Rings[0]
					for _, point := range ring.Points {
						points = append(points, domain.SegmentationPoint{
							X:   int(point.X),
							Y:   int(point.Y),
							UID: 0, // Будет установлен при создании
						})
					}
				}
			}
		}

		if len(points) == 0 {
			continue
		}

		// Создаем сегментацию с точками
		segArg := segmentation.CreateSegmentationArg{
			SegmentationGroupID: groupID,
			Points:              points,
		}

		_, err = h.services.Segmentation.CreateSegmentation(ctx, segArg)
		if err != nil {
			return fmt.Errorf("create segmentation: %w", err)
		}
	}

	return nil
}

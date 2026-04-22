import json
from confluent_kafka import Producer
from ml_service.internal.s3.s3 import S3
from ml_service.config.default import get_settings
import ml_service.internal.events.kafka_pb2 as pb_event

settings = get_settings()


def geojson_to_protobuf(geojson_data):
    """
    Конвертирует GeoJSON FeatureCollection в protobuf структуру.
    """
    feature_collection = pb_event.FeatureCollection()

    if "features" not in geojson_data:
        return feature_collection

    for feature_json in geojson_data["features"]:
        feature = pb_event.Feature()

        # Обрабатываем geometry
        geometry_json = feature_json.get("geometry", {})
        geometry_type = geometry_json.get("type", "")
        coordinates = geometry_json.get("coordinates", [])

        if geometry_type == "Point":
            if len(coordinates) >= 2:
                point = pb_event.Point(
                    x=float(coordinates[0]),
                    y=float(coordinates[1])
                )
                feature.geometry.point.CopyFrom(point)
        elif geometry_type == "Polygon":
            polygon = pb_event.Polygon()
            # Polygon coordinates: [[[x, y], ...], ...] - массив колец
            for ring_coords in coordinates:
                if isinstance(ring_coords, list) and len(ring_coords) > 0:
                    ring = pb_event.Ring()
                    for point_coords in ring_coords:
                        if isinstance(point_coords, list) and len(point_coords) >= 2:
                            point = pb_event.Point(
                                x=float(point_coords[0]),
                                y=float(point_coords[1])
                            )
                            ring.points.append(point)
                    polygon.rings.append(ring)
            feature.geometry.polygon.CopyFrom(polygon)

        # Обрабатываем properties
        properties_json = feature_json.get("properties", {})
        properties = pb_event.Properties()
        properties.is_locked = properties_json.get("isLocked", False)

        classification_json = properties_json.get("classification", {})
        if classification_json:
            classification = pb_event.Classification()
            classification.name = classification_json.get("name", "")

            color_json = classification_json.get("color", [])
            if isinstance(color_json, list) and len(color_json) >= 3:
                color = pb_event.Color(
                    r=int(color_json[0]),
                    g=int(color_json[1]),
                    b=int(color_json[2])
                )
                classification.color.CopyFrom(color)

            properties.classification.CopyFrom(classification)

        feature.properties.CopyFrom(properties)
        feature_collection.features.append(feature)

    return feature_collection


class cytologyUseCase:
    def __init__(self, store: S3):
        self.store = store

    def processCytologyImage(self, cytology_id: str, original_image_id: str):
        """
        Обрабатывает цитологическое изображение и возвращает GeoJSON FeatureCollection.
        Пока что это заглушка - в реальности здесь должна быть обработка через нейронную сеть.
        """
        print(f"Processing cytology image: cytology_id={cytology_id}, original_image_id={original_image_id}")

        # Загружаем изображение из S3
        print("Going to S3...")
        try:
            image_data = self.store.load(f"{cytology_id}/{original_image_id}")
            print(f"Image loaded from S3, size: {len(image_data) if image_data else 0}")
        except Exception as e:
            print(f"Error loading image from S3: {e}")
            # Для тестирования создаем пустой GeoJSON
            image_data = None

        # TODO: Здесь должна быть реальная обработка через нейронную сеть
        # Пока возвращаем пустой GeoJSON FeatureCollection
        # В реальности здесь должен быть вызов модели, аналогичный run_wsi_analysis_parallel
        geojson_features = {
            "type": "FeatureCollection",
            "features": []
        }

        # Если есть данные, можно добавить обработку
        if image_data:
            # Здесь должна быть обработка изображения
            # geojson_features = process_image_with_models(image_data)
            pass

        # Конвертируем GeoJSON в protobuf структуру
        feature_collection = geojson_to_protobuf(geojson_features)

        # Создаем сообщение для Kafka
        msg_event = pb_event.CytologyProcessed(
            cytology_id=cytology_id,
            original_image_id=original_image_id,
            geojson_features=feature_collection
        )

        content = msg_event.SerializeToString()

        # Отправляем в Kafka
        producer_config = {
            "bootstrap.servers": settings.kafka_host + ":" + str(settings.kafka_port)
        }
        producer = Producer(producer_config)

        producer.produce("cytologyprocessed", content)
        producer.flush()

        print(f"✅ Cytology processed and sent to Kafka: cytology_id={cytology_id}")

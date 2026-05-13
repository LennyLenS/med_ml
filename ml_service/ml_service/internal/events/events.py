from confluent_kafka import Consumer
import ml_service.internal.events.kafka_pb2 as pb
from ml_service.config.default import get_settings

settings = get_settings()


class EventsYo:
    def __init__(self, uzi, cytology=None):
        self.uzi = uzi
        self.cytology = cytology

    def run(self):
        consumer_config = {
            "bootstrap.servers": settings.kafka_host + ":" + str(settings.kafka_port),  # Адрес Kafka-брокера
            "group.id": "mlService",  # Имя consumer group
            "auto.offset.reset": "earliest",  # Начинать с самого начала, если оффсет не найден
            "security.protocol": "PLAINTEXT",  # Установка протокола безопасности на PLAINTEXT для отключения SASL
            "broker.version.fallback": "2.3.0",
        }

        consumer = Consumer(consumer_config)
        topics = ["uzisplitted"]
        if self.cytology:
            topics.append("cytologysplitted")
        consumer.subscribe(topics)

        while True:
            msg = consumer.poll(timeout=1.0)
            if msg is None:
                continue  # Если сообщения нет, то пропускаем итерацию

            topic = msg.topic()

            if topic == "uzisplitted":
                uzi_splitted_event = pb.UziSplitted()
                uzi_splitted_event.ParseFromString(msg.value())

                print("UZI ID: ", uzi_splitted_event.uzi_id)

                self.uzi.segmentClassificateSave(
                    uzi_splitted_event.uzi_id, uzi_splitted_event.pages_id
                )
            elif topic == "cytologysplitted" and self.cytology:
                cytology_splitted_event = pb.CytologySplitted()
                cytology_splitted_event.ParseFromString(msg.value())

                print("CYTOLOGY ID: ", cytology_splitted_event.cytology_id)
                print("ORIGINAL IMAGE ID: ", cytology_splitted_event.original_image_id)

                self.cytology.processCytologyImage(
                    cytology_splitted_event.cytology_id,
                    cytology_splitted_event.original_image_id
                )

            consumer.commit(msg)

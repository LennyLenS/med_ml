from confluent_kafka import Producer
from ml_service.config.default import get_settings
from ml_service.internal.s3.s3 import S3
from ml_service.internal.ml_model.kt.Segmentation_YOLO import SegmentationModel
from ml_service.internal.ml_model.kt.Classification_ResNet import ClassificationModel
import ml_service.internal.events.kafka_pb2 as pb_event
import uuid
import cv2
import tempfile
import os
import numpy as np

settings = get_settings()

class ktUseCase:
    def __init__(self, segment_model: SegmentationModel, classificate_model: ClassificationModel, store: S3):
        self.segment_model = segment_model
        self.classification_model = classificate_model
        self.store = store
        
    def segmentClassificateSave(self, kt_id):
        print("Going to S3...")
        data = self.store.load(kt_id + "/" + kt_id)

        # segment
        segment_data = self.segment_model.predict(data)
        # Временно убираем сохранение видео - в ответе будем возвращать прямоугольник с областью сегментации
        self._store_video_to_s3(segment_data.get('mask'), kt_id, 'mask', True)
        self._store_video_to_s3(segment_data.get('detection'), kt_id, 'detected_video')
        video_mask = self._video_mask_mult(segment_data.get('original'), segment_data.get('mask'))

        # classificate
        result = self.classification_model.predict(video_mask)

        self.reconstruct_video(video_mask, kt_id, 'video_mask_multiplied')

        msg_event = pb_event.KtProcessed(
            kt_id=kt_id, class_probabilities = result
        )
        content = msg_event.SerializeToString()

        producer_config = {
            "bootstrap.servers": settings.kafka_host + ":" + str(settings.kafka_port)
        }
        producer = Producer(producer_config)

        producer.produce("ktprocessed", content)
        producer.flush()


        print(type(result))
        print(result)

    def _store_video_to_s3(self, video_array, kt_id, name, is_mask = False):
        print("Saving to S3...")
        # Сохраняем numpy array как видео .mp4
        height, width = video_array.shape[1:3]
        fourcc = cv2.VideoWriter_fourcc(*'mp4v')
        fps = 15
        # Создаем временный файл
        temp_file = tempfile.NamedTemporaryFile(suffix=".mp4", delete=False)
        temp_path = temp_file.name
        temp_file.close()

        # Записываем видео
        out = cv2.VideoWriter(temp_path, fourcc, fps, (width, height))
        for frame in video_array:
            if is_mask:
                # Для масок конвертируем в BGR
                frame = cv2.cvtColor(frame, cv2.COLOR_GRAY2BGR)
            out.write(frame)
        out.release()

        # Загружаем временный файл в MinIO
        path = kt_id + '/' + name + ".mp4"
        try:
            with open(temp_path, "rb") as file_data:
                file_stat = os.stat(temp_path)
                self.store.store_as_is(
                    file_data,
                    path,
                    file_stat.st_size,
                    content_type="video/mp4"
                )
            print("Видео успешно загружено в MinIO")
        finally:
            os.unlink(temp_path)  # Удаляем временный файл

    def _video_mask_mult(self, video, mask):
        video = self._convert_to_grayscale(video)
        mask = self._convert_to_grayscale(mask)
        new_video = np.zeros(shape=video.shape)
        for frame in range(video.shape[0]):
            for x in range(video.shape[1]):
                for y in range(video.shape[2]):
                    new_video[frame][x][y] = video[frame][x][y] * mask[frame][x][y]
        return new_video

    def _convert_to_grayscale(self, video_array, target_shape=(54, 224, 224, 1)):
        """
        Конвертирует видео массив в оттенки серого с нормализацией

        Args:
            video_array: Входной массив видео (N, H, W, 3) или (N, H, W)
            target_shape: Целевая форма выходного массива (frames, height, width, channels)

        Returns:
            Нормализованный массив в оттенках серого (N, H, W, 1) в диапазоне [0, 1]
        """
        # Если видео уже в оттенках серого (N, H, W)
        if video_array.ndim == 3:
            gray_frames = video_array
        # Если цветное видео (N, H, W, 3)
        elif video_array.ndim == 4 and video_array.shape[-1] == 3:
            gray_frames = np.array([cv2.cvtColor(frame, cv2.COLOR_RGB2GRAY) for frame in video_array])
        else:
            raise ValueError("Неподдерживаемая размерность входного массива")

        # Нормализация [0, 1]
        gray_frames = gray_frames.astype(np.float32) / 255.0

        # Выборка или дополнение кадров до target_shape[0]
        n_frames = gray_frames.shape[0]
        if n_frames < target_shape[0]:
            # Дополнение нулями
            pad = np.zeros((target_shape[0] - n_frames, *gray_frames.shape[1:]))
            gray_frames = np.vstack([gray_frames, pad])
        elif n_frames > target_shape[0]:
            # Равномерная выборка кадров
            indices = np.linspace(0, n_frames-1, target_shape[0], dtype=int)
            gray_frames = gray_frames[indices]

        # Изменение размера если нужно
        if gray_frames.shape[1:3] != target_shape[1:3]:
            gray_frames = np.array([cv2.resize(frame, (target_shape[2], target_shape[1]))
                                for frame in gray_frames])

        # Добавление оси канала если нужно
        if gray_frames.ndim == 3:
            gray_frames = gray_frames[..., np.newaxis]

        return gray_frames

    def reconstruct_video(self, video, kt_id, name, fps=30):
        # Определяем параметры видео
        height, width = video[0].shape[:2]

        # Создаем временный файл
        temp_file = tempfile.NamedTemporaryFile(suffix=".mp4", delete=False)
        temp_path = temp_file.name
        temp_file.close()

        # Инициализируем VideoWriter
        fourcc = cv2.VideoWriter_fourcc(*'mp4v')
        out = cv2.VideoWriter(temp_path, fourcc, fps, (width, height))

        # Обрабатываем каждый кадр
        for frame in video:
            # Конвертируем обратно в BGR, если нужно
            # print(frame.dtype)
            frame = (frame * 255).astype(np.uint8)
            # print(frame.dtype)
            if frame.ndim == 2 or (frame.ndim == 3 and frame.shape[2] == 1):
                frame = cv2.cvtColor(frame, cv2.COLOR_GRAY2BGR)

            # Записываем кадр
            out.write(frame)

        out.release()
        
        # Загружаем временный файл в MinIO
        path = kt_id + '/' + name + ".mp4"
        try:
            with open(temp_path, "rb") as file_data:
                file_stat = os.stat(temp_path)
                self.store.store_as_is(
                    file_data,
                    path,
                    file_stat.st_size,
                    content_type="video/mp4"
                )
            print("Видео успешно загружено в MinIO")
        finally:
            os.unlink(temp_path)  # Удаляем временный файл

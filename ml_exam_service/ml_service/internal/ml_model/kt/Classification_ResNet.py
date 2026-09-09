import torch
import cv2
import numpy as np
from torchvision import transforms
from torchvision.models import resnet50
import torch.nn as nn
import tempfile

class TemporalResNet(nn.Module):
    def __init__(self, num_classes):
        super().__init__()
        self.base_model = resnet50(pretrained=True)

        # Заменяем первый сверточный слой для входных данных с 1 каналом
        self.base_model.conv1 = nn.Conv2d(
            1, 64, kernel_size=7, stride=2, padding=3, bias=False
        )

        # Определяем RNN слой
        self.rnn = nn.GRU(
            input_size=2048,  # ResNet50 output features after layer4
            hidden_size=512,
            batch_first=True
        )

        # Финальный полносвязный слой
        self.fc = nn.Linear(512, num_classes)

    def forward(self, x):
        batch, timesteps, C, H, W = x.shape

        # Объединяем батч и временные шаги
        x = x.view(batch * timesteps, C, H, W)

        # Пропускаем через слои ResNet
        x = self.base_model.conv1(x)
        x = self.base_model.bn1(x)
        x = self.base_model.relu(x)
        x = self.base_model.maxpool(x)

        x = self.base_model.layer1(x)
        x = self.base_model.layer2(x)
        x = self.base_model.layer3(x)
        x = self.base_model.layer4(x)

        # Применяем адаптивный пулинг и сжимаем размерности
        x = nn.functional.adaptive_avg_pool2d(x, 1)
        x = x.squeeze(-1).squeeze(-1)

        # Восстанавливаем временную размерность
        x = x.view(batch, timesteps, -1)

        # Пропускаем через RNN
        x, _ = self.rnn(x)

        # Берем последний временной шаг и классифицируем
        x = self.fc(x[:, -1, :])
        return x

class ClassificationModel:
    def __init__(self) -> None:
        pass

    def load_model(self, weights_path, num_classes):
        model = TemporalResNet(num_classes)
        model.load_state_dict(torch.load(weights_path, map_location=torch.device('cpu')))
        model.eval()
        return model
    

    def process_video(self, video_data, img_size=224, frame_skip=3):
        transform = transforms.Compose([
            transforms.ToPILImage(),
            # transforms.Grayscale(num_output_channels=1),
            transforms.Resize((img_size, img_size)),
            transforms.ToTensor(),
            transforms.Normalize([0.5], [0.5])
        ])

        frames = []
        frame_count = 0

        # Создаем временный файл
        with tempfile.NamedTemporaryFile(suffix='.mp4') as temp_file:
            temp_file.write(video_data)
            temp_file.flush()  # Убедимся, что данные записаны на диск
            
            # Открываем видео с помощью OpenCV
            cap = cv2.VideoCapture(temp_file.name)

            while cap.isOpened():
                ret, frame = cap.read()
                if not ret:
                    break
                if frame_count % frame_skip == 0:
                    frame = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
                    frame = transform(frame).unsqueeze(0)  # [1, C, H, W]
                    frames.append(frame)
                frame_count += 1
            cap.release()

        frames = torch.cat(frames, dim=0)

        return frames.unsqueeze(0)  # [1, timesteps, C, H, W]
    
    def predict(self, video_mask):
        classes = ["Доброкачественный КТ фенотип", "Неопределенный КТ фенотип", "Злокачественный КТ фенотип"]
        model = self.load_model("./ml_service/internal/ml_model/models/kt/RESNET50AUG300.pth", num_classes=3)

        video_mask_tensor = torch.tensor(video_mask, dtype=torch.float32)
        video_mask_tensor = video_mask_tensor.unsqueeze(0)
        video_mask_tensor = video_mask_tensor.permute(0, 1, 4, 2, 3)

        input_data = video_mask_tensor
        with torch.no_grad():
            output = model(input_data)
            outputs_probs = torch.softmax(output, dim=1)

        print(f"Predicted class: {outputs_probs}")
        print(classes[np.argmax(outputs_probs)])
        
        probs = outputs_probs.squeeze()
        print(probs)

        result_dict = {k: v.item() for k, v in zip(classes, probs)}
        return result_dict
    

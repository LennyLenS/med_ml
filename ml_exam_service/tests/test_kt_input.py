import unittest
from unittest.mock import patch
import numpy as np
import torch
from ml_service.internal.ml_model.kt.Classification_ResNet import ClassificationModel


class KtInputTest(unittest.TestCase):
    def test_grayscale_inputs_match_temporal_model(self):
        for shape in ((54, 8, 8), (54, 8, 8, 1)):
            with self.subTest(shape=shape):
                received = []

                def model(value):
                    received.append(tuple(value.shape))
                    return torch.zeros((1, 3))

                with patch.object(ClassificationModel, 'load_model', return_value=model):
                    result = ClassificationModel().predict(np.zeros(shape, dtype=np.float32))
                self.assertEqual(received, [(1, 54, 1, 8, 8)])
                self.assertEqual(len(result), 3)
                self.assertAlmostEqual(sum(result.values()), 1, places=6)

    def test_color_input_is_rejected(self):
        with patch.object(ClassificationModel, 'load_model'):
            with self.assertRaisesRegex(ValueError, 'Expected grayscale frames'):
                ClassificationModel().predict(np.zeros((54, 8, 8, 3), dtype=np.float32))


if __name__ == '__main__':
    unittest.main()

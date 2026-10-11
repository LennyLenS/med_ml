import unittest

from ml_service.internal.events import kafka_pb2


class KtEventContractTest(unittest.TestCase):
    def test_result_uses_backend_wire_fields(self):
        fields = kafka_pb2.KtProcessed.DESCRIPTOR.fields_by_name
        self.assertEqual(fields["kt_id"].number, 100)
        self.assertEqual(fields["class_probabilities"].number, 200)
        message = kafka_pb2.KtProcessed(
            kt_id="study", class_probabilities={"normal": 1.0}
        )
        restored = kafka_pb2.KtProcessed.FromString(message.SerializeToString())
        self.assertEqual(restored.kt_id, "study")
        self.assertEqual(dict(restored.class_probabilities), {"normal": 1.0})

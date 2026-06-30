INSERT INTO models (name, task_type, description)
VALUES ('mobilenet', 'classification', 'MobileNetV2 image classifier — ImageNet 1000 classes');

INSERT INTO model_versions (model_id, version, file_path, input_shape, output_type)
VALUES (
    (SELECT id FROM models WHERE name = 'mobilenet'),
    'v1',
    'mobilenet/v1/model.onnx',
    '[1, 3, 224, 224]',
    'float32[1,1000]'
);
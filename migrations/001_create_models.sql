CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE models(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(255) NOT NULL UNIQUE,
    task_type   VARCHAR(100) NOT NULL,
    description TEST,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()    
);

CREATE TABLE model_version (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    model_id    UUID NOT NULL REFERENCES models(id) ON DELETE CASECADE,
    version     VARCHAR(50) NOT NULL,
    file_path   TEXT NOT NULL,
    input_shape JSONB NOT NULL,
    output_type VARCHAR(100) NOT NULL,
    is_actuve   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE(model_id,version)
)

CREATE TABLE inference_logs (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    model_name   VARCHAR(255) NOT NULL,
    version      VARCHAR(50)  NOT NULL,
    latency_ms   FLOAT        NOT NULL,
    status       VARCHAR(50)  NOT NULL, 
    error_msg    TEXT,
    result       JSONB,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE INDEX idx_inference_logs_model ON inference_logs(model_name, version);
CREATE INDEX idx_inference_logs_created ON inference_logs(created_at DESC);


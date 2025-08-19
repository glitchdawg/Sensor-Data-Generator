CREATE TABLE IF NOT EXISTS sensor_readings (
    id SERIAL PRIMARY KEY,
    id1 VARCHAR(10) NOT NULL,
    id2 INT NOT NULL,
    sensor_type VARCHAR(50) NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    ts TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_sensor_id1_id2_ts ON sensor_readings (id1, id2, ts);

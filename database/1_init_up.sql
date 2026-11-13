CREATE TABLE IF NOT EXISTS coins 
(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    request_id VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    price INT NOT NULL,
    currency VARCHAR(100) NOT NULL,
    source VARCHAR(100) NOT NULL,
    request_time DATE NOT NULL
);

CREATE INDEX idx_coins_name_request_time ON coins(name, request_time DESC);

-- CREATE INDEX idx_coins_request_id ON coins(request_id);

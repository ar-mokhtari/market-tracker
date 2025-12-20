CREATE DATABASE IF NOT EXISTS market_tracker CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE market_tracker;

CREATE TABLE IF NOT EXISTS prices (
    id INT AUTO_INCREMENT PRIMARY KEY,
    date VARCHAR(20),
    time VARCHAR(20),
    time_unix BIGINT,
    symbol VARCHAR(50) NOT NULL,
    name_en VARCHAR(100),
    name_fa VARCHAR(100),
    price VARCHAR(50),
    change_value VARCHAR(50),
    change_percent DECIMAL(10, 2),
    unit VARCHAR(20),
    type VARCHAR(20) NOT NULL,
    market_cap BIGINT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_symbol (symbol),
    INDEX idx_type (type),
    INDEX idx_created_at (created_at),
    UNIQUE KEY unique_symbol_type (symbol, type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

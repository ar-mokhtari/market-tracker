CREATE TABLE
IF NOT EXISTS prices
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    date VARCHAR
(20),
    time VARCHAR
(20),
    symbol VARCHAR
(50) NOT NULL,
    name_en VARCHAR
(100),
    name_fa VARCHAR
(100),
    price VARCHAR
(50),
    change_value VARCHAR
(50),
    change_percent DECIMAL
(10, 2),
    unit VARCHAR
(20),
    type VARCHAR
(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON
UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_symbol_type (symbol, type
)
);

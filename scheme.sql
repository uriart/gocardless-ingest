CREATE TABLE transactions (
    id VARCHAR(32) PRIMARY KEY,
    booking_date DATE NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    currency CHAR(3) NOT NULL,
    creditor_name VARCHAR(255),
    purpose_code VARCHAR(50),
    description TEXT,
    balance_after DECIMAL(15, 2) NOT NULL
);
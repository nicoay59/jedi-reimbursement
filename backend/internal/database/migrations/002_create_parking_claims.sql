CREATE TABLE IF NOT EXISTS parking_claims (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    employee_id BIGINT UNSIGNED NOT NULL,
    parking_date DATE NOT NULL,
    parking_location VARCHAR(200) NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    description TEXT NULL,
    receipt_path VARCHAR(255) NULL,
    status ENUM('PENDING', 'APPROVED', 'REJECTED')
        NOT NULL DEFAULT 'PENDING',
    admin_note TEXT NULL,
    reviewed_by BIGINT UNSIGNED NULL,
    reviewed_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT chk_parking_amount CHECK (amount > 0),
    CONSTRAINT fk_parking_employee
        FOREIGN KEY (employee_id) REFERENCES users(id),
    CONSTRAINT fk_parking_reviewer
        FOREIGN KEY (reviewed_by) REFERENCES users(id),
    INDEX idx_parking_employee (employee_id),
    INDEX idx_parking_status (status),
    INDEX idx_parking_date (parking_date)
)

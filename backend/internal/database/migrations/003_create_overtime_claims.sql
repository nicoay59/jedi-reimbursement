CREATE TABLE IF NOT EXISTS overtime_claims (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    employee_id BIGINT UNSIGNED NOT NULL,
    overtime_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    duration_hours DECIMAL(5,2) NOT NULL,
    work_description TEXT NOT NULL,
    status ENUM('PENDING', 'APPROVED', 'REJECTED')
        NOT NULL DEFAULT 'PENDING',
    admin_note TEXT NULL,
    reviewed_by BIGINT UNSIGNED NULL,
    reviewed_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT chk_overtime_duration CHECK (duration_hours > 0),
    CONSTRAINT fk_overtime_employee
        FOREIGN KEY (employee_id) REFERENCES users(id),
    CONSTRAINT fk_overtime_reviewer
        FOREIGN KEY (reviewed_by) REFERENCES users(id),
    INDEX idx_overtime_employee (employee_id),
    INDEX idx_overtime_status (status),
    INDEX idx_overtime_date (overtime_date)
)

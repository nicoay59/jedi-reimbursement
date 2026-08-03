CREATE TABLE IF NOT EXISTS claim_status_histories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    claim_type ENUM('PARKING', 'OVERTIME') NOT NULL,
    claim_id BIGINT UNSIGNED NOT NULL,
    previous_status ENUM('PENDING', 'APPROVED', 'REJECTED') NULL,
    new_status ENUM('PENDING', 'APPROVED', 'REJECTED') NOT NULL,
    note TEXT NULL,
    updated_by BIGINT UNSIGNED NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_history_updated_by
        FOREIGN KEY (updated_by) REFERENCES users(id),
    INDEX idx_history_claim (claim_type, claim_id),
    INDEX idx_history_updated_by (updated_by),
    INDEX idx_history_created_at (created_at)
)

ALTER TABLE parking_claims
    ADD COLUMN receipt_original_name VARCHAR(255) NULL AFTER receipt_path,
    ADD COLUMN receipt_mime_type VARCHAR(100) NULL AFTER receipt_original_name,
    ADD COLUMN receipt_size BIGINT UNSIGNED NULL AFTER receipt_mime_type

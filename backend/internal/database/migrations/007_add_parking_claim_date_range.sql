ALTER TABLE parking_claims
    ADD COLUMN parking_end_date DATE NULL AFTER parking_date,
    ADD INDEX idx_parking_employee_month (
        employee_id,
        parking_date,
        status
    )

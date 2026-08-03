INSERT INTO claim_status_histories (
    claim_type,
    claim_id,
    previous_status,
    new_status,
    note,
    updated_by,
    created_at
)
SELECT
    source.claim_type,
    source.claim_id,
    NULL,
    source.current_status,
    'Pengajuan dibuat',
    source.employee_id,
    source.created_at
FROM (
    SELECT
        'PARKING' AS claim_type,
        parking.id AS claim_id,
        parking.status AS current_status,
        parking.employee_id,
        parking.created_at
    FROM parking_claims AS parking

    UNION ALL

    SELECT
        'OVERTIME' AS claim_type,
        overtime.id AS claim_id,
        overtime.status AS current_status,
        overtime.employee_id,
        overtime.created_at
    FROM overtime_claims AS overtime
) AS source
WHERE NOT EXISTS (
    SELECT 1
    FROM claim_status_histories AS history
    WHERE history.claim_type = source.claim_type
      AND history.claim_id = source.claim_id
)

-- name: ListMaintenanceOutstandingFlats :many
SELECT f.id, f.flat_number, f.block, f.flat_type, f.area_sqft_hundredths,
  COALESCE(sum(v.outstanding_amount_paise), 0)::bigint AS total_outstanding_paise,
  count(*)::bigint AS unpaid_bill_count,
  min(b.billing_month) AS oldest_unpaid_month
FROM flats f
JOIN maintenance_bills b ON b.flat_id = f.id AND b.society_id = f.society_id
JOIN maintenance_bill_balances v ON v.bill_id = b.id AND v.outstanding_amount_paise > 0
WHERE f.society_id = sqlc.arg('society_id')
  AND f.is_active
  AND (sqlc.arg('before_id')::bigint = 0 OR f.id < sqlc.arg('before_id'))
  AND (sqlc.narg('block')::text IS NULL OR COALESCE(f.block, '') ILIKE '%' || sqlc.narg('block')::text || '%')
  AND (sqlc.narg('flat_number')::text IS NULL OR f.flat_number ILIKE '%' || sqlc.narg('flat_number')::text || '%')
  AND (
    sqlc.arg('search')::text = ''
    OR f.flat_number ILIKE '%' || sqlc.arg('search')::text || '%'
    OR COALESCE(f.block, '') ILIKE '%' || sqlc.arg('search')::text || '%'
  )
GROUP BY f.id, f.flat_number, f.block, f.flat_type, f.area_sqft_hundredths
ORDER BY f.id DESC
LIMIT sqlc.arg('limit');

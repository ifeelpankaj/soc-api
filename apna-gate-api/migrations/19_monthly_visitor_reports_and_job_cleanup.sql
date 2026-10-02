-- +migrate Up

CREATE TABLE monthly_visitor_report_deliveries (
    id BIGSERIAL PRIMARY KEY,
    society_id BIGINT NOT NULL REFERENCES societies(id) ON DELETE CASCADE,
    report_month DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    attempt_count INTEGER NOT NULL DEFAULT 0,
    processing_until TIMESTAMPTZ,
    recipients TEXT[] NOT NULL DEFAULT '{}',
    provider_message_id TEXT,
    last_error TEXT,
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_monthly_visitor_report_delivery
        UNIQUE (society_id, report_month),
    CONSTRAINT chk_monthly_visitor_report_month_first_day
        CHECK (report_month = date_trunc('month', report_month::timestamp)::date),
    CONSTRAINT chk_monthly_visitor_report_status
        CHECK (status IN ('pending', 'processing', 'failed', 'sent')),
    CONSTRAINT chk_monthly_visitor_report_attempt_count
        CHECK (attempt_count >= 0),
    CONSTRAINT chk_monthly_visitor_report_sent_state
        CHECK (sent_at IS NULL OR status = 'sent')
);

CREATE INDEX idx_monthly_visitor_reports_pending
    ON monthly_visitor_report_deliveries (report_month, society_id)
    WHERE sent_at IS NULL;

CREATE INDEX idx_monthly_visitor_reports_processing
    ON monthly_visitor_report_deliveries (processing_until)
    WHERE sent_at IS NULL AND processing_until IS NOT NULL;

CREATE INDEX idx_monthly_visitor_reports_sent
    ON monthly_visitor_report_deliveries (society_id, report_month)
    WHERE sent_at IS NOT NULL;

CREATE INDEX idx_visitor_entries_cleanup
    ON visitor_entries (society_id, created_at, updated_at, id)
    WHERE status IN ('rejected', 'checked_out', 'cancelled', 'expired', 'auto_closed');

CREATE INDEX idx_visitor_invites_cleanup
    ON visitor_invites (updated_at, id)
    WHERE status IN ('used', 'expired', 'cancelled');

CREATE INDEX idx_flat_member_invites_cleanup
    ON flat_member_invites (updated_at, id)
    WHERE status IN ('accepted', 'expired', 'cancelled');

CREATE INDEX idx_notifications_cleanup
    ON notifications (created_at, id);

CREATE TRIGGER monthly_visitor_report_deliveries_updated_at
BEFORE UPDATE ON monthly_visitor_report_deliveries
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- +migrate Down

DROP TRIGGER IF EXISTS monthly_visitor_report_deliveries_updated_at
    ON monthly_visitor_report_deliveries;

DROP INDEX IF EXISTS idx_notifications_cleanup;
DROP INDEX IF EXISTS idx_flat_member_invites_cleanup;
DROP INDEX IF EXISTS idx_visitor_invites_cleanup;
DROP INDEX IF EXISTS idx_visitor_entries_cleanup;
DROP INDEX IF EXISTS idx_monthly_visitor_reports_sent;
DROP INDEX IF EXISTS idx_monthly_visitor_reports_processing;
DROP INDEX IF EXISTS idx_monthly_visitor_reports_pending;

DROP TABLE IF EXISTS monthly_visitor_report_deliveries;

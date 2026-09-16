CREATE TABLE audit_logs (
    id INTEGER NOT NULL PRIMARY KEY,
    user_id INTEGER,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INTEGER,
    message TEXT,
    created_at DATETIME NOT NULL
);

CREATE INDEX idx_audit_logs_created_at
ON audit_logs(created_at DESC);
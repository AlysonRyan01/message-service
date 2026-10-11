CREATE TABLE message_type (
    message_type_id SMALLINT NOT NULL,
    name            TEXT NOT NULL UNIQUE,
    CONSTRAINT pk_message_type PRIMARY KEY (message_type_id)
);

CREATE TABLE message_log_status (
    message_log_status_id SMALLINT NOT NULL,
    name                  TEXT NOT NULL UNIQUE,
    CONSTRAINT pk_message_log_status PRIMARY KEY (message_log_status_id)
);

CREATE TABLE message_log (
    message_log_id        UUID NOT NULL DEFAULT gen_random_uuid(),
    message_log_status_id SMALLINT NOT NULL,
    attempts              INT NOT NULL DEFAULT 0,
    last_error            TEXT,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    sent_at               TIMESTAMPTZ,
    CONSTRAINT pk_message_log PRIMARY KEY (message_log_id),
    CONSTRAINT fk_message_log_status
        FOREIGN KEY (message_log_status_id)
        REFERENCES message_log_status (message_log_status_id)
);

CREATE TABLE message (
    message_id       UUID NOT NULL DEFAULT gen_random_uuid(),
    message_log_id   UUID NOT NULL UNIQUE,
    message_type_id  SMALLINT NOT NULL,
    service_order_id INT,
    phone            TEXT NOT NULL,
    payload          TEXT NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT pk_message PRIMARY KEY (message_id),
    CONSTRAINT fk_message_message_type
        FOREIGN KEY (message_type_id) REFERENCES message_type (message_type_id),
    CONSTRAINT fk_message_message_log
        FOREIGN KEY (message_log_id) REFERENCES message_log (message_log_id)
);

CREATE INDEX idx_message_service_order_id ON message (service_order_id);
CREATE INDEX idx_message_log_status_id ON message_log (message_log_status_id);
CREATE INDEX idx_message_log_pending
    ON message_log (created_at)
    WHERE message_log_status_id = 1;

INSERT INTO message_type (message_type_id, name) VALUES
    (1,  'welcome'),
    (2,  'estimate'),
    (3,  'estimate_link'),
    (4,  'estimate_approved'),
    (5,  'estimate_rejected'),
    (6,  'device_ready'),
    (7,  'rejected_ready'),
    (8,  'guarantee'),
    (9,  'product_sale'),
    (10, 'owner_estimate_decision'),
    (11, 'no_repair');

INSERT INTO message_log_status (message_log_status_id, name) VALUES
    (1, 'pending'),
    (2, 'processing'),
    (3, 'sent'),
    (4, 'failed');
-- +migrate Up
CREATE TABLE gamer_profile
(
    first_name      VARCHAR(50)       NOT NULL,
    last_name       VARCHAR(50)       NOT NULL,
    student_number  VARCHAR(8)        NOT NULL PRIMARY KEY,
    membership_tier INTEGER DEFAULT 0 NOT NULL,
    banned          BOOLEAN,
    notes           VARCHAR(250),
    created_at      DATE,
    membership_expiry_date DATE
);

CREATE TABLE gamer_activity
(
    student_number VARCHAR(8) NOT NULL CONSTRAINT student_number REFERENCES gamer_profile,
    pc_number      INTEGER,
    game           VARCHAR(250),
    started_at     TIMESTAMP,
    ended_at       TIMESTAMP,
    exec_name      VARCHAR(250)
);

-- +migrate Down
DROP TABLE gamer_activity;
DROP TABLE gamer_profile;

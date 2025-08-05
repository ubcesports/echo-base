-- +migrate Up
ALTER TABLE gamer_activity
    ADD id UUID PRIMARY KEY DEFAULT gen_random_uuid();

CREATE INDEX gamer_activity_started_at_idx ON gamer_activity(started_at);
CREATE INDEX gamer_activity_ended_at_idx ON gamer_activity(ended_at);

ALTER TABLE gamer_activity DROP CONSTRAINT student_number;

ALTER TABLE gamer_profile ADD id UUID UNIQUE DEFAULT gen_random_uuid();
ALTER TABLE gamer_profile DROP CONSTRAINT gamer_profile_pkey;
ALTER TABLE gamer_profile 
    ADD CONSTRAINT gamer_profile_student_number_key UNIQUE (student_number);
ALTER TABLE gamer_profile ADD CONSTRAINT gamer_profile_pkey PRIMARY KEY (id); 
ALTER TABLE gamer_profile DROP CONSTRAINT gamer_profile_id_key;

ALTER TABLE gamer_activity 
    ADD CONSTRAINT gamer_activity_student_number_fkey 
        FOREIGN KEY (student_number)
        REFERENCES gamer_profile (student_number);

-- +migrate Down
ALTER TABLE gamer_activity DROP COLUMN id; 

DROP INDEX gamer_activity_started_at_idx;
DROP INDEX gamer_activity_ended_at_idx;

ALTER TABLE gamer_activity DROP CONSTRAINT gamer_activity_student_number_fkey;

ALTER TABLE gamer_profile DROP COLUMN id;
ALTER TABLE gamer_profile DROP CONSTRAINT gamer_profile_student_number_key;
ALTER TABLE gamer_profile 
    ADD CONSTRAINT gamer_profile_pkey PRIMARY KEY (student_number); 

ALTER TABLE gamer_activity
    ADD CONSTRAINT student_number
        FOREIGN KEY (student_number)
        REFERENCES gamer_profile (student_number);

CREATE TABLE IF NOT EXISTS user_passphrases (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id bigint not null,
    passphrase varchar(192) not null,
    inserted_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at time zone 'utc'),

    CONSTRAINT fk_user_passphrases_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE cascade ON DELETE cascade
);

CREATE UNIQUE INDEX IF NOT EXISTS user_passphrases_passphrase_unique_index ON user_passphrases USING btree(passphrase);
CREATE INDEX IF NOT EXISTS user_passphrases_inserted_at ON user_passphrases USING btree(inserted_at);

CREATE TABLE IF NOT EXISTS user_passphrase_invalidations (
    passphrase_id bigint PRIMARY KEY,
    source_user_id bigint null,
    inserted_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at time zone 'utc'),

    CONSTRAINT fk_user_passphrases_invalidations_passphrase_id FOREIGN KEY (passphrase_id)
        REFERENCES user_passphrases(id) ON UPDATE cascade ON DELETE cascade,
    CONSTRAINT fk_user_passphrases_invalidations_source_user_id FOREIGN KEY (source_user_id)
        REFERENCES users(id) ON UPDATE cascade ON DELETE SET NULL
);
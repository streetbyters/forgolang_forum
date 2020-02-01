CREATE TABLE IF NOT EXISTS languages (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name varchar(64) NOT NULL,
    code varchar(10) NOT NULL,
    date_format varchar(64),
    inserted_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at time zone 'utc'),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at time zone 'utc')
);

CREATE UNIQUE INDEX IF NOT EXISTS languages_code ON languages USING btree(code);

CREATE TABLE IF NOT EXISTS category_languages (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    category_id bigint not null,
    language_id bigint not null,
    source_user_id bigint null,
    title varchar(128) not null,
    description text null,
    slug varchar(200) not null,
    inserted_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at time zone 'utc'),

    CONSTRAINT fk_category_languages_category_id FOREIGN KEY (category_id)
        REFERENCES categories(id) ON UPDATE cascade ON DELETE cascade,
    CONSTRAINT fk_category_languages_language_id FOREIGN KEY (language_id)
        REFERENCES languages(id) ON UPDATE cascade ON DELETE cascade,
    CONSTRAINT fk_category_languages_source_user_id FOREIGN KEY (source_user_id)
        REFERENCES users(id) ON UPDATE cascade ON DELETE set null
);

CREATE INDEX IF NOT EXISTS category_languages_slug ON category_languages USING btree(lower(slug));

CREATE TABLE IF NOT EXISTS tag_languages (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    tag_id bigint not null,
    language_id bigint not null,
    source_user_id bigint null,
    name varchar(32) not null,
    inserted_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at time zone 'utc'),

    CONSTRAINT fk_tag_languages_tag_id FOREIGN KEY (tag_id)
        REFERENCES tags(id) ON UPDATE cascade ON DELETE cascade,
    CONSTRAINT fk_tag_languages_language_id FOREIGN KEY (language_id)
        REFERENCES languages(id) ON UPDATE cascade ON DELETE cascade,
    CONSTRAINT fk_tag_languages_source_user_id FOREIGN KEY (source_user_id)
        REFERENCES users(id) ON UPDATE cascade ON DELETE set null
);
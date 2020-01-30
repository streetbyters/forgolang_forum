CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    username varchar(64) not null,
    password_digest varchar(200) null,
    email    varchar(64),
    email_hidden boolean default true,
    bio text null,
    url varchar(200) null,
    is_active boolean default true,
    avatar varchar(200) null,
    inserted_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at time zone 'utc'),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at time zone 'utc')
);

CREATE UNIQUE INDEX IF NOT EXISTS users_username_unique_index ON users USING btree (username);
CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique_index ON users USING btree (email);

CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    code varchar(20) not null,
    name varchar(20) not null,
    description text null,
    inserted_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at time zone 'utc'),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at time zone 'utc')
);

CREATE UNIQUE INDEX IF NOT EXISTS roles_code_unique ON roles USING btree(code);

CREATE TABLE IF NOT EXISTS role_permissions (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    role_id bigint not null,
    controller varchar(200) not null,
    method varchar(64) not null,
    inserted_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at time zone 'utc'),

    CONSTRAINT fk_role_permissions_role_id FOREIGN KEY (role_id)
        REFERENCES roles(id) ON UPDATE cascade ON DELETE cascade
);

CREATE UNIQUE INDEX IF NOT EXISTS role_permissions_role_controller_method_unique ON
    role_permissions USING btree(role_id, controller, method);
CREATE INDEX IF NOT EXISTS role_permissions_controller_index ON role_permissions USING btree(controller);
CREATE INDEX IF NOT EXISTS role_permissions_method_index ON role_permissions USING btree(method);

CREATE TABLE IF NOT EXISTS user_role_assignments (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id bigint not null,
    role_id bigint not null,
    source_user_id bigint null,
    inserted_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at time zone 'utc'),

    CONSTRAINT fk_user_role_assignments_user_id FOREIGN KEY (user_id)
        REFERENCES users(id) ON UPDATE cascade ON DELETE cascade,
    CONSTRAINT fk_user_role_assignments_role_id FOREIGN KEY (role_id)
        REFERENCES roles(id) ON UPDATE cascade ON DELETE cascade,
    CONSTRAINT fk_user_role_assignments_source_user_id FOREIGN KEY (source_user_id)
        REFERENCES users(id) ON UPDATE cascade ON DELETE set null
);

CREATE TABLE IF NOT EXISTS user_role_assignment_invalidations (
    assignment_id bigint primary key,
    source_user_id bigint null,

    CONSTRAINT fk_user_role_assignment_invalidations_assignment_id FOREIGN KEY (assignment_id)
        REFERENCES user_role_assignments(id) ON UPDATE cascade ON DELETE cascade,
    CONSTRAINT fk_user_role_assignment_invalidations_source_user_id FOREIGN KEY (source_user_id)
        REFERENCES users(id) ON UPDATE cascade ON DELETE set null
);
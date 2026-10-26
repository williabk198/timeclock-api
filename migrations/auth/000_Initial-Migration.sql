CREATE SCHEMA auth;
COMMENT ON SCHEMA auth IS 'Holds user login credentials and account metadata. A user''s broader profile data (name, contact info, employment, etc...) is owned by the admin service.';

CREATE TABLE auth.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    person_id uuid NOT NULL,
    email character varying(254) NOT NULL,
    pass text NOT NULL
);
COMMENT ON TABLE auth.users IS 'Holds the login credentials (email/password) for a person';

-- The password for this seed user is "Password123!"
INSERT INTO auth.users VALUES ('b3b5f0b4-0000-4000-8000-000000000001', 'f43d8d6f-bd21-4fd7-8c7e-61ca1077b789', 'testy.mctesterson@example.com', '$2a$10$E9katTZaEiN7tAqXbi5dEuqv5eNiUjjagS.5xQzavRaTj4QmpUsRa');

ALTER TABLE ONLY auth.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

ALTER TABLE ONLY auth.users
    ADD CONSTRAINT users_email_key UNIQUE (email);

ALTER TABLE ONLY auth.users
    ADD CONSTRAINT users_person_id_fkey FOREIGN KEY (person_id) REFERENCES person.persons(id);

CREATE TABLE auth.user_metadata (
    user_id uuid NOT NULL,
    created_at timestamptz NOT NULL,
    failure_count smallint NOT NULL DEFAULT 0,
    last_login timestamptz NOT NULL,
    last_pw_reset timestamptz NOT NULL,
    status smallint NOT NULL
);
COMMENT ON TABLE auth.user_metadata IS 'Holds metadata about a user''s account, such as login history and account status. An unset last_login is represented by the zero-value timestamp, not NULL.';

INSERT INTO auth.user_metadata VALUES ('b3b5f0b4-0000-4000-8000-000000000001', now(), 0, 'epoch', now(), 1);

ALTER TABLE ONLY auth.user_metadata
    ADD CONSTRAINT user_metadata_pkey PRIMARY KEY (user_id);

ALTER TABLE ONLY auth.user_metadata
    ADD CONSTRAINT user_metadata_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth.users(id);

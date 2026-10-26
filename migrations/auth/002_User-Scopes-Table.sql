CREATE TABLE auth.user_scopes (
    user_id uuid NOT NULL,
    scope text NOT NULL,
    granted_at timestamptz NOT NULL DEFAULT now()
);
COMMENT ON TABLE auth.user_scopes IS 'Holds the scopes explicitly granted to a user. Every authenticated user implicitly also carries the "default" scope on their tokens; it is not stored here.';

ALTER TABLE ONLY auth.user_scopes
    ADD CONSTRAINT user_scopes_pkey PRIMARY KEY (user_id, scope);

ALTER TABLE ONLY auth.user_scopes
    ADD CONSTRAINT user_scopes_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth.users(id);

-- Bootstrap: grant the seed user from 000_Initial-Migration.sql admin access so there's
-- at least one usable admin account after running migrations.
INSERT INTO auth.user_scopes (user_id, scope) VALUES ('b3b5f0b4-0000-4000-8000-000000000001', 'user:admin');

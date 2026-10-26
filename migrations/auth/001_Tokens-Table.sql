CREATE TABLE auth.tokens (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    token_hash character varying(64) NOT NULL,
    scope text NOT NULL,
    kind smallint NOT NULL,
    issued_at timestamptz NOT NULL,
    expires_at timestamptz NOT NULL,
    revoked boolean NOT NULL DEFAULT false
);
COMMENT ON TABLE auth.tokens IS 'Holds the hash, and metadata, of every opaque token issued by the auth service. The raw token value is never stored; only the hex-encoded SHA-256 hash of it is kept.';

ALTER TABLE ONLY auth.tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (id);

ALTER TABLE ONLY auth.tokens
    ADD CONSTRAINT tokens_token_hash_key UNIQUE (token_hash);

ALTER TABLE ONLY auth.tokens
    ADD CONSTRAINT tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth.users(id);

CREATE INDEX tokens_user_id_idx ON auth.tokens (user_id);

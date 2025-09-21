CREATE TABLE person.addresses (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    person_id uuid NOT NULL,
    street1 character varying(64) NOT NULL,
    street2 character varying(64) DEFAULT ''::character varying NOT NULL,
    locality character varying(32) NOT NULL,
    region character varying(32) NOT NULL,
    postal_code character varying(16) NOT NULL,
    country character varying(32) NOT NULL,
    kind character varying(16) NOT NULL,
    "primary" boolean NOT NULL
);
COMMENT ON COLUMN person.addresses.locality IS 'AKA town/city...';
COMMENT ON COLUMN person.addresses.region IS 'AKA state/province/prefecture... ';

CREATE TABLE person.emails (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    person_id uuid NOT NULL,
    username character varying(64) NOT NULL,
    provider character varying(64) NOT NULL,
    "primary" boolean NOT NULL
);

CREATE TABLE person.phones (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    person_id uuid NOT NULL,
    country_code integer NOT NULL,
    phone_number character varying(16) NOT NULL,
    kind character varying(16) NOT NULL,
    "primary" boolean NOT NULL
);

INSERT INTO person.addresses VALUES ('11e4207e-eb48-4e94-99ef-f4c791dad15d', 'f43d8d6f-bd21-4fd7-8c7e-61ca1077b789', '123 Testing St', 'APT 3', 'Testersvile', 'Testonia', '12345-6789', 'Testaria', 'mailing', true);
INSERT INTO person.addresses VALUES ('dd6aff6f-1614-4b86-9786-bad014edb414', 'f43d8d6f-bd21-4fd7-8c7e-61ca1077b789', '987 Test Ln', '', 'Testervile', 'Testonia', '12345-6789', 'Testaria', 'mailing', false);
INSERT INTO person.addresses VALUES ('f13ec596-a0f6-43e7-b939-c37a4b9f0f8e', 'f43d8d6f-bd21-4fd7-8c7e-61ca1077b789', '456 Test Dr', 'APT 1', 'Testervile', 'Testonia', '12345-6789', 'Testaria', 'mailing', false);

INSERT INTO person.emails VALUES ('16ed0d4f-5210-4b32-8f1a-a2979a0bc57f', 'f43d8d6f-bd21-4fd7-8c7e-61ca1077b789', 'tester', 'test.com', true);

INSERT INTO person.phones VALUES ('8e4381e4-f849-4bc2-9e2a-6d13028ed19a', 'f43d8d6f-bd21-4fd7-8c7e-61ca1077b789', 1, '(555) 555-1234', 'home', false);
INSERT INTO person.phones VALUES ('06f20b0f-6d1a-4860-b1c2-af661a7270f0', 'f43d8d6f-bd21-4fd7-8c7e-61ca1077b789', 1, '(555) 555-9876', 'cell', true);

ALTER TABLE ONLY person.addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (id);

ALTER TABLE ONLY person.emails
    ADD CONSTRAINT emails_pkey PRIMARY KEY (id);

ALTER TABLE ONLY person.phones
    ADD CONSTRAINT phones_pkey PRIMARY KEY (id);

CREATE UNIQUE INDEX addresses_person_id_kind_primary_idx ON person.addresses USING btree (person_id, kind, "primary") WHERE "primary";

CREATE UNIQUE INDEX emails_person_id_primary_idx ON person.emails USING btree (person_id, "primary") WHERE "primary";

CREATE UNIQUE INDEX phones_person_id_kind_primary_idx ON person.phones USING btree (person_id, kind, "primary") WHERE "primary";

ALTER TABLE ONLY person.addresses
    ADD CONSTRAINT addresses_person_id_fkey FOREIGN KEY (person_id) REFERENCES person.persons(id);

ALTER TABLE ONLY person.emails
    ADD CONSTRAINT emails_person_id_fkey FOREIGN KEY (person_id) REFERENCES person.persons(id);

ALTER TABLE ONLY person.phones
    ADD CONSTRAINT phones_person_id_fkey FOREIGN KEY (person_id) REFERENCES person.persons(id);

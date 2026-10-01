CREATE SCHEMA person;
COMMENT ON SCHEMA person IS 'Holds all the data relating to a person like name, date of birth, contact info, mailing address, etc...';

CREATE TABLE person.persons (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    given_name character varying(32) NOT NULL,
    family_name character varying(32) NOT NULL,
    first_name character varying(8) NOT NULL,
    dob date NOT NULL,
    gender character varying(16) NOT NULL,
    pronouns character varying(16) NOT NULL
);
COMMENT ON TABLE person.persons IS 'Holds the basic information of a person like name, birthday, and preferred gender & pronouns';

INSERT INTO person.persons VALUES ('f43d8d6f-bd21-4fd7-8c7e-61ca1077b789', 'Testy', 'McTesterson', 'given', '1970-01-01', 'non-binary', 'they/them');

ALTER TABLE ONLY person.persons
    ADD CONSTRAINT persons_pkey PRIMARY KEY (id);

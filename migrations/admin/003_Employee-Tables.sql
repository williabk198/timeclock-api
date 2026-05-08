CREATE SCHEMA metadata;
CREATE SCHEMA public;

CREATE TABLE metadata.employees (
    eid uuid NOT NULL,
    pay text NOT NULL,
    hire_date date NOT NULL,
    start_date date NOT NULL,
    sick_time_hrs numeric(5,2) NOT NULL,
    time_off_hrs numeric(5,2) NOT NULL,
    exempt boolean NOT NULL,
    status smallint NOT NULL
);

CREATE TABLE public.employees (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    pid uuid NOT NULL,
    reports_to_eid uuid,
    title text NOT NULL
);

INSERT INTO metadata.employees VALUES ('2fb4286a-1036-49ce-a228-0f1e30ea850a', '250000 USD/year', '2017-05-08', '2017-06-05', 52.25, 112.75, true, 1);

INSERT INTO public.employees VALUES ('2fb4286a-1036-49ce-a228-0f1e30ea850a', 'f43d8d6f-bd21-4fd7-8c7e-61ca1077b789', NULL, 'Owner');

ALTER TABLE ONLY metadata.employees
    ADD CONSTRAINT employees_pkey PRIMARY KEY (eid);

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT employees_pkey PRIMARY KEY (id);

ALTER TABLE ONLY metadata.employees
    ADD CONSTRAINT employees_eid_fkey FOREIGN KEY (eid) REFERENCES public.employees(id);

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT employees_pid_fkey FOREIGN KEY (pid) REFERENCES person.persons(id);

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT employees_reports_to_eid_fkey FOREIGN KEY (reports_to_eid) REFERENCES public.employees(id) NOT VALID;

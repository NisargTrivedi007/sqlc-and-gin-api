
CREATE TABLE IF NOT EXISTS public.todos
(
    id integer NOT NULL DEFAULT nextval('todos_id_seq'::regclass),
    task text COLLATE pg_catalog."default",
    created_by bigint,
    created_date date DEFAULT CURRENT_DATE,
    updated_date date,
    CONSTRAINT todos_pkey PRIMARY KEY (id)
)

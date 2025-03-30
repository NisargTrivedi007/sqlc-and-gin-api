
CREATE TABLE IF NOT EXISTS public.todos
(
    id integer NOT NULL DEFAULT nextval('todos_id_seq'::regclass),
    task text COLLATE pg_catalog."default",
    created_by bigint,
    created_date date DEFAULT CURRENT_DATE,
    updated_date date,
    done boolean DEFAULT false,
    CONSTRAINT todos_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.users
(
    id integer NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    username text COLLATE pg_catalog."default" NOT NULL,
    email_id text COLLATE pg_catalog."default" NOT NULL,
    phone_no numeric,
    created_date date NOT NULL DEFAULT CURRENT_DATE,
    password text COLLATE pg_catalog."default" NOT NULL ,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
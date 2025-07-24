set role zulu_updater;

-- DROP TABLE public.zws_to_object_type;

CREATE TABLE public.zws_to_object_type (
	rec_id serial4 NOT NULL,
	inserted_ts timestamp DEFAULT (now() AT TIME ZONE 'utc'::text) NULL,
	zws_type int4 NULL,
	object_type varchar(32) NULL,
	only_calculated bool NULL
);
set role zulu_updater;

-- DROP TABLE public.object_state_configuration;

CREATE TABLE public.object_state_configuration (
	zws_type_id int4 NULL,
	collapsed_category_id int4 NULL,
	full_category_id int4 NULL,
	object_type varchar(50) NULL
);
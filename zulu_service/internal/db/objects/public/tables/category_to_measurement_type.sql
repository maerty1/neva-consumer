set role zulu_updater;

-- DROP TABLE public.category_to_measurement_type;

CREATE TABLE public.category_to_measurement_type (
	rn int4 NULL,
	category_id int4 NULL,
	measurement_types_id int4 NULL,
	cut_unique_name text NULL
);
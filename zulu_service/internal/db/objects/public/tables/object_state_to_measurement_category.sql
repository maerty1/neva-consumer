set role zulu_updater;

-- DROP TABLE public.object_state_to_measurement_category;

CREATE TABLE public.object_state_to_measurement_category (
	measurement_category_id int4 NULL,
	zws_type_id int4 NULL,
	object_type varchar(50) NULL,
	rn int4 NULL
);
set role zulu_updater;

-- DROP TABLE public.measurement_categories;

CREATE TABLE public.measurement_categories (
	id int4 NULL,
	"name" varchar(50) NULL,
	cut_type varchar(50) NULL,
	max_values int4 NULL,
	is_open int4 NULL,
	expanded_rows_qty int4 NULL,
	object_type varchar(64) NULL
);
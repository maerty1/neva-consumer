set role zulu_updater;

-- DROP TABLE public.measurement_types;

CREATE TABLE public.measurement_types (
	id int4 NULL,
	zulu_desc text NULL,
	lers_desc text NULL,
	front_desc text NULL,
	zulu_var varchar(50) NULL,
	zulu_un varchar(50) NULL,
	ler_var varchar(50) NULL,
	rest_var varchar(50) NULL,
	lers_un varchar(50) NULL,
	scada_var varchar(50) NULL,
	scada_un varchar(50) NULL,
	is_calculated bool NULL
);
set role zulu_updater;

-- DROP TABLE public.measurement_groups;

CREATE TABLE public.measurement_groups (
	id int4 NULL,
	group_name text NULL,
	group_front_desc text NULL,
	"in" int4 NULL,
	"out" int4 NULL,
	measurement_unit varchar(32) NULL
);
CREATE INDEX idx_measurement_groups_id ON public.measurement_groups USING btree (id);
CREATE INDEX idx_measurement_groups_in ON public.measurement_groups USING btree ("in");
CREATE INDEX idx_measurement_groups_out ON public.measurement_groups USING btree ("out");
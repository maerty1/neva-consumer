set role zulu_updater;

-- DROP TABLE public.category_to_group;

CREATE TABLE public.category_to_group (
	rn int4 NULL,
	category_id int4 NULL,
	group_id int4 NULL,
	cut_unique_name varchar(64) NULL
);
CREATE INDEX idx_category_to_group_category_id ON public.category_to_group USING btree (category_id);
CREATE INDEX idx_category_to_group_group_id ON public.category_to_group USING btree (group_id);
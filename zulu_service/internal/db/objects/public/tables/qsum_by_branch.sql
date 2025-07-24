set role zulu_updater;

-- DROP TABLE public.qsum_by_branch;

CREATE TABLE public.qsum_by_branch (
	istok_id int4 NULL,
	entrance_elem_id int4 NULL,
	entrance_pipe_id int4 NULL,
	elem_id int4 NULL,
	qsum numeric(10, 5) NULL,
	q_ot numeric(10, 5) NULL
);
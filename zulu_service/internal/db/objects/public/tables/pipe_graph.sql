set role zulu_updater;

-- DROP TABLE public.pipe_graph;

CREATE TABLE public.pipe_graph (
	parent_elem_id int4 NULL,
	child_elem_id int4 NULL,
	joint_point public.geometry NULL
);
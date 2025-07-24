--SET ROLE zulu_updater;
--
--create table public.qsum_by_branch
--						  (istok_id	INT4,
--                           entrance_elem_id INT4,
--                           entrance_pipe_id INT4,
--                           elem_id	INT4,
--                           qsum	numeric(10,5),
--                           q_ot numeric(10,5));

SET ROLE zulu_updater;

CREATE OR REPLACE FUNCTION public.get_qsum_by_object(_task_parameters json)
 RETURNS json
 LANGUAGE plpgsql
 SECURITY DEFINER
AS $$
DECLARE 
    _error_context    TEXT;
    _error_msg        TEXT;
BEGIN

TRUNCATE TABLE public.qsum_by_branch;

WITH RECURSIVE pipes_in_branch AS (
    SELECT 
        ep.elem_id entrance_elem_id, 
        zws_geometry as entrance_point,
        child_elem_id as entrance_pipe_id, 
        child_elem_id as pipe_id 
        FROM pipe_graph g
        JOIN (select elem_id, zws_geometry 
                            from zulu.objects_geometry_log
                            where elem_id in (select json_array_elements_text(_task_parameters->'entrance_elem_ids')::int4)) ep
        ON g.joint_point = zws_geometry
    UNION ALL
    SELECT 
        entrance_elem_id,
        entrance_point,
        entrance_pipe_id,
        child_elem_id as pipe_id
    FROM pipes_in_branch b
    JOIN pipe_graph g ON b.pipe_id = g.parent_elem_id
    WHERE joint_point not in (select zws_geometry 
                        from zulu.objects_geometry_log
                        where elem_id in (select json_array_elements_text(_task_parameters->'excluded_elem_ids')::int4))
),
pipes_data as (
    select 
        elem_id,
        ST_StartPoint(zws_geometry) start_point, 
        ST_EndPoint(zws_geometry) end_point
	from zulu.objects_geometry_log
	where zws_type = 6
),
points_in_branch as (
    SELECT distinct
    entrance_elem_id,
    entrance_point,
    entrance_pipe_id,
    d.start_point geo_point
    FROM pipes_in_branch p
    JOIN pipes_data d
    ON p.pipe_id = d.elem_id
    union 
    SELECT distinct
    entrance_elem_id,
	entrance_point,
    entrance_pipe_id,
    d.end_point geo_point
    FROM pipes_in_branch p
    JOIN pipes_data d
    ON p.pipe_id = d.elem_id
),
objects as (
    select distinct p.entrance_elem_id, p.entrance_pipe_id,
       g.elem_id
    from points_in_branch p
    join zulu.objects_geometry_log g
    on g.zws_geometry = p.geo_point
    and zws_type IN (8,3)
    and g.zws_geometry::text != entrance_point::text
)
insert into public.qsum_by_branch(istok_id,
                           entrance_elem_id,
                           entrance_pipe_id,
                           elem_id,
                           qsum)
select r2.val::int istok_id,
       entrance_elem_id,
       entrance_pipe_id,
       o.elem_id,
       r.val::numeric qsum
from objects o
left join zulu.object_records r
on o.elem_id = r.elem_id
and r."parameter" = 'Qsum'
join zulu.object_records r2
on o.elem_id = r2.elem_id
and r2."parameter" = 'Nist';

update public.qsum_by_branch
set istok_id = r.elem_id
from zulu.zulu.object_records r 
where r.elem_id in (select elem_id 
                   from zulu.zulu.objects_geometry_log i 
                   where zws_type = 1)
and r."parameter" in ('Nist')
and r.val::int = istok_id;

RETURN JSON_BUILD_OBJECT('task_status', 200);

EXCEPTION
    WHEN OTHERS THEN
        GET STACKED DIAGNOSTICS _error_msg = MESSAGE_TEXT,
                                _error_context = PG_EXCEPTION_CONTEXT;

        RETURN JSON_BUILD_OBJECT('task_status',    415,
                                'errors',         _error_msg||' '|| _error_context
                                );

END;
$$
;

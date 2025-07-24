WITH RECURSIVE pipes_in_branch AS (
    SELECT 
        ep.elem_id entrance_elem_id, 
        zws_geometry as entrance_point,
        child_elem_id as entrance_pipe_id, 
        child_elem_id as pipe_id 
        FROM pipe_graph g
        JOIN (select elem_id, zws_geometry 
                            from zulu.objects_geometry_log
                            where elem_id = $1::integer) ep
        ON g.joint_point = zws_geometry
    UNION ALL
    SELECT 
        entrance_elem_id,
        entrance_point,
        entrance_pipe_id,
        child_elem_id as pipe_id
    FROM pipes_in_branch b
    JOIN pipe_graph g ON b.pipe_id = g.parent_elem_id
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
)
select distinct 
    p.entrance_elem_id,
    g.elem_id,
    g.zws_type
from points_in_branch p
join zulu.objects_geometry_log g
    on g.zws_geometry = p.geo_point
    and zws_type != 6
    and g.zws_geometry::text != entrance_point::text
union all
select distinct 
    entrance_elem_id,
    pipe_id,
    6 zws_type
from pipes_in_branch
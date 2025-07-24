--select public.get_graph()

SET ROLE zulu_updater;

CREATE OR REPLACE FUNCTION public.get_graph()
 RETURNS json
 LANGUAGE plpgsql
 SECURITY DEFINER
AS $$
DECLARE
    _processed_elems    INT4;
    _unprocessed_elems 	INT4[];
    _graph_depth        INT4 = 0;
    _error_context    TEXT;
    _error_msg        TEXT;
BEGIN

--делаем тмп таблицу всех труб - начало и конец
CREATE TEMP TABLE tmp_src
(
elem_id         INT4,
start_point     GEOMETRY,
end_point       GEOMETRY
)
ON COMMIT DROP;

--наполняем
INSERT INTO tmp_src (elem_id, start_point, end_point)
SELECT elem_id, 
ST_StartPoint(zws_geometry) start_point, 
ST_EndPoint(zws_geometry) end_point 
FROM zulu.zulu.objects_geometry_log ogl 
WHERE zws_type = 6;

--делаем тмп таблицу графа
CREATE TEMP TABLE tmp_graph
(
parent_elem_id         INT4,
child_elem_id          INT4,
joint_point         GEOMETRY
)
ON COMMIT DROP;

--берем котельные и трубы котельных
--заносим в таблицу графа
insert into tmp_graph(child_elem_id, joint_point)
select s.elem_id, i.zws_geometry
from zulu.zulu.objects_geometry_log i 
join tmp_src s on s.start_point = i.zws_geometry
or s.end_point = i.zws_geometry
where zws_type = 1;

--цикл:
--джойним трубы из таблицы по началу И концу к началу И концу если их нет в таблице графа
--записывам результат в таблицу графа
--останавливаемся когда ничего не джойнится
LOOP
    insert into tmp_graph(parent_elem_id, child_elem_id)
    select child_elem_id parent_elem_id, s2.elem_id child_elem_id from tmp_graph
    join tmp_src s1
    on s1.elem_id = child_elem_id
    join tmp_src s2
    on (s1.start_point = s2.end_point
    or s1.start_point = s2.start_point
    or s1.end_point = s2.start_point
    or s1.end_point = s2.end_point)
    and s2.elem_id not in ((select child_elem_id from tmp_graph
    where child_elem_id is not null)
    union 
    (select parent_elem_id from tmp_graph
    where parent_elem_id is not null))
    where child_elem_id not in (select parent_elem_id from tmp_graph where parent_elem_id is not null)
    and child_elem_id is not null;

    IF NOT FOUND THEN
        EXIT;
    END IF; 

    select _graph_depth + 1
    INTO _graph_depth;

END LOOP;

update tmp_graph g
set joint_point = j.joint_point
from (
    select 
        child_elem_id, 
        parent_elem_id, 
        case 
            when ST_StartPoint(g.zws_geometry)=ST_StartPoint(g1.zws_geometry)
            or ST_StartPoint(g.zws_geometry)=ST_EndPoint(g1.zws_geometry) 
            then ST_StartPoint(g.zws_geometry)
            when ST_EndPoint(g.zws_geometry)=ST_StartPoint(g1.zws_geometry)
            or ST_EndPoint(g.zws_geometry)=ST_EndPoint(g1.zws_geometry)
            then ST_EndPoint(g.zws_geometry)
        end joint_point
    from tmp_graph
    join zulu.zulu.objects_geometry_log g
    on child_elem_id = g.elem_id
    join zulu.zulu.objects_geometry_log g1
    on child_elem_id = g1.elem_id
    ) j
where j.child_elem_id = g.child_elem_id
and j.parent_elem_id = g.parent_elem_id
and g.joint_point is null;

update tmp_graph t
set joint_point = g.zws_geometry
from zulu.zulu.objects_geometry_log g
join zulu.zulu.objects_geometry_log g1
    on (
        ST_EndPoint(g1.zws_geometry) = g.zws_geometry
        or ST_StartPoint(g1.zws_geometry) = g.zws_geometry
        )
    and g.zws_type = 1
where joint_point is null
and child_elem_id = g1.elem_id;

TRUNCATE TABLE public.pipe_graph;

INSERT INTO public.pipe_graph(parent_elem_id, child_elem_id, joint_point)
SELECT parent_elem_id, 
	   child_elem_id,
       joint_point
FROM tmp_graph;

select count(*)
INTO _processed_elems
FROM ((select distinct child_elem_id elem_id from tmp_graph
    where child_elem_id is not null)
    union 
    (select distinct parent_elem_id elem_id from tmp_graph
    where parent_elem_id is not null)) e;

with processed_elems as ((select distinct child_elem_id elem_id from tmp_graph
    where child_elem_id is not null)
    union 
    (select distinct parent_elem_id elem_id from tmp_graph
    where parent_elem_id is not null))
select array_agg(s.elem_id)::text
INTO _unprocessed_elems
from tmp_src s
left join processed_elems p on p.elem_id = s.elem_id
where p.elem_id is null;

RETURN JSON_BUILD_OBJECT('task_status', 200,
                         'processed_elems',  _processed_elems ,
                         'unprocessed_elem_ids', _unprocessed_elems,
                         'graph_depth', _graph_depth);

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



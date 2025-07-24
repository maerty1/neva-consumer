with border_objects as (
    select 
        i.elem_id, 
        ST_StartPoint(i.zws_geometry) start_point, 
        max(i2.elem_id) start_elem_id, 
        ST_EndPoint(i.zws_geometry) end_point, 
        max(i3.elem_id) end_elem_id 
    from zulu.zulu.objects_geometry_log i 
    left join zulu.zulu.objects_geometry_log i2
        on i2.zws_type != 6
        and ST_StartPoint(i.zws_geometry) = i2.zws_geometry
    left join zulu.zulu.objects_geometry_log i3
        on i3.zws_type != 6
        and ST_EndPoint(i.zws_geometry) = i3.zws_geometry
    where i.zws_type = 6
    group by i.elem_id, i.zws_geometry
)
select 
    b.elem_id,  
    start_point,
    r.val::numeric                          start_h_geo,
    end_point,
    r2.val::numeric                         end_h_geo,
    (r.val::numeric + r2.val::numeric)/2    avg_h_geo
from border_objects b
join zulu.zulu.object_records r
    on start_elem_id = r.elem_id
    and r."parameter" = 'H_geo'
join zulu.zulu.object_records r2
    on end_elem_id = r2.elem_id
    and r2."parameter" = 'H_geo';
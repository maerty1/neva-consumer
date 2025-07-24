--кадр 1 (Каждый блок - источник (котельные))
select r.elem_id                            block_id, 
       coalesce(em.title, r.elem_id::text)  block_name, 
       val::numeric                         qsum 
from zulu.zulu.object_records r 
left join zulu.zulu.elems_metadata em 
on em.elem_id = r.elem_id
where r.elem_id in(select elem_id 
                   from zulu.zulu.objects_geometry_log i 
                   where zws_type = 1)
and "parameter" in ('Qsum');

--кадр 2, переменная - elem_id источника (из кадра выше, он же block_id)
with branches as (
            select distinct entrance_elem_id branch_id from public.qsum_by_branch
where istok_id = $1
),
src as (
select 
            branch_id
 			block_id, 
	        entrance_pipe_id, 
	        q.elem_id,
	        qsum
        from public.qsum_by_branch q
        join branches 
        on entrance_elem_id = branch_id
        and istok_id = $1     
        where q.elem_id not in (select * from branches)
        )
select 
    block_id, 
    case 
        when block_id = $1
        then coalesce (b.branch_name, 'Остальные')
        else coalesce (b.branch_name, em.title, block_id::text)
    end         block_name, 
    SUM(qsum)   qsum
from src
left join zulu.zulu.elems_metadata em 
on em.elem_id = block_id
left join branch_names b
on b.entrance_elem_id = block_id
group by block_id, em.title, branch_name;

--кадр 3 (ЦТП/ветка), переменная - elem_id ЦТП или ветки (При нажатии на цтп из кадра выше)
select 
    q.elem_id   block_id, 
    em.title    block_name,       
    q.qsum
from public.qsum_by_branch q
join zulu.zulu.elems_metadata em 
on em.elem_id = q.elem_id
where q.entrance_elem_id = $1::integer
and qsum is not null;

--кадр 3 (Остальные), переменная - elem_id источника
with branches as (
            select distinct entrance_elem_id branch_id from public.qsum_by_branch
where istok_id = $1
and entrance_elem_id != $1
)
select q.elem_id 			            block_id, 
	   coalesce(em.title, q.elem_id::text)    block_name,
	   qsum
from public.qsum_by_branch q
left join branches 
on entrance_elem_id = branch_id
or q.elem_id = branch_id
left join zulu.zulu.elems_metadata em 
on em.elem_id = q.elem_id 
where branch_id is null
and istok_id = $1::integer;
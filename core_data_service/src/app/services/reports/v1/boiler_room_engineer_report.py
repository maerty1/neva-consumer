import asyncio

from aiocache import Cache, cached


@cached(ttl=60*10, key_builder=lambda f, self, temperature: f"boiler_room:{temperature}")
async def boiler_room_engineer_report(self, temperature: float):
    less_task = self.reports_repository.boiler_room_engineer_report_v2(temperature=temperature, direction=-1)
    current_task = self.reports_repository.boiler_room_engineer_report_v2(temperature=temperature, direction=0)
    more_task = self.reports_repository.boiler_room_engineer_report_v2(temperature=temperature, direction=1)

    less, current, more = await asyncio.gather(less_task, current_task, more_task)

    return {"less": less, "current": current, "more": more}

from fastapi import FastAPI

from app.builder import Application
from app.dblayer.connection import provide_pool
import logging
import asyncio
from datetime import datetime
from apscheduler.schedulers.asyncio import AsyncIOScheduler
from fastapi import FastAPI

from aiocache import caches


app = FastAPI(
    title="Core Data Service",
    version="0.1.0",
    docs_url="/core/api/docs",
    redoc_url="/core/api/redoc",
    openapi_url="/core/api/openapi.json",
)


@app.on_event("startup")
async def initialize_application_resources():
    pool = await provide_pool()
    application_builder = Application(app=app, pool=pool)
    final_application = application_builder.build()
    await final_application.startup()

    return final_application.app


@app.on_event("shutdown")
async def shutdown_event():
    # Здесь можно освободить асинхронные ресурсы
    print("Shutdown event")

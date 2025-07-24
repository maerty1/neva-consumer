from aiocache import caches

caches.set_config(
    {
        "default": {
            "cache": "aiocache.SimpleMemoryCache",
            "serializer": {"class": "aiocache.serializers.PickleSerializer"},
            "ttl": 60 * 60 * 12,  # 12 часов
        }
    }
)

cache = caches.get("default")

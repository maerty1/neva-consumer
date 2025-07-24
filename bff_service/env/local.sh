export BFF__HTTP_HOST='0.0.0.0'
export BFF__HTTP_PORT=8003
export JWT_SECRET='your_jwt_secret'
export GIN_MODE=debug

export USERS_HTTP_SERVICE_URL=http://host.docker.internal:8001
export CORE_DATA_SERVICE_HTTP_SERVICE_URL=http://host.docker.internal:8002
export ZULU_SERVICE_HTTP_SERVICE_URL=http://host.docker.internal:8004
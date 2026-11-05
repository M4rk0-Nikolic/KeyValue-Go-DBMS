## Key Value DataBase Managment System written in Golang
Written for linux and other POSIX or POSIX-compliant operating systems. Optimized for large amounts of write operations therefore poorer for read heavy applications.
All Docker configuration files (`Dockerfile`, `docker-compose.yml`, `.dockerignore`) can be found in the root directory. `configuration.yml` allows for changes to the behaviour of the system.
# Supported Operations:
- PUT(key,  value)
- GET(key) -> value
- DELETE(key)
- PREFIX_SCAN(prefix, pageNumber, pageSize) -> [page](key, value)
- RANGE_SCAN(range, pageNumber, pageSize) -> [page](key, value)
- PREFIX_ITERATE(prefix) -> iterator
- RANGE_ITERATE(range) -> iterato


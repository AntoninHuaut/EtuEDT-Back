# EtuEDT-Deno
docker-compose.yml
```YAML
version: "3.9"
services:
  etuedt_deno:
    container_name: etuedt_deno
    restart: always
    build: .
    ports:
      - "4000:8080"
    environment:
      - MYSQL_HOST=
      - MYSQL_PORT=3306
      - MYSQL_DATABASE=
      - MYSQL_USER=
      - MYSQL_PASSWORD=
```
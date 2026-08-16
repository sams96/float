# float

Trigger commands in containers with webhooks. I use it to trigger rclone from
paperless. I wrote some more about this setup at
[samsm.ch/paperless-cloud-backups/](https://samsm.ch/paperless-cloud-backups/)

```yaml
services:
  broker: ...
  db: ...
  webserver:
    ...
    volumes:
      - export:/usr/src/paperless/export
      ...

  rclone:
    container_name: rclone
    image: rclone/rclone:latest
    profiles: ["rclone"]
    volumes:
      - export:/data
      - ./rclone/config:/config/rclone
    command: "copy /data [remote:location]"

  float:
    container_name: float
    image: ghcr.io/sams96/float:latest
    ports:
      - "41232:41232"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    environment:
      FLOAT_CMD: "docker exec -d paperless-webserver document_exporter /usr/src/paperless/export && docker start rclone"

volumes:
  export:
```

# EtuEDT-Back

Type `docker compose up -d`

University settings are loaded from `config.json` at startup. ADE remains the
default provider. To use the timetable JSON API for a university, add an optional
`sourceOverride` object to that university:

```json
"sourceOverride": {
  "provider": "reparser",
  "baseUrl": "https://your-timetable-api.example",
  "cacheMinutes": 30
}
```

`baseUrl` must be an HTTPS origin without credentials, query parameters, or a path.
`cacheMinutes` defaults to 30 and accepts values from 5 to 720. Remove the override
to restore ADE. Restart the backend after editing configuration; when using the
published image, mount your configuration at `/app/config.json` or rebuild the image.

The override fetches events for the existing local group and room catalogue from
the provider's v3 endpoints. Public response schemas remain unchanged. Cache entries
are isolated by source and academic year, with the last successful result retained
if a refresh fails. Free-room search returns 503 while an override is active because
calendar coverage cannot establish room availability. Calendar subscription/export
continues to depend on ADE; no new ICS endpoint is provided.

Run `go test ./...` and `go vet ./...` to validate changes.

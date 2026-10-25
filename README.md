# Wii/DS NAS Server

OpenPak's Wii/DS NAS Server.

| | |
| --- | --- |
| Network | Wii, DS |
| Role | NAS authentication and native HTTP routing: the WFC HTTP front door for Wii and DS |
| Depends on | [`wfc-sake`](https://github.com/openpak/wfc-sake), [`wfc-race`](https://github.com/openpak/wfc-race), [`wfc-gamestats`](https://github.com/openpak/wfc-gamestats), [`wfc-api`](https://github.com/openpak/wfc-api), [`wfc-dls`](https://github.com/openpak/wfc-dls), shared WFC database (`WFC_DB_ADDRESS`) |
| Image | `ghcr.io/openpak/wfc-nas` |
| Based on | [WiiLink WFC](https://github.com/WiiLink24/wfc-server) (AGPL-3.0) |
| Licence | AGPL-3.0-only |

## What it serves

| Endpoint | What |
| --- | --- |
| HTTP front door on `WFC_NAS_PORT` (default 20120) | Native auth and HTTP routing for Wii and DS |
| SAKE and race hosts | Forwarded to `WFC_SAKE_URL` / `WFC_RACE_URL` |
| Gamestats and API routes | Forwarded to `WFC_GAMESTATS_URL` / `WFC_API_URL` |

## Build

```sh
docker build -t wfc-nas .
GOWORK=off go test -mod=vendor ./...
```

## Docs

- [Protocol](docs/protocol.md)
- [Evidence](docs/evidence.md)
- [Deploy](docs/deploy.md)
- [Notes](docs/notes.md)

## Licence

AGPL-3.0-only. See [LICENSE](LICENSE).

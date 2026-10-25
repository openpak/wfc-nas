# wfc-nas

Independent Wii/DS NAS authentication and native HTTP routing, extracted from nn-wfc.
Preserves the native auth codec, duplicate Host handling, connection-test pages,
profanity checks and operator-mounted payload assets. AGPL-3.0 (WiiLink WFC).

Serves the existing WFC HTTP front door on WFC_NAS_PORT (default 20120).
Forwards SAKE/race hosts to WFC_SAKE_URL / WFC_RACE_URL and gamestats/API routes to
WFC_GAMESTATS_URL / WFC_API_URL, with an authenticated internal request header. Each upstream
must be an HTTP origin without path/query/credentials. Host and native body/reply
are preserved; upstream failure returns 503.

Required: WFC_DB_USER, WFC_DB_PASSWORD, WFC_DB_ADDRESS, WFC_API_SECRET,
WFC_TOKEN_KEY_FILE (shared binary key >=32 bytes), WFC_INTERNAL_KEY (>=32 bytes),
WFC_SAKE_URL / WFC_RACE_URL and WFC_GAMESTATS_URL / WFC_API_URL. Optional WFC_DB_NAME=wwfc, WFC_LOG_LEVEL=4.
Mount the existing payload volume at /app/payload. The wfc-dls owner mounts DLC
at /app/dlc.
The shared database schema and account IDs remain unchanged.

Build/test independently: `GOWORK=off go test -mod=vendor ./...`,
`GOWORK=off go vet -mod=vendor ./...`, `docker build -t wfc-nas .`.
Shared code is vendored from sibling wfc-common. Image ghcr.io/openpak/wfc-nas.
Read [the full split guide](https://github.com/openpak/wfc-common/blob/main/docs/service-split.md) before deploying.

DLS1 now belongs to wfc-dls; configure WFC_DLS_URL. NAS no longer carries DLC data.

# Deploy: wfc-nas

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

## DLS

Set `WFC_DLS_URL` to the wfc-dls service.

## Link-first (prds#139)

`WFC_LINK_FIRST=1` refuses unlinked consoles at NAS login with their link number. It needs the
WFC database settings (`WFC_DB_*`, already in the shared env) and optionally `WFC_CORE_ADDRESS`/
`WFC_CORE_KEY` (heal a half-done approval). `WFC_NET_TAG_KEY` (else `WFC_INTERNAL_KEY`) must be
the same as wfc-api's.


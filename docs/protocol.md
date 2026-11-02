# Protocol: wfc-nas

Serves the existing WFC HTTP front door on WFC_NAS_PORT (default 20120).
Forwards SAKE/race hosts to WFC_SAKE_URL / WFC_RACE_URL and gamestats/API routes to
WFC_GAMESTATS_URL / WFC_API_URL, with an authenticated internal request header. Each upstream
must be an HTTP origin without path/query/credentials. Host and native body/reply
are preserved; upstream failure returns 503.

## Scope

NAS authentication and native HTTP routing, extracted from the former nn-wfc monolith.
Preserves the native auth codec, duplicate Host handling, connection-test pages, profanity
checks and operator-mounted payload assets. DLS1 belongs to wfc-dls; NAS carries no DLC data.

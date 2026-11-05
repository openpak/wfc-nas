# Changelog

## v0.3.1 - 2026-11-05

- Fix: `.gitignore` ignored every `state/` directory, so a vendored gRPC package was missing and CI could not build v0.3.0. Only the root `state/` is ignored now.

## v0.3.0 - 2026-11-05

- Link-first (prds#139), behind `WFC_LINK_FIRST=1` (needs the WFC database settings): an unlinked console's NAS login is refused with its link number as the return code (4000-9999), which every DWC title shows as "Error Code: 24000-29999" on the Wii and the DS. Mario Kart Wii is let through to wfc-gpcm, which shows it text. A banned console gets 3002 (23002), an unreadable link state 3102 (23102). The console's network is kept only as a keyed tag on its code (for a claim by number); the address comes from `X-Real-Ip`/`X-Forwarded-For` when the peer is the box's own front.
- Fix: `/w1` serves `payload/stage1.bin` when `stage1v1.bin` is absent (the wii-patcher releases ship one stage1.bin).
- Fix: the stage-2 salt hash covers the console's whole query up to `&h=` (the `d=` parameter of the v0.1.2+ payloads made every stage-2 request fail).
- Vendors wfc-common v0.8.0.

## v0.2.2 - 2026-11-03

- Internal `/internal/` polls are counted and summarised once a minute instead of logged per call (wfc-common v0.3.1).

## v0.2.1 - 2026-11-03

- Every HTTP request now logs one line to stderr, `http <method> <path> -> <status> <ms>ms` (path only, never the query, headers or bodies; `/healthz` skipped), through wfc-common v0.3.0 `httplog`. Covers the HTTP and HTTPS console listeners.


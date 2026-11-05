package nas

// Link-first (prds repos/wfc-link-first-design.md, card prds#139, step 0 results): a refused
// NAS login is the one refusal every DWC title shows the same way, "Error Code: 2xxxx" (20000 +
// the return code, for return codes 1..9999), on the Wii and the DS. So an unlinked console is
// refused here, at NAS login, with its link number as the return code (4000-9999, shown as
// 24000-29999). Mario Kart Wii is let through: wfc-gpcm refuses it with text (the 8-symbol code
// and the same number). WFC_LINK_FIRST=1 turns the gate on; it needs the WFC database.

import (
	"context"
	"encoding/binary"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"openpak/wfc-common/common"
	"openpak/wfc-common/database"
	"openpak/wfc-common/logging"
	"openpak/wfc-common/mkwboard"
)

const (
	returnBanned      = "3002" // shown 23002, like WiiLink's 22002 ban
	returnUnavailable = "3102" // shown 23102, like 22102
)

var (
	linkDB   *database.Connection
	linkCore database.CoreChecker
)

func startLinkGate(cfg common.Config) {
	if !database.LinkFirst() {
		logging.Notice("NAS", "link-first gate is off (WFC_LINK_FIRST unset)")
		return
	}
	c := database.Start(cfg)
	if err := c.EnsurePairing(); err != nil {
		logging.Error("NAS", "link gate: pairing tables:", err)
	}
	linkDB = &c
	if cfg.CoreAddress != "" {
		if links, _, err := mkwboard.DialCore(cfg.CoreAddress, cfg.CoreKey); err == nil {
			linkCore = links
		}
	}
	logging.Notice("NAS", "link-first gate is on: unlinked consoles are refused with their link number")
}

// clientIP is the console's address: the proxy's X-Real-Ip / X-Forwarded-For when the peer is a
// private address (the box's own front), the peer otherwise.
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	if ip := net.ParseIP(host); ip != nil && (ip.IsPrivate() || ip.IsLoopback()) {
		if v := strings.TrimSpace(r.Header.Get("X-Real-Ip")); v != "" {
			return v
		}
		if v := r.Header.Get("X-Forwarded-For"); v != "" {
			return strings.TrimSpace(strings.Split(v, ",")[0])
		}
	}
	return host
}

// gateLogin rewrites a successful NAS login reply into a refusal for an unlinked console.
func gateLogin(r *http.Request, moduleName string, fields map[string][]byte, reply map[string]string) {
	if linkDB == nil || !database.LinkFirst() {
		return
	}
	if rc := reply["returncd"]; rc != "001" && rc != "040" {
		return
	}
	gamecd := string(fields["gamecd"])
	if strings.HasPrefix(gamecd, "RMC") {
		return // Mario Kart Wii: refused at GPCM with text
	}
	uid, err := strconv.ParseUint(string(fields["userid"]), 10, 64)
	if err != nil || uid == 0 {
		return
	}
	refuse := func(rc string) {
		for k := range reply {
			delete(reply, k)
		}
		reply["retry"] = "0"
		reply["datetime"] = getDateTime()
		reply["returncd"] = rc
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	switch res, err := linkDB.LinkGate(ctx, uid, linkCore); res {
	case database.GateLinked:
		return
	case database.GateUnavailable:
		logging.Warn(moduleName, "link gate: refused user", uid, "reason=core_unavailable:", err)
		refuse(returnUnavailable)
		return
	}
	ip := clientIP(r)
	if banned, _, err := linkDB.ConsoleBanned(ctx, uid, 0, ip); err != nil {
		logging.Warn(moduleName, "link gate: ban check:", err)
		refuse(returnUnavailable)
		return
	} else if banned {
		logging.Notice(moduleName, "link gate: refused user", uid, "reason=console_banned")
		refuse(returnBanned)
		return
	}
	unit := byte(database.UnitWii)
	var order binary.ByteOrder = binary.BigEndian
	if u := fields["unitcd"]; len(u) == 1 && u[0] == '0' {
		unit, order = database.UnitDS, binary.LittleEndian
	}
	// The name the player can read on the console: the in-game name, else (older DS games send
	// it empty or not at all) the DS user name.
	name := ""
	if v, ok := fields["ingamesn"]; ok {
		name = strings.TrimSpace(strings.TrimRight(common.UTF16Decode(v, order), "\x00"))
	}
	if v, ok := fields["devname"]; ok && name == "" {
		name = strings.TrimSpace(strings.TrimRight(common.UTF16Decode(v, order), "\x00"))
	}
	gsbrcd := string(fields["gsbrcd"])
	if gsbrcd == "" && len(gamecd) >= 3 {
		gsbrcd = gamecd[:3] + "J"
	}
	p, err := linkDB.IssuePairingCode(ctx, database.PairingRequest{UserID: uid, Unit: unit, GameCode: gsbrcd, Name: name,
		NetTag: database.NetTag(database.NetTagKey(), ip, time.Now())})
	if err != nil {
		logging.Warn(moduleName, "link gate: refused user", uid, "reason=no_code:", err)
		refuse(returnUnavailable)
		return
	}
	unitName := "wii"
	if unit == database.UnitDS {
		unitName = "ds"
	}
	logging.Notice(moduleName, "link gate: refused user", uid, "unit", unitName, "game", gamecd,
		"reason=not_linked code="+database.CodeHash(p.Code))
	refuse(strconv.Itoa(p.Number))
}

var timeNow = time.Now

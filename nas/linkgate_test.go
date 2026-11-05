package nas

import (
	"context"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"openpak/wfc-common/common"
	"openpak/wfc-common/database"
)

// The NAS gate against a real Postgres (WFC_TEST_DSN): an unlinked console's successful login
// becomes a refusal with its link number; Mario Kart Wii and linked consoles pass.
func TestGateLogin(t *testing.T) {
	dsn := os.Getenv("WFC_TEST_DSN")
	if dsn == "" {
		t.Skip("WFC_TEST_DSN not set")
	}
	t.Setenv("WFC_LINK_FIRST", "1")
	t.Setenv("WFC_NET_TAG_KEY", "k")
	cfg := common.Config{}
	rest := strings.TrimPrefix(dsn, "postgres://")
	cred, hostdb, _ := strings.Cut(rest, "@")
	cfg.Username, cfg.Password, _ = strings.Cut(cred, ":")
	cfg.DatabaseAddress, cfg.DatabaseName, _ = strings.Cut(hostdb, "/")
	c := database.Start(cfg)
	defer c.Close()
	if err := c.EnsurePairing(); err != nil {
		t.Fatal(err)
	}
	linkDB = &c
	defer func() { linkDB = nil }()
	ctx := context.Background()
	const uid = 6660001
	clean := func() { _ = c.ForgetPairing(ctx, uid); _, _ = c.EraseUsers(ctx, []int64{uid}) }
	clean()
	defer clean()

	ok := func() map[string]string { return map[string]string{"returncd": "001", "token": "t", "challenge": "c"} }
	ds := map[string][]byte{"gamecd": []byte("AMCE"), "userid": []byte(strconv.Itoa(uid)), "unitcd": []byte("0"),
		"gsbrcd": []byte("AMCE"), "devname": []byte{'D', 0, 'S', 0}, "ingamesn": []byte{}}
	r := httptest.NewRequest("POST", "/ac", nil)
	r.RemoteAddr = "10.89.2.96:5000"
	r.Header.Set("X-Real-Ip", "198.51.100.9")

	reply := ok()
	gateLogin(r, "test", ds, reply)
	n, err := strconv.Atoi(reply["returncd"])
	if err != nil || n < 4000 || n > 9999 || reply["token"] != "" {
		t.Fatalf("unlinked DS: %v", reply)
	}
	shown := strconv.Itoa(database.DisplayNumber(n))
	// From the console's network the number resolves; from elsewhere it needs the name.
	if _, err := c.LookupPairingCode(ctx, shown, database.NumberProof{NetTags: database.NetTags([]byte("k"), "198.51.100.9", timeNow())}); err != nil {
		t.Fatalf("same network: %v", err)
	}
	if _, err := c.LookupPairingCode(ctx, shown, database.NumberProof{NetTags: database.NetTags([]byte("k"), "203.0.113.1", timeNow())}); err == nil {
		t.Fatal("other network resolved")
	}
	if tg, err := c.LookupPairingCode(ctx, shown, database.NumberProof{Name: "ds"}); err != nil || tg.Unit != database.UnitDS || tg.GameCode != "AMCE" {
		t.Fatalf("by name: %+v %v", tg, err)
	}

	mkw := map[string][]byte{"gamecd": []byte("RMCP"), "userid": []byte(strconv.Itoa(uid)), "unitcd": []byte("1")}
	reply = ok()
	gateLogin(r, "test", mkw, reply)
	if reply["returncd"] != "001" {
		t.Fatalf("Mario Kart Wii must reach GPCM: %v", reply)
	}

	tg, _, err := c.BeginApproval(ctx, shown, "acct-n", database.NumberProof{Name: "DS"})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.FinishApproval(ctx, tg, "acct-n", "link-n"); err != nil {
		t.Fatal(err)
	}
	reply = ok()
	gateLogin(r, "test", ds, reply)
	if reply["returncd"] != "001" || reply["token"] != "t" {
		t.Fatalf("linked DS refused: %v", reply)
	}
}

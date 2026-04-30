package house

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

func TestComposeDoorNoFromParts(t *testing.T) {
	got := composeDoorNoFromParts("1", "A", "103")
	if got != "1 A 103" {
		t.Fatalf("composeDoorNoFromParts failed, got=%q", got)
	}

	got = composeDoorNoFromParts(" 1号楼 ", " A单元 ", " 103室 ")
	if got != "1号楼 A单元 103室" {
		t.Fatalf("composeDoorNoFromParts trim failed, got=%q", got)
	}

	got = composeDoorNoFromParts("1", "", "103")
	if got != "" {
		t.Fatalf("composeDoorNoFromParts should return empty on missing part, got=%q", got)
	}
}

func TestValidateDoorNoByDict_FormatOnly(t *testing.T) {
	// community_id=0 时不会查库，但格式校验依然生效。
	xq := system.XiaoQu{CommunityId: 0}

	if err := validateDoorNoByDict(xq, "1 A 103"); err != nil {
		t.Fatalf("expected valid door no, got err=%v", err)
	}

	if err := validateDoorNoByDict(xq, "1号楼 A单元 103室"); err != nil {
		t.Fatalf("expected valid door no with suffix, got err=%v", err)
	}

	if err := validateDoorNoByDict(xq, "1号楼103室"); err == nil {
		t.Fatal("expected format error for missing split parts")
	}

	if err := validateDoorNoByDict(xq, "1 A"); err == nil {
		t.Fatal("expected format error for only two parts")
	}
}


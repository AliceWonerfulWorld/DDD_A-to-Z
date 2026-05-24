package postgres

import (
	"testing"

	guildtowndomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guildtown"
)

func TestShouldSpendPurchaseSP(t *testing.T) {
	t.Run("Commonは永続SP口座がないためSP消費しない", func(t *testing.T) {
		if shouldSpendPurchaseSP(guildtowndomain.BuildingMaster{PurchaseSP: 50, TargetSP: "Common"}) {
			t.Fatal("shouldSpendPurchaseSP() = true, 期待値 false")
		}
	})

	t.Run("言語別SPは消費対象にする", func(t *testing.T) {
		if !shouldSpendPurchaseSP(guildtowndomain.BuildingMaster{PurchaseSP: 50, TargetSP: "Go"}) {
			t.Fatal("shouldSpendPurchaseSP() = false, 期待値 true")
		}
	})
}

func TestShouldSpendUpgradeSP(t *testing.T) {
	t.Run("Commonは永続SP口座がないためSP消費しない", func(t *testing.T) {
		if shouldSpendUpgradeSP(guildtowndomain.BuildingLevelCost{SP: 50, TargetSP: "Common"}) {
			t.Fatal("shouldSpendUpgradeSP() = true, 期待値 false")
		}
	})

	t.Run("言語別SPは消費対象にする", func(t *testing.T) {
		if !shouldSpendUpgradeSP(guildtowndomain.BuildingLevelCost{SP: 50, TargetSP: "Go"}) {
			t.Fatal("shouldSpendUpgradeSP() = false, 期待値 true")
		}
	})
}

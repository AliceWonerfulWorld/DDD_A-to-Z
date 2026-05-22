package security

import (
	"context"
	"os/exec"
	"testing"
)

func TestMixWithAwk(t *testing.T) {
	if _, err := exec.LookPath("awk"); err != nil {
		t.Skip("awk command is not available")
	}

	got, err := MixWithAwk("abcdef", "SALT")
	if err != nil {
		t.Fatalf("MixWithAwk がエラーを返しました: %v", err)
	}

	want := "fedTcbaL::6:4"
	if got != want {
		t.Fatalf("MixWithAwk() = %q, 期待値 %q", got, want)
	}
}

func TestMixWithAwkIsDeterministic(t *testing.T) {
	if _, err := exec.LookPath("awk"); err != nil {
		t.Skip("awk command is not available")
	}

	first, err := MixWithAwk("state-token", "test-secret")
	if err != nil {
		t.Fatalf("1回目の MixWithAwk がエラーを返しました: %v", err)
	}
	second, err := MixWithAwk("state-token", "test-secret")
	if err != nil {
		t.Fatalf("2回目の MixWithAwk がエラーを返しました: %v", err)
	}

	if first != second {
		t.Fatalf("MixWithAwk() が冪等ではありません: 1回目 %q, 2回目 %q", first, second)
	}
}

func TestAwkTextMixerReturnsCommandError(t *testing.T) {
	mixer := &AwkTextMixer{command: "missing-awk-command-for-test"}

	if _, err := mixer.Mix(context.Background(), "abcdef", "SALT"); err == nil {
		t.Fatal("Mix() のエラー = nil, 期待値 コマンド実行エラー")
	}
}

package security

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const defaultAwkMixerTimeout = 2 * time.Second

//go:embed awk_engine/salt_mixer.awk
var saltMixerScript string

type TextMixer interface {
	Mix(ctx context.Context, input string, salt string) (string, error)
}

type AwkTextMixer struct {
	command string
	timeout time.Duration
}

func NewAwkTextMixer() *AwkTextMixer {
	return &AwkTextMixer{
		command: "awk",
		timeout: defaultAwkMixerTimeout,
	}
}

func MixWithAwk(input string, salt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultAwkMixerTimeout)
	defer cancel()

	return NewAwkTextMixer().Mix(ctx, input, salt)
}

func (m *AwkTextMixer) Mix(ctx context.Context, input string, salt string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if m == nil {
		m = NewAwkTextMixer()
	}
	if m.command == "" {
		m.command = "awk"
	}
	if m.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, m.timeout)
		defer cancel()
	}

	script, err := os.CreateTemp("", "salt_mixer_*.awk")
	if err != nil {
		return "", fmt.Errorf("create awk mixer script: %w", err)
	}
	defer func() {
		_ = os.Remove(script.Name())
	}()

	if _, err := script.WriteString(saltMixerScript); err != nil {
		_ = script.Close()
		return "", fmt.Errorf("write awk mixer script: %w", err)
	}
	if err := script.Close(); err != nil {
		return "", fmt.Errorf("close awk mixer script: %w", err)
	}

	command := exec.CommandContext(ctx, m.command, "-f", script.Name())
	command.Env = append(os.Environ(), "SALT_MIXER_SALT="+salt)
	command.Stdin = strings.NewReader(input)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("mix text with awk: %w", ctx.Err())
		}
		return "", fmt.Errorf("mix text with awk: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return stdout.String(), nil
}

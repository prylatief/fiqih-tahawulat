package handler

import (
	"context"
	"os/exec"
	"strings"
)

func translateText(ctx context.Context, text, targetLang string) string {
	cmd := exec.CommandContext(ctx, "trans", "-b", ":"+targetLang)
	cmd.Stdin = strings.NewReader(text)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

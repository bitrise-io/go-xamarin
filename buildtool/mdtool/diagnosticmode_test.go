package mdtool

import (
	"fmt"
	"testing"
	"time"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/stretchr/testify/require"
)

func TestRunCommandInDiagnosticMode(t *testing.T) {
	cmd := cmdex.NewCommand("/bin/bash", "-c", "echo pattern && sleep 100")
	now := time.Now()
	err := runCommandInDiagnosticMode(*cmd, "pattern", 10*time.Second, 10*time.Second, false)
	require.Equal(t, "timed out", err.Error())
	diff := time.Now().Sub(now)
	require.Equal(t, true, diff.Seconds() < 100, fmt.Sprintf("diff: %v", diff.Seconds()))
}

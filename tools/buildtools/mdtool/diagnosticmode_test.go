package mdtool

import (
	"fmt"
	"testing"
	"time"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/stretchr/testify/require"
)

func TestRunCommandInDiagnosticMode(t *testing.T) {
	t.Log("test without retry")
	{
		cmd := cmdex.NewCommand("/bin/bash", "-c", "echo pattern && sleep 100")
		now := time.Now()
		err := runCommandInDiagnosticMode(*cmd, "pattern", 2*time.Second, 2*time.Second, false)
		require.Equal(t, "timed out", err.Error())
		diff := time.Now().Sub(now)
		require.Equal(t, true, diff.Seconds() < 10, fmt.Sprintf("diff: %v", diff.Seconds()))
	}

	t.Log("test with retry")
	{
		cmd := cmdex.NewCommand("/bin/bash", "-c", "echo pattern && sleep 100")
		now := time.Now()
		err := runCommandInDiagnosticMode(*cmd, "pattern", 2*time.Second, 2*time.Second, true)
		require.Equal(t, "timed out", err.Error())
		diff := time.Now().Sub(now)
		require.Equal(t, true, diff.Seconds() < 20, fmt.Sprintf("diff: %v", diff.Seconds()))
	}
}

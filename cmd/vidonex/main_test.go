package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestCLI_CompletionTable(t *testing.T) {
	tests := []struct {
		name           string
		shell          string
		expectedSubstr string
	}{
		{
			name:           "bash_completion",
			shell:          "bash",
			expectedSubstr: "_vidonex_completion",
		},
		{
			name:           "zsh_completion",
			shell:          "zsh",
			expectedSubstr: "#compdef vidonex",
		},
		{
			name:           "fish_completion",
			shell:          "fish",
			expectedSubstr: "complete -c vidonex",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			oldStdout := os.Stdout
			reader, writer, err := os.Pipe()
			if err != nil {
				t.Fatalf("failed to create pipe: %v", err)
			}
			os.Stdout = writer

			executeCompletionCommand([]string{tc.shell})

			_ = writer.Close()
			os.Stdout = oldStdout

			var buffer bytes.Buffer
			_, _ = io.Copy(&buffer, reader)
			output := buffer.String()

			if !strings.Contains(output, tc.expectedSubstr) {
				t.Errorf("expected completion output to contain %q, got: %s", tc.expectedSubstr, output)
			}
		})
	}
}

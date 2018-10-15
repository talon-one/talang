package examples_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExamples(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	dir, err := ioutil.ReadDir(wd)
	require.NoError(t, err)

	for _, p := range dir {
		if p.IsDir() && !strings.HasPrefix(p.Name(), ".") {
			cmd := exec.Command("go", "run", filepath.Join(p.Name(), "main.go"))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			require.NoError(t, cmd.Run())
		}
	}
}

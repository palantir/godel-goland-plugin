// Copyright 2022 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package goland

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var (
	update = flag.Bool("update", false, "update the golden files of this test")
)

func TestRender(t *testing.T) {
	for _, tc := range []struct {
		name          string
		goldenFile    string
		usesGoModules bool
		envVars       map[string]string
	}{
		{
			name:          "no go modules",
			goldenFile:    "no-go-mod.ipr",
			usesGoModules: false,
			envVars:       map[string]string{},
		},
		{
			name:          "uses go modules with no env flags",
			goldenFile:    "go-mod-no-flags.ipr",
			usesGoModules: true,
			envVars:       map[string]string{},
		},
		{
			name:          "uses go modules with subset of env vars",
			goldenFile:    "go-mod-with-flags-subset.ipr",
			usesGoModules: true,
			envVars: map[string]string{
				"GOPRIVATE": "github.com/palantir/*",
				"GOPROXY":   "https://proxy.golang.org",
			},
		},
		{
			name:          "uses go modules with all env vars",
			goldenFile:    "go-mod-with-flags.ipr",
			usesGoModules: true,
			envVars: map[string]string{
				"GOFLAGS":   "-mod=vendor",
				"GONOPROXY": "none",
				"GOPRIVATE": "github.com/palantir/*",
				"GOPROXY":   "https://proxy.golang.org",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			projectDir := filepath.Join(t.TempDir(), "test-project")
			requireNoError(t, os.Mkdir(projectDir, 0777))

			if tc.usesGoModules {
				err := os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte("module foo"), 0777)
				requireNoError(t, err)
			}

			// setup relevant env vars before testing
			for _, v := range goEnvFlags {
				requireNoError(t, os.Unsetenv(v))
			}
			for k, v := range tc.envVars {
				requireNoError(t, os.Setenv(k, v))
			}

			err := CreateProjectFiles(projectDir)
			requireNoError(t, err)

			// read written file and compare
			gotBytes, err := os.ReadFile(filepath.Join([]string{projectDir, "test-project.ipr"}...))
			requireNoError(t, err)

			expected := goldenValue(t, tc.goldenFile, string(gotBytes), *update)
			if expected != string(gotBytes) {
				t.FailNow()
			}
		})
	}
}

func requireNoError(t *testing.T, err error) {
	if err != nil {
		t.FailNow()
	}
}

func goldenValue(t *testing.T, goldenFile string, actual string, update bool) string {
	goldenPath := "testdata/" + goldenFile + ".golden"

	if update {
		requireNoError(t, os.WriteFile(goldenPath, []byte(actual), 0644))
	}

	content, err := os.ReadFile(goldenPath)
	requireNoError(t, err)
	return string(content)
}

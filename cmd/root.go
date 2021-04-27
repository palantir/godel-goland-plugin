// Copyright 2016 Palantir Technologies, Inc.
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

package cmd

import (
	"github.com/palantir/godel-goland-plugin/goland"
	"github.com/palantir/godel/v2/framework/pluginapi"
	"github.com/palantir/pkg/cobracli"
	"github.com/spf13/cobra"
)

var (
	projectDirFlagVal string
)

var rootCmd = &cobra.Command{
	Use:   "goland-plugin",
	Short: "Create GoLand project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return goland.CreateProjectFiles(projectDirFlagVal)
	},
}

func Execute() int {
	return cobracli.ExecuteWithDefaultParams(rootCmd)
}

func init() {
	pluginapi.AddProjectDirPFlagPtr(rootCmd.PersistentFlags(), &projectDirFlagVal)
	if err := rootCmd.MarkPersistentFlagRequired(pluginapi.ProjectDirFlagName); err != nil {
		panic(err)
	}
}

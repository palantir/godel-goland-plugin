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

package goland

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/nmiyake/pkg/dirs"
	"github.com/pkg/errors"
)

const (
	defaultGoSDK             = "Go"
	imlGolandTemplateContent = `<?xml version="1.0" encoding="UTF-8"?>
<module type="WEB_MODULE" version="4">
  <component name="Go" enabled="true" />
  <component name="NewModuleRootManager" inherit-compiler-output="true">
    <exclude-output />
    <content url="file://$MODULE_DIR$" />
    <orderEntry type="sourceFolder" forTests="false" />
  </component>
</module>
`
	iprGolandTemplateContent = `<?xml version="1.0" encoding="UTF-8"?>
<project version="4">
  <component name="ProjectModuleManager">
    <modules>
      <module fileurl="file://$PROJECT_DIR$/{{.ProjectName}}.iml" filepath="$PROJECT_DIR$/{{.ProjectName}}.iml" />
    </modules>
  </component>
  <component name="ProjectTasksOptions">
    <TaskOptions isEnabled="true">
      <option name="arguments" value="format runAll $FilePathRelativeToProjectRoot$" />
      <option name="checkSyntaxErrors" value="true" />
      <option name="description" value="" />
      <option name="exitCodeBehavior" value="ERROR" />
      <option name="fileExtension" value="go" />
      <option name="immediateSync" value="false" />
      <option name="name" value="godel" />
      <option name="output" value="" />
      <option name="outputFilters">
        <array />
      </option>
      <option name="outputFromStdout" value="false" />
      <option name="program" value="$PROJECT_DIR$/godelw" />
      <option name="runOnExternalChanges" value="true" />
      <option name="scopeName" value="Changed Files" />
      <option name="trackOnlyRoot" value="false" />
      <option name="workingDir" value="$ProjectFileDir$" />
      <envs />
    </TaskOptions>
  </component>
</project>
`
)

func CreateProjectFiles(rootDir string) error {
	return createIDEAFiles(rootDir, imlGolandTemplateContent, iprGolandTemplateContent)
}

func createIDEAFiles(rootDir string, imlContent, iprContent string) error {
	projectName, err := projectNameFromDir(rootDir)
	if err != nil {
		return err
	}

	goRoot, err := dirs.GoRoot()
	if err != nil {
		return errors.Wrapf(err, "failed to determine GOROOT")
	}
	buffer := bytes.Buffer{}
	templateValues := map[string]string{
		"GoSDK":       defaultGoSDK,
		"GoRoot":      goRoot,
		"ProjectName": projectName,
	}
	imlTemplate := template.Must(template.New("iml").Parse(imlContent))
	if err := imlTemplate.Execute(&buffer, templateValues); err != nil {
		return errors.Wrapf(err, "failed to execute template %s with values %v", imlContent, templateValues)
	}

	imlFilePath := path.Join(rootDir, projectName+".iml")
	if err := ioutil.WriteFile(imlFilePath, buffer.Bytes(), 0644); err != nil {
		return errors.Wrapf(err, "failed to write .iml file to %s", imlFilePath)
	}

	iprTemplate := template.Must(template.New("modules").Parse(iprContent))
	buffer = bytes.Buffer{}
	if err := iprTemplate.Execute(&buffer, templateValues); err != nil {
		return errors.Wrapf(err, "failed to execute template %s with values %v", iprContent, templateValues)
	}

	iprFilePath := path.Join(rootDir, projectName+".ipr")
	if err := ioutil.WriteFile(iprFilePath, buffer.Bytes(), 0644); err != nil {
		return errors.Wrapf(err, "failed to write .ipr file to %s", iprFilePath)
	}

	return nil
}

func CleanProjectFiles(rootDir string) error {
	projectName, err := projectNameFromDir(rootDir)
	if err != nil {
		return err
	}

	for _, ext := range []string{"iml", "ipr", "iws"} {
		currPath := path.Join(rootDir, fmt.Sprintf("%v.%v", projectName, ext))
		if err := os.Remove(currPath); err != nil && !os.IsNotExist(err) {
			return errors.Wrapf(err, "failed to remove file %s", currPath)
		}
	}
	return nil
}

func projectNameFromDir(dir string) (string, error) {
	if !filepath.IsAbs(dir) {
		wd, err := os.Getwd()
		if err != nil {
			return "", errors.Wrapf(err, "failed to determine working directory")
		}
		dir = path.Join(wd, dir)
	}
	return path.Base(dir), nil
}

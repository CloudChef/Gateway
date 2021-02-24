// Copyright (c) 2021 上海骞云信息科技有限公司. All rights reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
)

type ScriptHandler struct {
	ScriptType       string `json:"scriptType"`
	ScriptContent    string `json:"scriptContent"`
	ScriptParameters string `json:"scriptParameters"`
	ScriptName       string `json:"scriptName"`
	ScriptPath       string
}

func (scriptHandler *ScriptHandler) Execute() ([]byte, error) {
	if !scriptHandler.Exist() {
		scriptHandler.Create()
	}

	if scriptHandler.ScriptType == "python" {
		log.Info("python", scriptHandler.ScriptPath, scriptHandler.ScriptParameters)
		cmd := exec.Command("python", scriptHandler.ScriptPath, scriptHandler.ScriptParameters)
		return cmd.CombinedOutput()
	} else if scriptHandler.ScriptType == "shell" {
		cmd := exec.Command("sh", scriptHandler.ScriptPath, scriptHandler.ScriptParameters)
		return cmd.CombinedOutput()
	} else {
		cmd := exec.Command("")
		return cmd.CombinedOutput()
	}
}

func (scriptHandler *ScriptHandler) Exist() bool {
	_, err := os.Stat(path.Join(SCRIPT_PATH, scriptHandler.ScriptName))
	if err == nil {
		scriptHandler.ScriptPath = path.Join(SCRIPT_PATH, scriptHandler.ScriptName)
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (scriptHandler *ScriptHandler) Create() {
	scriptHandler.ScriptPath = path.Join(SCRIPT_PATH, scriptHandler.ScriptName)
	log.Println(scriptHandler.ScriptPath)
	file, err := os.Create(scriptHandler.ScriptPath)
	if err != nil {
		log.Info("create", err)
	}

	_, err = file.Write([]byte(scriptHandler.ScriptContent))
	if err != nil {
		log.Info("write", err)
	}
	_ = file.Close()
}

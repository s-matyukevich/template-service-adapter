package utils

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

func ExecuteScript(script string, params interface{}) (interface{}, string, error) {
	if script == "" {
		return nil, "", nil
	}
	serializedParams, err := json.Marshal(params)
	if err != nil {
		return nil, "", err
	}
	cmd := exec.Command(script, string(serializedParams))
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return nil, "", err
	}
	output := stdout.Bytes()
	stderrOutput := string(stderr.Bytes())
	if output != nil {
		var res interface{}
		err = json.Unmarshal(output, &res)
		return res, stderrOutput, err
	}
	return nil, stderrOutput, nil
}

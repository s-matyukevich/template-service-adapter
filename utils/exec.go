package utils

import (
	"encoding/json"
	"os/exec"
)

func ExecuteScript(script string, params interface{}) (interface{}, error) {
	if script == "" {
		return "", nil
	}
	serializedParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(script, string(serializedParams))
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	if output != nil {
		var res interface{}
		err = json.Unmarshal(output, &res)
		return res, err
	}
	return nil, nil
}

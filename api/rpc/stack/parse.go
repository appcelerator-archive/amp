package stack

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type stackSpec struct {
	Services map[string]serviceSpec `yaml:"services"`
}

type serviceSpec struct {
	Image  string     `yaml:"image"`
	Deploy deploySpec `yaml:"deploy"`
}

type deploySpec struct {
	Labels interface{} `yaml:"labels"`
}

// parseStackfile update stack structure with the services labels
func parseStack(stack *Stack) error {
	if stack.FileData == "" {
		fullFileName := fmt.Sprintf("/var/lib/amp/%s.yml", stack.Name)

		data, err := ResolvedComposeFileVariables(fullFileName, fmt.Sprintf("%s/%s", stackFilePath, StackFileVarName), stack.AmpTag)
		if err != nil {
			return err
		}
		stack.FileData = data
	}
	var specs stackSpec
	if err := yaml.Unmarshal([]byte(stack.FileData), &specs); err != nil {
		return err
	}
	for name, spec := range specs.Services {
		// try to parse labels as a map
		// else try to parse labels as string entries
		var labels = map[string]string{}
		labelSpec := spec.Deploy.Labels
		if labelMap, ok := labelSpec.(map[interface{}]interface{}); ok {
			for k, v := range labelMap {
				labels[k.(string)] = v.(string)
			}
		} else if labelList, ok := labelSpec.([]interface{}); ok {
			for _, s := range labelList {
				a := strings.Split(s.(string), "=")
				labels[a[0]] = a[1]
			}
		}

		stack.Services = append(stack.Services, &ServiceSpec{
			Name:   name,
			Image:  spec.Image,
			Labels: labels,
		})
	}
	return nil
}

// LoadInfraVariables load variable from amp.var file to map
func LoadInfraVariables(varFilePath string, ampTag string) (map[string]string, error) {
	data, err := ioutil.ReadFile(varFilePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")

	//Load map
	varMap := make(map[string]string)
	for _, line := range lines {
		cmd := strings.Split(line, "=")
		if len(cmd) == 2 {
			name := strings.Trim(cmd[0], " ")
			val := strings.Trim(cmd[1], " ")
			varMap[name] = val
		}
	}
	if ampTag != "" {
		varMap["versionAmp"] = ampTag
	}

	//update var
	for name, val := range varMap {
		if lb := strings.Index(val, "${"); lb >= 0 {
			if le := strings.Index(val[lb+1:], "}"); le > 0 {
				namer := val[lb+2 : lb+le+1]
				if valv, ok := varMap[namer]; ok {
					varMap[name] = strings.Replace(val, "${"+namer+"}", valv, -1)
				}
			}
		}
	}
	return varMap, nil
}

// ResolvedComposeFileVariables replace variables by their values
func ResolvedComposeFileVariables(filePath string, varFilePath string, ampTag string) (string, error) {
	var varMap map[string]string
	if varFilePath == "" {
		varMap = make(map[string]string)
	} else {
		varm, errl := LoadInfraVariables(varFilePath, ampTag)
		if errl != nil {
			return "", errl
		}
		varMap = varm
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	retData := string(data)
	lb := strings.Index(retData, "${")
	var resolveErr error
	for lb >= 0 {
		if retData, resolveErr = ResolveVariableAt(retData, lb, varMap); resolveErr != nil {
			return "", resolveErr
		}
		lb = strings.Index(retData, "${")
	}
	return retData, nil
}

// ResolveVariableAt Resolve compose file variable, compose format: ${variableName:-defaultValue} where :-defaultValue is optionnal
func ResolveVariableAt(data string, lb int, varMap map[string]string) (string, error) {
	if le := strings.Index(data[lb+1:], "}"); le > 0 {
		varName := data[lb+2 : lb+le+1]
		if strings.Index(varName, "${") > 0 {
			le = lb + 30
			if lb+30 > len(data) {
				le = len(data)
			}
			return "", fmt.Errorf("bad formated compose file: finc '${' without '}' arround: ${%s", data[lb:le-1])
		}
		varList := strings.Split(varName, ":-")
		if val, ok := varMap[varList[0]]; ok {
			//replace variable by variable value
			data = strings.Replace(data, "${"+varName+"}", val, -1)
		} else if len(varList) > 1 {
			//replace variable by default value
			data = strings.Replace(data, "${"+varName+"}", varList[1], -1)
		} else {
			return "", fmt.Errorf("Found variable without value or default value %s", varName)
		}

	} else {
		le = lb + 30
		if lb+30 > len(data) {
			le = len(data)
		}
		return "", fmt.Errorf("bad formated compose file2: finc '${' without '}' arround: %s", data[lb:le-1])
	}
	return data, nil
}

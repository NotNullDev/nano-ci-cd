package config

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var EnvFromFileMap map[string]string

func areEnvValid(envs ...string) bool {
	for _, env := range envs {
		if env == "" {
			return false
		}
	}
	return true
}

func LoadEnvs(envs map[string]string) {
	for key, val := range envs {
		os.Setenv(key, val)
	}
}

func ParseEnvFiles(failOnMissingEnvFile bool, envFilePath ...string) (map[string]string, error) {
	var result map[string]string = make(map[string]string)

	for _, filePath := range envFilePath {
		file, err := os.Open(filePath)

		if err != nil {
			if failOnMissingEnvFile {
				return nil, err
			}
			continue
		}

		scanner := bufio.NewScanner(file)

		var lines []string

		for {
			lines = append(lines, scanner.Text())
			if !scanner.Scan() {
				break
			}
		}

		for lineNumber, line := range lines {

			commentRegex := regexp.MustCompile("#.+$")

			comment := commentRegex.Find([]byte(line))

			line = strings.Replace(line, string(comment), "", -1)
			line = strings.TrimSpace(line)

			if line == "" {
				continue
			}

			splitted := strings.Split(line, "=")

			if len(splitted) != 2 {
				return nil, fmt.Errorf("could not parse line %d in file [%s]", lineNumber, filePath)
			}

			envKey := splitted[0]
			envKey = strings.Replace(envKey, "export", "", -1)
			envKey = strings.TrimSpace(envKey)

			envVal := splitted[1]
			envVal = strings.Replace(envVal, "\"", "", -1)
			envVal = strings.Replace(envVal, "'", "", -1)
			envVal = strings.Replace(envVal, "`", "", -1)
			envVal = strings.TrimSpace(envVal)

			if envKey == "" {
				return nil, fmt.Errorf("key at line %d is empty in file [%s]", lineNumber, filePath)
			}

			if envVal == "" {
				return nil, fmt.Errorf("value at line %d is empty in file [%s]", lineNumber, filePath)
			}

			result[envKey] = envVal
		}
	}

	return result, nil
}

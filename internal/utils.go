package internal

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func (config *ConfigFile) LoadConfig(path string) {
	// read config file
	configFile, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("[ERROR] Could not read config file: %s", err)
	}

	// unmarshal config file into config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("[ERROR] Could not unmarshal config file: %s", err)
	}
}

func CorrectedArch() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x64"
	case "386":
		return "ia32" // afaik only windows got 32 bit support
	default:
		return runtime.GOARCH
	}
}

func CorrectedOS() string {
	if runtime.GOOS == "windows" {
		return "win32"
	}
	return runtime.GOOS
}

func ShellCommand() (program string, input string, sep string) {
	program = "bash"
	input = "-c"
	sep = ":"

	if runtime.GOOS == "windows" {
		program = "cmd"
		input = "/c"
		sep = ";"
	}

	return program, input, sep
}

func CreateLog(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(path)
}

func (config ConfigFile) SetEnv() {
	for _, val := range config.EnvVars {
		log.Printf("[ENV] Set Variable: %s = %s", val.Key, val.Value)
		os.Setenv(val.Key, val.Value)
	}
}

func AssetIndex(version string) (index string) {
	if version == "1.7.10" {
		return version
	}

	l := strings.Split(version, ".")

	return strings.Join(l[:len(l)-1], ".")
}

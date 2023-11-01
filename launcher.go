package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	internal "github.com/KotonBads/llg/internal"
	utils "github.com/KotonBads/llgutils"
)

func main() {
	// logging
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	file, err := utils.CreateLog(fmt.Sprintf("launcherlogs/%s.log", timestamp))

	if err != nil {
		log.Printf("[WARN] Could not create a log file: %s", err)
	}

	log.SetOutput(file)

	var config internal.ConfigFile

	// read config file
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("[ERROR] Could not read config file: %s", err)
	}

	// unmarshal config file into config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("[ERROR] Could not unmarshal config file: %s", err)
	}

	launchbody := utils.LaunchBody{
		OS:      "linux",
		Arch:    "x64",
		Version: "1.8.9",
		Module:  "forge",
	}

	launchmeta, _ := utils.FetchLaunchMeta(launchbody)
	launchmeta.DownloadArtifacts(config.WorkingDirectory)
	launchmeta.DownloadCosmetics(config.WorkingDirectory + "/textures")

	classpath, ichorClassPath, external, natives := launchmeta.SortFiles(config.WorkingDirectory)

	for _, val := range natives {
		utils.Unzip(val, config.WorkingDirectory+"/natives")
	}

	args := internal.MinecraftArgs{
		BaseArgs: []string{"--add-modules",
			"jdk.naming.dns",
			"--add-exports",
			"jdk.naming.dns/com.sun.jndi.dns=java.naming",
			"-Djna.boot.library.path=" + config.WorkingDirectory + "/natives",
			"-Djava.library.path=" + config.WorkingDirectory + "/natives",
			"-Dlog4j2.formatMsgNoLookups=true",
			"--add-opens",
			"java.base/java.io=ALL-UNNAMED",
			"-XX:+UseStringDeduplication",
			"-Dichor.prebakeClasses=false",
			"-Dlunar.webosr.url=file:index.html"},
		JVMArgs:            config.JVMArgs,
		Classpath:          classpath,
		IchorClassPath:     ichorClassPath,
		IchorExternalFiles: external,
		RAM:                config.Memory,
		Width:              config.Width,
		Height:             config.Height,
		MainClass:          launchmeta.LaunchTypeData.MainClass,
		Version:            launchbody.Version,
		AssetIndex:         internal.AssetIndex(launchbody.Version),
		GameDir:            config.GameDirectory,
		TexturesDir:        config.WorkingDirectory + "/textures",
		WebOSRDir:          config.WorkingDirectory + "/natives",
		WorkingDir:         config.WorkingDirectory,
		ClassPathDir:       config.WorkingDirectory,
		Fullscreen:         config.Fullscreen,
	}

	program := "bash"
	input := "-c"
	sep := ":"

	if runtime.GOOS == "win32" {
		program = "cmd"
		input = "/c"
		sep = ";"
	}

	cmd := exec.Command(program, input, fmt.Sprintf("%s %s", config.JRE, args.CompileArgs(sep)))

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	fmt.Printf("\nExecuting: \n%s\n\n", strings.Join(cmd.Args, " "))

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

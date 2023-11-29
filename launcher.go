package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	internal "github.com/KotonBads/llg/internal"
	utils "github.com/KotonBads/llgutils"
)

func main() {
	// logging
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	file, err := internal.CreateLog(fmt.Sprintf("launcherlogs/%s.log", timestamp))

	if err != nil {
		fmt.Printf("[WARN] Could not create a log file: %s", err)
	} else {
		log.SetOutput(file)
	}

	version := flag.String("version", "1.8.9", "Minecraft Version")
	module := flag.String("module", "lunar", "Module to use")
	cf := flag.String("config", "config.json", "Config file")

	flag.Parse()

	var config internal.ConfigFile
	config.LoadConfig(*cf)

	launchbody := utils.LaunchBody{
		OS:      internal.CorrectedOS(),
		Arch:    internal.CorrectedArch(),
		Version: *version,
		Module:  *module,
	}

	fmt.Printf("Minecraft %s\n", *version)
	fmt.Printf("%s module\n\n", *module)
	fmt.Printf("See logs: launcherlogs/%s.log\n\n", timestamp)
	fmt.Println("Downloading Assets...")

	launchmeta, _ := launchbody.FetchLaunchMeta()
	launchmeta.DownloadArtifacts(config.WorkingDirectory)
	launchmeta.DownloadCosmetics(config.WorkingDirectory + "/textures")

	classpath, ichorClassPath, external, natives := launchmeta.SortFiles(config.WorkingDirectory)

	fmt.Println("Extracting Natives...")
	for _, val := range natives {
		utils.Unzip(val, config.WorkingDirectory+"/natives")
	}

	config.SetEnv()

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
			"-Dichor.prebakeClasses=false",
			"-Dlunar.webosr.url=file:index.html"},
		JVMArgs:            config.JVMArgs,
		Classpath:          classpath,
		IchorClassPath:     ichorClassPath,
		IchorExternalFiles: external,
		JavaAgents:         config.JavaAgents,
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

	program, input, sep := internal.ShellCommand()

	cmd := exec.Command(program, input, fmt.Sprintf("%s %s", config.JRE, args.CompileArgs(sep)))

	if len(config.PreJava) != 0 {
		cmd = exec.Command(program, input, fmt.Sprintf("%s %s %s", config.PreJava, config.JRE, args.CompileArgs(sep)))
	}

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	fmt.Printf("\nExecuting: \n%s\n\n", strings.Join(cmd.Args, " "))
	log.Printf("[LAUNCH] Full cmdline: %s", strings.Join(cmd.Args, " "))

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

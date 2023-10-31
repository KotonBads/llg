package main

import (
	"bytes"
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

	launchbody := utils.LaunchBody{
		OS:      "linux",
		Arch:    "x64",
		Version: "1.8.9",
		Module:  "lunar",
	}
	launchmeta, _ := utils.FetchLaunchMeta(launchbody)
	launchmeta.DownloadArtifacts("temp")
	launchmeta.DownloadCosmetics("temp/textures")

	classpath, external, natives := launchmeta.SortFiles("temp")

	for _, val := range natives {
		utils.Unzip(val, "temp/natives")
	}

	ram := internal.RamAlloc{
		Xmx: 3072,
		Xms: 3072,
		Xmn: 1024,
		Xss: 2,
	}
	args := internal.MinecraftArgs{
		BaseArgs: []string{"--add-modules",
			"jdk.naming.dns",
			"--add-exports",
			"jdk.naming.dns/com.sun.jndi.dns=java.naming",
			"-Djna.boot.library.path=temp/natives",
			"-Djava.library.path=temp/natives",
			"-Dlog4j2.formatMsgNoLookups=true",
			"--add-opens",
			"java.base/java.io=ALL-UNNAMED",
			"-XX:+UseStringDeduplication",
			"-Dichor.prebakeClasses=false",
			"-Dlunar.webosr.url=file:index.html"},
		JVMArgs: []string{"-XX:+UnlockExperimentalVMOptions",
			"-XX:+UnlockDiagnosticVMOptions",
			"-XX:+AlwaysActAsServerClassMachine",
			"-XX:+AlwaysPreTouch",
			"-XX:+DisableExplicitGC",
			"-XX:+UseNUMA",
			"-XX:AllocatePrefetchStyle=3",
			"-XX:NmethodSweepActivity=1",
			"-XX:+UseG1GC",
			"-XX:MaxGCPauseMillis=37",
			"-XX:+PerfDisableSharedMem",
			"-XX:G1HeapRegionSize=16M",
			"-XX:G1NewSizePercent=23",
			"-XX:G1ReservePercent=20",
			"-XX:SurvivorRatio=32",
			"-XX:G1MixedGCCountTarget=3",
			"-XX:G1HeapWastePercent=20",
			"-XX:InitiatingHeapOccupancyPercent=10",
			"-XX:G1RSetUpdatingPauseTimePercent=0",
			"-XX:MaxTenuringThreshold=1",
			"-XX:G1SATBBufferEnqueueingThresholdPercent=30",
			"-XX:G1ConcMarkStepDurationMillis=5.0",
			"-XX:G1ConcRSHotCardLimit=16",
			"-XX:G1ConcRefinementServiceIntervalMillis=150",
			"-XX:GCTimeRatio=99",
			"-XX:ReservedCodeCacheSize=400M",
			"-XX:NonNMethodCodeHeapSize=12M",
			"-XX:ProfiledCodeHeapSize=194M",
			"-XX:NonProfiledCodeHeapSize=194M",
			"-XX:-DontCompileHugeMethods",
			"-XX:MaxNodeLimit=240000",
			"-XX:NodeLimitFudgeFactor=8000",
			"-XX:+UseVectorCmov",
			"-XX:+PerfDisableSharedMem",
			"-XX:+UseFastUnorderedTimeStamps",
			"-XX:+UseCriticalJavaThreadPriority",
			"-XX:+EagerJVMCI",
			"-XX:+UseTransparentHugePages",
			"-Dgraal.TuneInlinerExploration=1",
			"-Dgraal.CompilerConfiguration=enterprise"},
		Classpath:      classpath,
		IchorClassPath: external,
		RAM:            ram,
		Width:          1366,
		Height:         768,
		MainClass:      "com.moonsworth.lunar.genesis.Genesis",
		Version:        "1.8.9",
		AssetIndex:     "1.8",
		GameDir:        "~/.minecraft",
		TexturesDir:    "temp/textures",
		WebOSRDir:      "temp/natives",
		WorkingDir:     "temp",
		ClassPathDir:   "temp",
		Fullscreen:     true,
	}

	program := "bash"
	input := "-c"

	if runtime.GOOS == "win32" {
		program = "cmd"
		input = "/c"
	}

	cmd := exec.Command(program, input, "/home/koton-bads/.lunarclient/custom/graalvm-ee-java17-22.3.0/bin/java "+args.CompileArgs(":"))

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	fmt.Printf("\nExecuting: \n%s\n\n", strings.Join(cmd.Args, " "))

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

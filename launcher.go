package main

import (
	"fmt"
	"log"
	"time"

	utils "github.com/KotonBads/llgutils"
	internal "llg/internal"
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
		Module:  "forge",
	}
	launchmeta, _ := utils.FetchLaunchMeta(launchbody)
	launchmeta.DownloadArtifacts("temp")
	launchmeta.DownloadCosmetics("temp/textures")
}

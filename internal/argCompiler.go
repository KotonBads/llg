package internal

import (
	"fmt"
	"strings"
)

func (args MinecraftArgs) CompileArgs(sep string) string {
	var final []string

	final = append(final, strings.Join(args.BaseArgs, " "))
	final = append(final, fmt.Sprintf("-Xms%dM -Xmx%dM -Xss%dM -Xmn%dM", args.RAM.Xms, args.RAM.Xmx, args.RAM.Xss, args.RAM.Xmn))
	final = append(final, strings.Join(args.JVMArgs, " "))
	final = append(final, "-cp")
	final = append(final, strings.Join(args.Classpath, sep))

	// special case for java agents
	if len(args.JavaAgents) != 0 {
		for _, val := range args.JavaAgents {
			final = append(final, "-javaagent:"+val)
		}
	}

	final = append(final, args.MainClass)
	final = append(final, "--version "+args.Version)
	final = append(final, "--assetIndex "+args.AssetIndex)
	final = append(final, "--userProperties {}")
	final = append(final, "--gameDir "+args.GameDir)
	final = append(final, "--texturesDir "+args.TexturesDir)
	final = append(final, "--webosrDir "+args.WebOSRDir)
	final = append(final, "--launcherVersion 3.0.0")
	final = append(final, fmt.Sprintf("--width %d", args.Width))
	final = append(final, fmt.Sprintf("--height %d", args.Height))
	final = append(final, "--workingDirectory "+args.WorkingDir)
	final = append(final, "--classpathDir "+args.ClassPathDir)
	final = append(final, "--ichorClassPath "+strings.Join(args.IchorClassPath, ","))
	final = append(final, "--ichorExternalFiles "+strings.Join(args.IchorExternalFiles, ","))
	final = append(final, fmt.Sprintf("--fullscreen %t", args.Fullscreen))

	return strings.Join(final, " ")
}

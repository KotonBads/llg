package internal

import (
	"fmt"
	"strings"
)

func (args MinecraftArgs) CompileArgs(sep string) string {
	var final []string

	final = append(final, strings.Join(args.BaseArgs, " "))
	final = append(final, fmt.Sprintf("-Xms%d -Xmx%d -Xss%d -Xmn%d", args.RAM.Xms, args.RAM.Xmx, args.RAM.Xss, args.RAM.Xmn))
	final = append(final, strings.Join(args.JVMArgs, " "))
	final = append(final, strings.Join(args.Classpath, sep))
	final = append(final, args.MainClass)
	final = append(final, "--version "+args.Version)
	final = append(final, "--accessToken 0")
	final = append(final, "--assetIndex "+args.AssetIndex)
	final = append(final, "--userProperties {}")
	final = append(final, "--gameDir "+args.GameDir)
	final = append(final, "--texturesDir "+args.TexturesDir)
	final = append(final, "--webosrDir "+args.WebOSRDir)
	final = append(final, "--launcherVersion 3.0.0")
	final = append(final, "--hwid 0")
	final = append(final, fmt.Sprintf("--width %d", args.Width))
	final = append(final, fmt.Sprintf("--height %d", args.Height))
	final = append(final, "--workingDirectory "+args.WorkingDir)
	final = append(final, "--classpathDir "+args.ClassPathDir)
	final = append(final, "--ichorClassPath "+strings.Join(args.Classpath, ","))
	final = append(final, "--ichorExternalFiles "+strings.Join(args.IchorClassPath, ","))

	return strings.Join(final, " ")
}

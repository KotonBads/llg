package internal

type RamAlloc struct {
	Xmx int
	Xms int
	Xmn int
	Xss int
}

type MinecraftArgs struct {
	BaseArgs       []string
	JVMArgs        []string
	Classpath      []string
	IchorClassPath []string
	RAM            RamAlloc
	Width          int
	Height         int
	MainClass      string
	Version        string
	AssetIndex     string
	GameDir        string
	TexturesDir    string
	UIDir          string
	WebOSRDir      string
	WorkingDir     string
	ClassPathDir   string
	Fullscreen     bool
}

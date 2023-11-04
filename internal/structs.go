package internal

type Memory struct {
	Xmx int `json:"Xmx"`
	Xms int `json:"Xms"`
	Xmn int `json:"Xmn"`
	Xss int `json:"Xss"`
}

type MinecraftArgs struct {
	BaseArgs           []string
	JVMArgs            []string
	Classpath          []string
	IchorClassPath     []string
	IchorExternalFiles []string
	JavaAgents         []string
	RAM                Memory
	Width              int
	Height             int
	MainClass          string
	Version            string
	AssetIndex         string
	GameDir            string
	TexturesDir        string
	UIDir              string
	WebOSRDir          string
	WorkingDir         string
	ClassPathDir       string
	Fullscreen         bool
}

type ConfigFile struct {
	JRE              string   `json:"JRE"`
	Memory           Memory   `json:"Memory"`
	WorkingDirectory string   `json:"WorkingDirectory"`
	GameDirectory    string   `json:"GameDirectory"`
	PreJava          string   `json:"PreJava"`
	JavaAgents       []string `json:"JavaAgents"`
	EnvVars          []struct {
		Key   string `json:"Key"`
		Value string `json:"Value"`
	} `json:"EnvironmentVariables"`
	Width      int      `json:"Width"`
	Height     int      `json:"Height"`
	Fullscreen bool     `json:"Fullscreen"`
	JVMArgs    []string `json:"JVMArgs"`
}

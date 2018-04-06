package popConstants

const (
	UsW = "us-w"
	UsE = "us-e"
	JP  = "jp"
	Cf3 = "cf3"
	AZ  = "az"
	FF  = "ff"
)

type PoP struct {
	Api      string
	Uaa      string
	Passcode string
	Name     string
}

var USW = PoP{
	Api:      "https://api.system.aws-usw02-pr.ice.predix.io",
	Uaa:      "https://uaa.system.aws-usw02-pr.ice.predix.io",
	Passcode: "https://login.system.aws-usw02-pr.ice.predix.io/passcode",
	Name:     "west",
}

var USE = PoP{
	Api:      "https://api.system.asv-pr.ice.predix.io",
	Uaa:      "https://uaa.system.asv-pr.ice.predix.io",
	Passcode: "",
	Name:     "east",
}

var CF3 = PoP{
	Api:      "https://api.system.aws-usw02-dev.ice.predix.io",
	Uaa:      "https://uaa.system.aws-usw02-dev.ice.predix.io",
	Passcode: "https://login.system.aws-usw02-dev.ice.predix.io/passcode",
	Name:     "cf3",
}

var JPN = PoP{
	Api:      "https://api.system.aws-jp01-pr.ice.predix.io",
	Uaa:      "https://uaa.system.aws-jp01-pr.ice.predix.io",
	Passcode: "https://login.system.aws-jp01-pr.ice.predix.io/passcode",
	Name:     "jp",
}

var FFT = PoP{
	Api:      "https://api.system.aws-eu-central-1-pr.ice.predix.io",
	Uaa:      "https://uaa.system.aws-eu-central-1-pr.ice.predix.io",
	Passcode: "https://login.system.aws-eu-central-1-pr.ice.predix.io/passcode",
	Name:     "fft",
}

var AZR = PoP{
	Api:      "https://api.system.azr-usw01-pr.ice.predix.io",
	Uaa:      "https://uaa.system.azr-usw01-pr.ice.predix.io",
	Passcode: "",
	Name:     "azr",
}

var PoPs map[string]*PoP

func init() {
	PoPs = map[string]*PoP{
		UsW: &USW,
		UsE: &USE,
		JP:  &JPN,
		Cf3: &CF3,
		AZ:  &AZR,
		FF:  &FFT,
	}
}

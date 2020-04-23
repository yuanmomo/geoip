package list

import (
	"github.com/yuanmomo/geoip/util"
	"strings"
)
type Telegram struct{}

func (f *Telegram) Name() string {
	return "telegram"
}


func (a *Telegram) FetchAndMerge() []string {
	var TelegramCIDRArray []string
	tgCidrTxt := util.HttpGet("https://core.telegram.org/resources/cidr.txt")

	if len(strings.TrimSpace(tgCidrTxt)) > 0 { // not empty
		TelegramCIDRArray = strings.Split(tgCidrTxt,"\n")
	}

	return Merge(a,TelegramCIDRArray);
}

func init() {
	RegisterCommand(&Telegram{})
}

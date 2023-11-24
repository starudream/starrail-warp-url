package main

import (
	"github.com/skip2/go-qrcode"

	"github.com/starudream/go-lib/core/v2/utils/osutil"
)

func genQRCode(content string) string {
	code, err := qrcode.New(content, qrcode.Medium)
	osutil.PanicErr(err)
	return code.ToSmallString(false)
}

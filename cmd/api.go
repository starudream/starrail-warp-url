package main

import (
	"fmt"

	"github.com/starudream/go-lib/core/v2/slog"
	"github.com/starudream/go-lib/resty/v2"
)

type baseResp struct {
	RetCode *int   `json:"retcode"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (t *baseResp) GetRetCode() int {
	if t == nil || t.RetCode == nil {
		return 999999
	}
	return *t.RetCode
}

func (t *baseResp) IsSuccess() bool {
	return t != nil && t.RetCode != nil && *t.RetCode == 0
}

func (t *baseResp) String() string {
	if t == nil || t.RetCode == nil {
		return "<nil>"
	}
	return fmt.Sprintf("retcode: %d, message: %s, data: %v", *t.RetCode, t.Message, t.Data)
}

func checkGachaIsValid(gachaURL string) bool {
	_, err := resty.ParseResp[*baseResp, *baseResp](
		resty.R().SetError(&baseResp{}).SetResult(&baseResp{}).Get(gachaURL),
	)
	if err != nil {
		slog.Error("cannot get gacha data: %v", err)
	}
	return err == nil
}

package models

import (
	"encoding/json"
	"log"

	"github.com/pebbe/util"
)

func Pretty(e interface{}) string {
	x, err := json.MarshalIndent(e, "", "  ")
	util.WarnErr(err)
	return string(x)
}

func Pp(e interface{}) {
	log.Println(Pretty(e))
}

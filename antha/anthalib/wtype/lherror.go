package wtype

import (
	"errors"
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

// consts for liquid handling planner errors

const (
	LH_OK = iota
	LH_ERR_NO_DECK_SPACE
	LH_ERR_NO_TIPS
	LH_ERR_NOT_IMPLEMENTED
	LH_ERR_CONC
	LH_ERR_DRIV
	LH_ERR_POLICY
	LH_ERR_VOL
	LH_ERR_DIRE
	LH_ERR_OTHER
	LH_ERR_TIP_WASTE
)

func ErrorName(code int) string {
	errornames := [...]string{"LH_OK", "LH_ERR_NO_DECK_SPACE", "LH_ERR_NO_TIPS", "LH_ERR_NOT_IMPLEMENTED", "LH_ERR_CONC", "LH_ERR_DRIV", "LH_ERR_POLICY", "LH_ERR_VOL", "LH_ERR_DIRE", "LH_ERR_OTHER", "LH_ERR_TIP_WASTE"}
	return errornames[code]
}

func ErrorDesc(code int) string {
	errorboilerplate := [...]string{"no problem", "insufficient deck space to fit all required items; this may be due to constraints", "ran out of tips", "a required command is not implemented", "error calculating required volume for target concentration", "driver error", "liquid handling policy error", "volume error", "an internal error", "", "no space to safely drop tips"}
	return errorboilerplate[code]
}

func LHError(code int, detail string) error {
	// format here: error code (name) : general description : specific problem
	s := fmt.Sprintf("%d (%s) : %s : %s", code, ErrorName(code), ErrorDesc(code), detail)

	return errors.New(s)
}

//LHErrorf convenience method for formatting error messages
func LHErrorf(code int, format string, values ...interface{}) error {
	return LHError(code, fmt.Sprintf(format, values...))
}

func LHErrorCodeFromErr(err error) int {
	tx := strings.Split(err.Error(), " ")

	i := wutil.ParseInt(tx[0])

	return i
}

func LHErrorIsInternal(err error) bool {
	c := LHErrorCodeFromErr(err)

	return c == LH_ERR_DIRE
}

package handlers

import (
	"fmt"
	"math/big"
	"net/http"
	"regexp"

	"github.com/SAP/go-hdb/driver"
)

//DecToFloat converts go-hdb Decimal value to float64
func decToFloat(d driver.Decimal) float64 {
	out, _ := (*big.Rat)(&d).Float64()
	return out
}

func floatToDec(f float64) driver.Decimal {
	var z big.Rat
	z.SetFloat64(f)
	out := (*driver.Decimal)(&z)
	return *out
}

func onReadFail(s *Service, err error, w http.ResponseWriter) {
	s.Error("couldn't read rows: %s", err)
}

func onEmptyTable(s *Service, w http.ResponseWriter, t string) {
	s.Warning("table %s is empty", t)
}

func (srv *Service) seqNextVal(t string) int {
	var pattern = regexp.MustCompile(`SWP_(.*)$`)
	tn := pattern.FindString(t)
	var out int
	srv.Db.Get(&out, fmt.Sprintf(`SELECT "YBS_HC.SWP.SEQ::SEQ_%s".NEXTVAL FROM DUMMY`, tn[4:]))
	return out
}

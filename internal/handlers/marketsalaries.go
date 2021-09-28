package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/ts.tarusov/swp/internal/dbmodels"
	"github.com/ts.tarusov/swp/pkg/models"
)

// TODO remove unneccessary types
func (srv *Service) PostMarketSalaries(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	fid, err := strconv.Atoi(q.Get("function_id"))
	var years []int
	// Get gets the first value associated with the given key.
	// If there are no values associated with the key,
	// Get returns the empty string.
	// To access multiple values, use the map directly.
	// https://golang.org/pkg/net/url/#Values
	for _, v := range q["year"] {
		yr, err := strconv.Atoi(v)
		if err != nil {
			srv.Error("wrong request: %s", err)
			http.Error(w, fmt.Sprintf("wrong request: %s", err), http.StatusBadRequest)
			years = append(years, yr)
		}
	}
	sort.Ints(years)

	if err != nil {
		srv.Error("wrong request: %s", err)
		http.Error(w, fmt.Sprintf("wrong request: %s", err), http.StatusBadRequest)
	}
	ind_rate, err := strconv.ParseFloat(r.URL.Query().Get("index_rate"), 64)
	if err != nil {
		srv.Error("wrong request: %s", err)
		http.Error(w, fmt.Sprintf("wrong request: %s", err), http.StatusBadRequest)
	}

	mss, err := srv.marketsalaries(fid)
	if err != nil {
		onReadFail(srv, err, w)
	}
	if len(mss) == 0 {
		onEmptyTable(srv, w, dbmodels.MarketSalariesTable)
	}

	marksal := models.NewMarketSalaries(ind_rate)

	for _, ms := range mss {
		gr, err := strconv.Atoi(ms.Grade)
		if err != nil {
			srv.Error("wrong grade value: %v - %v", gr, err)
			http.Error(w, fmt.Sprintf("wrong grade value: %v - %v", gr, err), http.StatusBadRequest)
		}
		marksal.AddPayline(
			models.City(ms.City),
			models.Grade(gr),
			models.Salary(decToFloat(ms.Salary)),
			years,
		)
	}
	idx := srv.seqNextVal(dbmodels.MarketSalariesIndexedTable)
	stmt := srv.writeMarketSalaries()
	defer stmt.Close()

	for y, sg := range marksal.YearsGrid {
		for c, pl := range sg {
			for g, sal := range pl {
				_, err := stmt.Exec(
					idx,
					q.Get("function_id"),
					string(c),
					strconv.Itoa(int(g)),
					strconv.Itoa(y),
					floatToDec(float64(sal)),
					time.Now(),
				)
				idx++
				if err != nil {
					srv.Error("failed to prepare stmt: %v", err)
				}
			}
		}

	}
	srv.deleteFromTableByFunction(dbmodels.ProjectedHeadcountTable, fid)
	_, err = stmt.Exec()
	if err != nil {
		srv.Error("failed to bulk insert: %v", err)
	}
}

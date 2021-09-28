package handlers

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ts.tarusov/swp/internal/dbmodels"
)

func (srv *Service) writeProjectHeadcount() *sql.Stmt {
	stmt, err := srv.Db.PrepareContext(context.Background(),
		fmt.Sprintf(`BULK INSERT INTO "%s" VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, dbmodels.ProjectedHeadcountTable),
	)
	if err != nil {
		srv.Error("%s", err)
	}
	return stmt

}

func (srv *Service) writeBusinessValue() *sql.Stmt {
	stmt, err := srv.Db.PrepareContext(context.Background(),
		fmt.Sprintf(`BULK INSERT INTO "%s" VALUES (?, ?, ?, ?, ?, ?, ?)`, dbmodels.BusinessValueTable),
	)
	if err != nil {
		srv.Error("%s", err)
	}
	return stmt
}

func (srv *Service) writeProjectSkillLevels() *sql.Stmt {
	stmt, err := srv.Db.PrepareContext(context.Background(),
		fmt.Sprintf(`BULK INSERT INTO "%s" VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, dbmodels.ProjectSkillLevelsTable),
	)
	if err != nil {
		srv.Error("%s", err)
	}
	return stmt
}

func (srv *Service) writeProjectSufficiency() *sql.Stmt {
	stmt, err := srv.Db.PrepareContext(context.Background(),
		fmt.Sprintf(`BULK INSERT INTO "%s" VALUES (?, ?, ?, ?, ?, ?, ?)`, dbmodels.ProjectSufficiencyTable),
	)
	if err != nil {
		srv.Error("%s", err)
	}
	return stmt
}

func (srv *Service) writeAvailabilityRanks() *sql.Stmt {
	stmt, err := srv.Db.PrepareContext(context.Background(),
		fmt.Sprintf(`BULK INSERT INTO "%s" VALUES (?, ?, ?, ?, ?, ?)`, dbmodels.AvailabilityRanksTable),
	)
	if err != nil {
		srv.Error("%s", err)
	}
	return stmt
}
func (srv *Service) writePreferredStrategy() *sql.Stmt {
	stmt, err := srv.Db.PrepareContext(context.Background(),
		fmt.Sprintf(`BULK INSERT INTO "%s" VALUES (?, ?, ?, ?, ?, ?, ?)`, dbmodels.PreferredStrategyTable),
	)
	if err != nil {
		srv.Error("%s", err)
	}
	return stmt
}

func (srv *Service) writeMarketSalaries() *sql.Stmt {
	stmt, err := srv.Db.PrepareContext(context.Background(),
		fmt.Sprintf(`BULK INSERT INTO "%s" VALUES (?, ?, ?, ?, ?, ?, ?)`, dbmodels.MarketSalariesIndexedTable),
	)
	if err != nil {
		srv.Error("%s", err)
	}
	return stmt
}

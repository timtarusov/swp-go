package handlers

import (
	"fmt"
)

func (srv *Service) deleteFromTableById(table string, prid int) {
	srv.Db.Exec(fmt.Sprintf(`DELETE FROM "%s" WHERE "PROJECT_ID" = %d`,
		table, prid))
}
func (srv *Service) deleteFromTableByFunction(table string, fid int) {
	srv.Db.Exec(fmt.Sprintf(`DELETE FROM "%s" WHERE "FUNCTION_ID" = %d`,
		table, fid))
}

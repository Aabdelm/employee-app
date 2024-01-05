package employeedb

import "github.com/Aabdelm/employee-app/security"

func (DbMap *DbMap) FetchInfo(user string) (*security.Info, error) {
	var err error
	DbMap.l.Printf("[INFO] Starting function FetchInfo")

	stmt, err := DbMap.Db.Prepare("SELECT hashedpass,salt FROM authentication WHERE user = ?")
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to prepare statement in function FetchInfo. Error %s", err)
		return nil, err
	}

	info := security.NewInfoStruct()

	row := stmt.QueryRow(user)

	if err := row.Scan(&info.Pw, &info.Salt); err != nil {
		DbMap.l.Printf("[INFO] Failed to scan row. Error %s", err)
		return nil, err
	}

	DbMap.l.Printf("[INFO] successfully fetched user")
	return info, err
}

package mysql_db

import "time"

func (m *MySqlDB) GetLastSyncTimeFromDB() (time.Time, error) {
	var lastSyncTime time.Time
	err := m.db.QueryRow("SELECT last_sync_time FROM sync_time").Scan(&lastSyncTime)
	if err != nil {
		return time.Time{}, err
	}
	return lastSyncTime, nil
}

func (m *MySqlDB) UpdateLastSyncTimeInDB(newTime time.Time) error {
	_, err := m.db.Exec("UPDATE sync_time SET last_sync_time=?", newTime)
	return err
}

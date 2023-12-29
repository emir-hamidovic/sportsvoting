package mysql_db

import "time"

func (m *MySqlDB) GetLastSyncTimeFromDB() (time.Time, error) {
	var lastSyncTime int64
	err := m.db.QueryRow("SELECT last_sync_time FROM sync_time").Scan(&lastSyncTime)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(lastSyncTime, 0), nil
}

func (m *MySqlDB) UpdateLastSyncTimeInDB(newTime time.Time) error {
	_, err := m.db.Exec("UPDATE sync_time SET last_sync_time=?", newTime)
	return err
}

func (m *MySqlDB) GetGOATLastSyncTimeFromDB() (time.Time, error) {
	var lastSyncTime int64
	err := m.db.QueryRow("SELECT goat_last_sync_time FROM sync_time").Scan(&lastSyncTime)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(lastSyncTime, 0), nil
}

func (m *MySqlDB) UpdateGOATLastSyncTimeInDB(newTime time.Time) error {
	_, err := m.db.Exec("UPDATE sync_time SET goat_last_sync_time=?", newTime)
	return err
}

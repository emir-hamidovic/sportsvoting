package mysql_db

import "time"

func (m *MySqlDB) GetLastSyncTime(name string) (time.Time, error) {
	var lastSyncTime int64
	err := m.db.QueryRow("SELECT last_sync_time FROM sync_time WHERE name = ?", name).Scan(&lastSyncTime)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(lastSyncTime, 0), nil
}

func (m *MySqlDB) InsertLastSyncTime(newTime time.Time, name string) error {
	_, err := m.db.Exec("INSERT INTO sync_time(name, last_sync_time) VALUES (?, ?)", name, newTime)
	return err
}

func (m *MySqlDB) UpdateLastSyncTime(newTime time.Time, name string) error {
	_, err := m.db.Exec("UPDATE sync_time SET last_sync_time=? WHERE name = ?", newTime, name)
	return err
}

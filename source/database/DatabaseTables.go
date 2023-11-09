package database

func CreateUserTable() error {
	_, err := Instance.Exec(`
			create table if not exists users
			(
			    id       INTEGER
			        constraint users_pk
			            primary key autoincrement,
			    username TEXT,
			    password TEXT,
			    theme TEXT,
			    concurrents INTEGER,
			    cooldown INTEGER,
			    max_time INTEGER,
			    max_sessions INTEGER,
			    expiry INTEGER,
			    roles TEXT,
			    created_by TEXT
			);
    `)
	return err
}

func CreateLogsTable() error {
	_, err := Instance.Exec(`
			create table if not exists logs
			(
			    id       INTEGER
			        constraint logs_pk
			            primary key autoincrement,
			    user_id INTEGER,
			    target TEXT,
			    duration INTEGER,
			    port INTEGER,
			    method TEXT,
			    time_created INTEGER,
			    time_end INTEGER
			);
    `)
	return err
}

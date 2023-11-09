package database

import (
	"time"
)

type Flood struct {
	Id       int       `swash:"id"`
	UserId   int       `swash:"user_id"`
	Target   string    `swash:"target"`
	Duration int       `swash:"duration"`
	Port     int       `swash:"port"`
	Method   string    `swash:"method"`
	Created  time.Time `swash:"created"`
	End      time.Time `swash:"end"`
}

func LogAttack(flood *Flood, user *UserProfile) error {
	_, err := Instance.Exec("INSERT INTO logs (user_id, target, duration, port, method, time_created, time_end) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.Id,
		flood.Target,
		flood.Duration,
		flood.Port,
		flood.Method,
		flood.Created.Unix(),
		flood.End.Unix(),
	)
	return err
}

func FloodByID(id int) (*Flood, error) {
	rows, err := Instance.Query("SELECT * FROM logs WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		var flood = &Flood{}
		var created int64
		var end int64

		err := rows.Scan(&flood.Id,
			&flood.UserId,
			&flood.Target,
			&flood.Duration,
			&flood.Port,
			&flood.Method,
			&created,
			&end,
		)

		if err != nil {
			return nil, err
		}

		flood.Created = time.Unix(created, 0)
		flood.End = time.Unix(end, 0)

		return flood, nil
	}

	return nil, nil
}

func FloodsDuring(duration time.Duration) ([]*Flood, error) {
	var floods []*Flood

	endTime := time.Now()
	startTime := endTime.Add(-duration)

	rows, err := Instance.Query("SELECT id FROM logs WHERE time_created BETWEEN ? AND ?", startTime.Unix(), endTime.Unix())
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int

		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		flood, err := FloodByID(id)
		if err != nil {
			return nil, err
		}

		floods = append(floods, flood)
	}

	return floods, nil
}

func Floods() ([]*Flood, error) {
	var floods []*Flood

	rows, err := Instance.Query("SELECT id FROM logs")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int

		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		flood, err := FloodByID(id)
		if err != nil {
			return nil, err
		}

		floods = append(floods, flood)
	}

	return floods, nil
}

func LastFlood() (*Flood, error) {
	rows, err := Instance.Query("SELECT id FROM logs WHERE rowid = (SELECT MAX(rowid) FROM logs)")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		var id int

		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		flood, err := FloodByID(id)
		if err != nil {
			return nil, err
		}

		return flood, nil
	}

	return nil, nil
}

func RunningAttacks() ([]*Flood, error) {
	var runningFloods []*Flood

	floods, err := Floods()
	if err != nil {
		return nil, err
	}

	for _, flood := range floods {
		if !flood.End.After(time.Now()) {
			continue
		}

		runningFloods = append(runningFloods, flood)
	}

	return runningFloods, nil
}

func (flood *Flood) Remove() error {
	_, err := Instance.Exec("DELETE FROM logs WHERE id = ? AND user_id = ?", flood.Id, flood.UserId)
	return err
}

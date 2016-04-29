package geography

import (
	"database/sql"
	"fmt"

	"github.com/curt-labs/API/helpers/redis"
	"github.com/curt-labs/API/helpers/sortutil"
)

var (
	getAllStatesStmt             = `select st.stateID, st.state, st.abbr, st.countryID, c.name, c.abbr from States st join Country c on c.countryID = st.countryID`
	getAllCountriesStmt          = `select c.countryID, c.name, c.abbr from Country c order by c.countryID`
	getAllCountriesAndStatesStmt = `select c.countryID, c.name, c.abbr,
									s.stateID, s.state, s.abbr
									from Country c
									left join States s on c.countryID = s.countryID
									order by c.countryID, s.state`
)

type State struct {
	ID           int      `json:"-" xml:"-"`
	State        string   `json:"state"`
	Abbreviation string   `json:"abbreviation"`
	Country      *Country `json:"country"`
}

type Country struct {
	ID           int     `json:"-" xml:"-"`
	Country      string  `json:"country"`
	Abbreviation string  `json:"abbreviation"`
	States       []State `json:"states"`
}

func GetAllCountriesAndStates(db *sql.DB) (countries []Country, err error) {

	stmt, err := db.Prepare(getAllCountriesAndStatesStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil && err != sql.ErrNoRows {
		return
	}

	countryMap := make(map[int]Country, 0)

	for rows.Next() {
		var c Country
		var s State
		var stateID *int
		var state, abbr *string

		err = rows.Scan(
			&c.ID,
			&c.Country,
			&c.Abbreviation,
			&stateID,
			&state,
			&abbr,
		)
		if err != nil {
			continue
		}

		if stateID != nil {
			s.ID = *stateID
			if state != nil {
				s.State = *state
			}
			if abbr != nil {
				s.Abbreviation = *abbr
			}
		}

		if _, exists := countryMap[c.ID]; !exists {
			countryMap[c.ID] = c
		}

		tmp := countryMap[c.ID]
		tmp.States = append(tmp.States, s)
		countryMap[c.ID] = tmp
	}
	defer rows.Close()

	for _, c := range countryMap {
		countries = append(countries, c)
	}

	sortutil.AscByField(countries, "ID")
	return countries, nil
}

func GetAllCountries(db *sql.DB) (countries []Country, err error) {

	stmt, err := db.Prepare(getAllCountriesStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil && err != sql.ErrNoRows {
		return
	}

	for rows.Next() {
		var c Country
		err = rows.Scan(
			&c.ID,
			&c.Country,
			&c.Abbreviation,
		)
		if err != nil {
			continue
		}
		countries = append(countries, c)
	}
	defer rows.Close()

	sortutil.AscByField(countries, "ID")

	return countries, nil
}

func GetAllStates(db *sql.DB) (states []State, err error) {

	stmt, err := db.Prepare(getAllStatesStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	defer rows.Close()

	for rows.Next() {
		var state State
		state.Country = &Country{}
		err = rows.Scan(
			&state.ID,
			&state.State,
			&state.Abbreviation,
			&state.Country.ID,
			&state.Country.Country,
			&state.Country.Abbreviation,
		)
		if err != nil {
			continue
		}
		states = append(states, state)
	}

	return
}

func GetStateMap(db *sql.DB) (map[int]State, error) {
	stateMap := make(map[int]State)
	states, err := GetAllStates(db)
	for _, state := range states {
		stateMap[state.ID] = state
		redisKey := fmt.Sprintf("state:%d", state.ID)
		err = redis.Set(redisKey, state)
	}
	return stateMap, err
}

func GetCountryMap(db *sql.DB) (map[int]Country, error) {
	countryMap := make(map[int]Country)
	countries, err := GetAllCountries(db)
	for _, country := range countries {
		countryMap[country.ID] = country
		redisKey := fmt.Sprintf("country:%d", country.ID)
		err = redis.Set(redisKey, country)
	}
	return countryMap, err
}

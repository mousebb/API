package geography

import (
	"strconv"

	"github.com/curt-labs/API/helpers/redis"
	"github.com/curt-labs/API/helpers/sortutil"
	"github.com/curt-labs/API/middleware"
)

var (
	getAllStatesStmt             = `select st.stateID, st.state, st.abbr, st.countryID, c.name, c.abbr from States st join Country c on c.countryID = st.countryID`
	getAllCountriesStmt          = `select c.countryID, c.name, c.abbr from Country c`
	getAllCountriesAndStatesStmt = `select C.*, S.stateID, S.state, S.abbr from Country C
									inner join States S on S.countryID = C.countryID
									order by C.countryID, S.state`
)

type States []State
type State struct {
	Id           int      `json:"state_id"`
	State        string   `json:"state"`
	Abbreviation string   `json:"abbreviation"`
	Country      *Country `json:"country,omitempty"`
}

type Countries []Country
type Country struct {
	Id           int     `json:"country_id"`
	Country      string  `json:"country"`
	Abbreviation string  `json:"abbreviation"`
	States       *States `json:"states,omitempty"`
}

func GetAllCountriesAndStates(ctx *middleware.APIContext) (countries Countries, err error) {

	stmt, err := ctx.DB.Prepare(getAllCountriesAndStatesStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	countryMap := make(map[int]Country, 0)

	for rows.Next() {
		var c Country
		var s State

		err = rows.Scan(
			&c.Id,
			&c.Country,
			&c.Abbreviation,
			&s.Id,
			&s.State,
			&s.Abbreviation,
		)
		if err != nil {
			return
		}

		country, exists := countryMap[c.Id]

		if !exists {
			c.States = &States{s}
			countryMap[c.Id] = c
		} else {
			*country.States = append(*country.States, s)
		}
	}
	defer rows.Close()

	for _, c := range countryMap {
		countries = append(countries, c)
	}

	sortutil.AscByField(countries, "Id")
	return
}

func GetAllCountries(ctx *middleware.APIContext) (countries Countries, err error) {

	stmt, err := ctx.DB.Prepare(getAllCountriesStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var c Country
		err = rows.Scan(
			&c.Id,
			&c.Country,
			&c.Abbreviation,
		)
		if err != nil {
			return
		}
		countries = append(countries, c)
	}
	defer rows.Close()

	sortutil.AscByField(countries, "Id")

	return
}

func GetAllStates(ctx *middleware.APIContext) (states States, err error) {

	stmt, err := ctx.DB.Prepare(getAllStatesStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var state State
		state.Country = &Country{}
		err = rows.Scan(
			&state.Id,
			&state.State,
			&state.Abbreviation,
			&state.Country.Id,
			&state.Country.Country,
			&state.Country.Abbreviation,
		)
		if err != nil {
			return
		}
		states = append(states, state)
	}
	defer rows.Close()

	return
}

func GetStateMap(ctx *middleware.APIContext) (map[int]State, error) {
	stateMap := make(map[int]State)
	states, err := GetAllStates(ctx)
	for _, state := range states {
		stateMap[state.Id] = state
		redis_key := "state:" + strconv.Itoa(state.Id)
		err = redis.Set(redis_key, state)
	}
	return stateMap, err
}

func GetCountryMap(ctx *middleware.APIContext) (map[int]Country, error) {
	countryMap := make(map[int]Country)
	countries, err := GetAllCountries(ctx)
	for _, country := range countries {
		countryMap[country.Id] = country
		redis_key := "country:" + strconv.Itoa(country.Id)
		err = redis.Set(redis_key, country)
	}
	return countryMap, err
}

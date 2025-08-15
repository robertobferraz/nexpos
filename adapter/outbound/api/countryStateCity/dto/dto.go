package dto

type Country struct {
	ID        *int    `json:"id"`
	Name      *string `json:"name"`
	Iso2      *string `json:"iso2"`
	Iso3      *string `json:"iso3"`
	PhoneCode *string `json:"phonecode"`
	Capital   *string `json:"capital"`
	Currency  *string `json:"currency"`
	Native    *string `json:"native"`
	Emoji     *string `json:"emoji"`
}
type GetCountryResponse []Country

type State struct {
	Id          *int    `json:"id"`
	Name        *string `json:"name"`
	CountryId   *int    `json:"country_id"`
	CountryCode *string `json:"country_code"`
	Iso2        *string `json:"iso2"`
	Type        *string `json:"type"`
	Latitude    *string `json:"latitude"`
	Longitude   *string `json:"longitude"`
}

type GetStateResponse []State

type GetCity struct {
	Id               *int    `json:"id"`
	Name             *string `json:"name"`
	StateCode        *string `json:"state_code"`
	CountryCode      *string `json:"country_code"`
	CountryLocalCode *string `json:"country_local_code"`
}

type GetCityResponse []GetCity

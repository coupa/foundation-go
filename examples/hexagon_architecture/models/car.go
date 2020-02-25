package models

type Car struct {
	ID           int64  `db:"id" json:"id"`
	LicensePlate string `db:"license_place" json:"license_plate" valid:"stringlength(6|6)"`
	Make         string `db:"make" json:"make"`
	Year         int    `db:"year" json:"year" valid:"range(1980|9999)"`
	Crashed      bool   `db:"crashed" json:"crashed" valid:"type(bool)"`
}

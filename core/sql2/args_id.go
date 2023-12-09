package CoreSQL2

type ArgsID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

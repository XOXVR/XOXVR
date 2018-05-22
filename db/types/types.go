package types

type (
	Buyer struct {
		Id           int
		EthAddr      string   `sql:"type:varchar(100),unique"`
		BtcAddr      string   `sql:"type:varchar(100),unique"`
		Transactions []*Transaction
	}

	Transaction struct {
		Id      int
		BuyerId int
		Buyer   *Buyer
		From    string `sql:"type:varchar(100),notnull"`
		Amount  float64
		Hash    string `sql:"type:varchar(100),unique"`
	}

	NotHandledTransaction struct {
		Id     int
		From   string `sql:"type:varchar(100),notnull"`
		Amount float64
		Hash   string `sql:"type:varchar(100),unique"`
	}
)

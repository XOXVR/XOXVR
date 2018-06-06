package types

type (
	Buyer struct {
		Id           int
		EthAddr      string   `sql:",unique,type:varchar(100)"`
		BtcAddr      string   `sql:",unique,type:varchar(100)"`
		BtcAddrWif	 string   `sql:",unique,type:varchar(100)"`
		Transactions []*Transaction
	}

	Transaction struct {
		Id      int
		BuyerId int
		Buyer   *Buyer
		From    string `sql:",notnull,type:varchar(100)"`
		Amount  float64
		Hash    string `sql:",unique,type:varchar(100)"`
	}

	NotHandledTransaction struct {
		Id     int
		From   string `sql:",notnull,type:varchar(100)"`
		Amount float64
		Hash   string `sql:",unique,type:varchar(100)"`
	}
)

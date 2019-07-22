package client

type Response struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		Hash         string        `json:"hash"`
		Height       string        `json:"height"`
		Time         string        `json:"time"`
		NumTXS       string        `json:"num_txs"`
		TotalTXS     string        `json:"total_txs"`
		Transactions []Transaction `json:"transactions"`
		BlockReward  string        `json:"block_reward"`
		Size         string        `json:"size"`
		Proposer     string        `json:"proposer"`
		Validators   []struct {
			PubKey string `json:"pub_key"`
			Signed bool   `json:"signed"`
		} `json:"validators"`
	} `json:"result"`
	Error struct {
		Code    uint64 `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	} `json:"error,omitempty"`
}

type Transaction struct {
	Hash     string `json:"hash"`
	RawTX    string `json:"raw_tx"`
	From     string `json:"from"`
	Nonce    string `json:"nonce"`
	GasPrice uint64 `json:"gas_price"`
	Type     uint8  `json:"type"`
	Data     struct {
		Coin  string `json:"coin"`
		To    string `json:"to"`
		Value string `json:"value"`
	} `json:"data"`
	Payload     string `json:"payload"`
	ServiceData string `json:"service_data"`
	Gas         string `json:"gas"`
	GasCoin     string `json:"gas_coin"`
	Tags        struct {
		TxFrom string `json:"tx.from"`
		TxTo   string `json:"tx.to"`
		TxCoin string `json:"tx.coin"`
		TxType string `json:"tx.type"`
	} `json:"tags"`
}

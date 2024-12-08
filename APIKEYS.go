package Binance

type APIKEYS struct {
	KEY    string
	SECRET string
}

func (keys *APIKEYS) Set(APIKEY string, APISECRET string) {
	keys.KEY = APIKEY
	keys.SECRET = APISECRET
}

func (keys *APIKEYS) Get() (KEY string, SECRET string) {
	return keys.KEY, keys.SECRET
}

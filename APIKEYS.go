package Binance

type APIKEYS struct {
	KEY    string
	SECRET string
}

func (keys *APIKEYS) Set(KEY string, SECRET string) {
	keys.KEY = KEY
	keys.SECRET = SECRET
}

func (keys *APIKEYS) Get() (KEY string, SECRET string) {
	return keys.KEY, keys.SECRET
}

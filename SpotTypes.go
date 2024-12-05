package main

var SPOT_Constants = struct {
	URLs                [6]string
	URL_Data_Only       string
	SecurityTypes       Spot_SecurityTypes_ENUM
	ExchangeFilterTypes Spot_Exchange_FilterTypes_ENUM
	SymbolFilterTypes   SPOT_Symbol_FilterTypes_ENUM
	SymbolStatuses      Spot_SymbolStatuses_ENUM
	Permissions         Spot_Permissions_ENUM
	OrderStatuses       Spot_OrderStatuses_ENUM
	ListStatusTypes     Spot_ListStatusTypes_ENUM
	ListOrderStatuses   Spot_ListOrderStatuses_ENUM
	ContingencyTypes    Spot_ContingencyTypes_ENUM
	AllocationTypes     Spot_AllocationTypes_ENUM
	OrderTypes          Spot_OrderTypes_ENUM
	NewOrderRespTypes   Spot_NewOrderRespTypes_ENUM
	WorkingFloors       Spot_WorkingFloors_ENUM
	OrderSides          Spot_OrderSides_ENUM
	TimeInForces        Spot_TimeInForces_ENUM
	RateLimitTypes      Spot_RateLimitTypes_ENUM
	RateLimitIntervals  Spot_RateLimitIntervals_ENUM
	STPModes            Spot_STPModes_ENUM
}{
	URLs:          [6]string{"https://api.binance.com", "https://api-gcp.binance.com", "https://api1.binance.com", "https://api2.binance.com", "https://api3.binance.com", "https://api4.binance.com"},
	URL_Data_Only: "https://data-api.binance.vision",
	SecurityTypes: Spot_SecurityTypes_ENUM{
		NONE:        "NONE",
		USER_STREAM: "USER_STREAM",
		TRADE:       "TRADE",
		USER_DATA:   "USER_DATA",
	},
	ExchangeFilterTypes: Spot_Exchange_FilterTypes_ENUM{
		EXCHANGE_MAX_NUM_ORDERS:         "EXCHANGE_MAX_NUM_ORDERS",
		EXCHANGE_MAX_NUM_ALGO_ORDERS:    "EXCHANGE_MAX_NUM_ALGO_ORDERS",
		EXCHANGE_MAX_NUM_ICEBERG_ORDERS: "EXCHANGE_MAX_NUM_ICEBERG_ORDERS",
	},
	SymbolFilterTypes: SPOT_Symbol_FilterTypes_ENUM{
		PRICE_FILTER:           "PRICE_FILTER",
		PERCENT_PRICE:          "PERCENT_PRICE",
		PERCENT_PRICE_BY_SIDE:  "PERCENT_PRICE_BY_SIDE",
		LOT_SIZE:               "LOT_SIZE",
		MIN_NOTIONAL:           "MIN_NOTIONAL",
		NOTIONAL:               "NOTIONAL",
		ICEBERG_PARTS:          "ICEBERG_PARTS",
		MARKET_LOT_SIZE:        "MARKET_LOT_SIZE",
		MAX_NUM_ORDERS:         "MAX_NUM_ORDERS",
		MAX_NUM_ALGO_ORDERS:    "MAX_NUM_ALGO_ORDERS",
		MAX_NUM_ICEBERG_ORDERS: "MAX_NUM_ICEBERG_ORDERS",
		MAX_POSITION:           "MAX_POSITION",
		TRAILING_DELTA:         "TRAILING_DELTA",
	},
	SymbolStatuses: Spot_SymbolStatuses_ENUM{
		PRE_TRADING:   "PRE_TRADING",
		TRADING:       "TRADING",
		POST_TRADING:  "POST_TRADING",
		END_OF_DAY:    "END_OF_DAY",
		HALT:          "HALT",
		AUCTION_MATCH: "AUCTION_MATCH",
		BREAK:         "BREAK",
	},
	Permissions: Spot_Permissions_ENUM{
		SPOT:        "SPOT",
		MARGIN:      "MARGIN",
		LEVERAGED:   "LEVERAGED",
		TRD_GRP_002: "TRD_GRP_002",
		TRD_GRP_003: "TRD_GRP_003",
		TRD_GRP_004: "TRD_GRP_004",
		TRD_GRP_005: "TRD_GRP_005",
		TRD_GRP_006: "TRD_GRP_006",
		TRD_GRP_007: "TRD_GRP_007",
		TRD_GRP_008: "TRD_GRP_008",
		TRD_GRP_009: "TRD_GRP_009",
		TRD_GRP_010: "TRD_GRP_010",
		TRD_GRP_011: "TRD_GRP_011",
		TRD_GRP_012: "TRD_GRP_012",
		TRD_GRP_013: "TRD_GRP_013",
		TRD_GRP_014: "TRD_GRP_014",
		TRD_GRP_015: "TRD_GRP_015",
		TRD_GRP_016: "TRD_GRP_016",
		TRD_GRP_017: "TRD_GRP_017",
		TRD_GRP_018: "TRD_GRP_018",
		TRD_GRP_019: "TRD_GRP_019",
		TRD_GRP_020: "TRD_GRP_020",
		TRD_GRP_021: "TRD_GRP_021",
		TRD_GRP_022: "TRD_GRP_022",
		TRD_GRP_023: "TRD_GRP_023",
		TRD_GRP_024: "TRD_GRP_024",
		TRD_GRP_025: "TRD_GRP_025",
	},
	OrderStatuses: Spot_OrderStatuses_ENUM{
		NEW:              "NEW",
		PENDING_NEW:      "PENDING_NEW",
		PARTIALLY_FILLED: "PARTIALLY_FILLED",
		FILLED:           "FILLED",
		CANCELED:         "CANCELED",
		PENDING_CANCEL:   "PENDING_CANCEL",
		REJECTED:         "REJECTED",
		EXPIRED:          "EXPIRED",
		EXPIRED_IN_MATCH: "EXPIRED_IN_MATCH",
	},
	ListStatusTypes: Spot_ListStatusTypes_ENUM{
		RESPONSE:     "RESPONSE",
		EXEC_STARTED: "EXEC_STARTED",
		ALL_DONE:     "ALL_DONE",
	},
	ListOrderStatuses: Spot_ListOrderStatuses_ENUM{
		EXECUTING: "EXECUTING",
		ALL_DONE:  "ALL_DONE",
		REJECT:    "REJECT",
	},
	ContingencyTypes: Spot_ContingencyTypes_ENUM{
		OCO: "OCO",
		OTO: "OTO",
	},
	AllocationTypes: Spot_AllocationTypes_ENUM{
		SOR: "SOR",
	},
	OrderTypes: Spot_OrderTypes_ENUM{
		LIMIT:             "LIMIT",
		MARKET:            "MARKET",
		STOP_LOSS:         "STOP_LOSS",
		STOP_LOSS_LIMIT:   "STOP_LOSS_LIMIT",
		TAKE_PROFIT:       "TAKE_PROFIT",
		TAKE_PROFIT_LIMIT: "TAKE_PROFIT_LIMIT",
		LIMIT_MAKER:       "LIMIT_MAKER",
	},
	NewOrderRespTypes: Spot_NewOrderRespTypes_ENUM{
		ACK:    "ACK",
		RESULT: "RESULT",
		FULL:   "FULL",
	},
	WorkingFloors: Spot_WorkingFloors_ENUM{
		EXCHANGE: "EXCHANGE",
		SOR:      "SOR",
	},
	OrderSides: Spot_OrderSides_ENUM{
		BUY:  "BUY",
		SELL: "SELL",
	},
	TimeInForces: Spot_TimeInForces_ENUM{
		GTC: "GTC",
		IOC: "IOC",
		FOK: "FOK",
	},
	RateLimitTypes: Spot_RateLimitTypes_ENUM{
		REQUEST_WEIGHT: "REQUEST_WEIGHT",
		ORDERS:         "ORDERS",
		RAW_REQUESTS:   "RAW_REQUESTS",
	},
	RateLimitIntervals: Spot_RateLimitIntervals_ENUM{
		SECOND: "SECOND",
		MINUTE: "MINUTE",
		DAY:    "DAY",
	},
	STPModes: Spot_STPModes_ENUM{
		NONE:         "NONE",
		EXPIRE_MAKER: "EXPIRE_MAKER",
		EXPIRE_TAKER: "EXPIRE_TAKER",
		EXPIRE_BOTH:  "EXPIRE_BOTH",
	},
}

type Spot_SecurityTypes_ENUM struct {
	NONE        string
	USER_STREAM string
	TRADE       string
	USER_DATA   string
}

type Spot_Exchange_FilterTypes_ENUM struct {
	EXCHANGE_MAX_NUM_ORDERS         string
	EXCHANGE_MAX_NUM_ALGO_ORDERS    string
	EXCHANGE_MAX_NUM_ICEBERG_ORDERS string
}

type SPOT_Symbol_FilterTypes_ENUM struct {
	PRICE_FILTER           string
	PERCENT_PRICE          string
	PERCENT_PRICE_BY_SIDE  string
	LOT_SIZE               string
	MIN_NOTIONAL           string
	NOTIONAL               string
	ICEBERG_PARTS          string
	MARKET_LOT_SIZE        string
	MAX_NUM_ORDERS         string
	MAX_NUM_ALGO_ORDERS    string
	MAX_NUM_ICEBERG_ORDERS string
	MAX_POSITION           string
	TRAILING_DELTA         string
}

type Spot_SymbolStatuses_ENUM struct {
	PRE_TRADING   string
	TRADING       string
	POST_TRADING  string
	END_OF_DAY    string
	HALT          string
	AUCTION_MATCH string
	BREAK         string
}

type Spot_Permissions_ENUM struct {
	SPOT        string
	MARGIN      string
	LEVERAGED   string
	TRD_GRP_002 string
	TRD_GRP_003 string
	TRD_GRP_004 string
	TRD_GRP_005 string
	TRD_GRP_006 string
	TRD_GRP_007 string
	TRD_GRP_008 string
	TRD_GRP_009 string
	TRD_GRP_010 string
	TRD_GRP_011 string
	TRD_GRP_012 string
	TRD_GRP_013 string
	TRD_GRP_014 string
	TRD_GRP_015 string
	TRD_GRP_016 string
	TRD_GRP_017 string
	TRD_GRP_018 string
	TRD_GRP_019 string
	TRD_GRP_020 string
	TRD_GRP_021 string
	TRD_GRP_022 string
	TRD_GRP_023 string
	TRD_GRP_024 string
	TRD_GRP_025 string
}

type Spot_OrderStatuses_ENUM struct {
	NEW              string
	PENDING_NEW      string
	PARTIALLY_FILLED string
	FILLED           string
	CANCELED         string
	PENDING_CANCEL   string
	REJECTED         string
	EXPIRED          string
	EXPIRED_IN_MATCH string
}

type Spot_ListStatusTypes_ENUM struct {
	RESPONSE     string
	EXEC_STARTED string
	ALL_DONE     string
}

type Spot_ListOrderStatuses_ENUM struct {
	EXECUTING string
	ALL_DONE  string
	REJECT    string
}

type Spot_ContingencyTypes_ENUM struct {
	OCO string
	OTO string
}

type Spot_AllocationTypes_ENUM struct {
	SOR string
}

type Spot_OrderTypes_ENUM struct {
	LIMIT             string
	MARKET            string
	STOP_LOSS         string
	STOP_LOSS_LIMIT   string
	TAKE_PROFIT       string
	TAKE_PROFIT_LIMIT string
	LIMIT_MAKER       string
}

type Spot_NewOrderRespTypes_ENUM struct {
	ACK    string
	RESULT string
	FULL   string
}

type Spot_WorkingFloors_ENUM struct {
	EXCHANGE string
	SOR      string
}

type Spot_OrderSides_ENUM struct {
	BUY  string
	SELL string
}

type Spot_TimeInForces_ENUM struct {
	GTC string
	IOC string
	FOK string
}

type Spot_RateLimitTypes_ENUM struct {
	REQUEST_WEIGHT string
	ORDERS         string
	RAW_REQUESTS   string
}

type Spot_RateLimitIntervals_ENUM struct {
	SECOND string
	MINUTE string
	DAY    string
}

type Spot_STPModes_ENUM struct {
	NONE         string
	EXPIRE_MAKER string
	EXPIRE_TAKER string
	EXPIRE_BOTH  string
}

//

type Spot_RateLimitType struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

type Spot_ExchangeFilters struct {
	EXCHANGE_MAX_NUM_ORDERS         *Spot_ExchangeFilter_EXCHANGE_MAX_NUM_ORDERS
	EXCHANGE_MAX_NUM_ALGO_ORDERS    *Spot_ExchangeFilter_EXCHANGE_MAX_NUM_ALGO_ORDERS
	EXCHANGE_MAX_NUM_ICEBERG_ORDERS *Spot_ExchangeFilter_EXCHANGE_MAX_NUM_ICEBERG_ORDERS
}

type Spot_ExchangeFilter_EXCHANGE_MAX_NUM_ORDERS struct {
	FilterType   string `json:"filterType"`
	MaxNumOrders int64  `json:"maxNumOrders"`
}

type Spot_ExchangeFilter_EXCHANGE_MAX_NUM_ALGO_ORDERS struct {
	FilterType       string `json:"filterType"`
	MaxNumAlgoOrders int64  `json:"maxNumAlgoOrders"`
}

type Spot_ExchangeFilter_EXCHANGE_MAX_NUM_ICEBERG_ORDERS struct {
	FilterType          string `json:"filterType"`
	MaxNumIcebergOrders int64  `json:"maxNumIcebergOrders"`
}

//

type Spot_Symbol struct {
	Symbol                          string   `json:"symbol"`
	Status                          string   `json:"status"`
	BaseAsset                       string   `json:"baseAsset"`
	BaseAssetPrecision              int      `json:"baseAssetPrecision"`
	QuoteAsset                      string   `json:"quoteAsset"`
	QuotePrecision                  int      `json:"quotePrecision"`
	BaseCommissionPrecision         int      `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision        int      `json:"quoteCommissionPrecision"`
	OrderTypes                      []string `json:"orderTypes"`
	IcebergAllowed                  bool     `json:"icebergAllowed"`
	OcoAllowed                      bool     `json:"ocoAllowed"`
	OtoAllowed                      bool     `json:"otoAllowed"`
	QuoteOrderQtyMarketAllowed      bool     `json:"quoteOrderQtyMarketAllowed"`
	AllowTrailingStop               bool     `json:"allowTrailingStop"`
	CancelReplaceAllowedbool        bool     `json:"cancelReplaceAllowedbool"`
	IsSpotTradingAllowed            bool     `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed          bool     `json:"isMarginTradingAllowed"`
	Filters                         Spot_SymbolFilters
	Permissions                     []string   `json:"permissions"`
	PermissionSets                  [][]string `json:"permissionSets"`
	DefaultSelfTradePreventionMode  string     `json:"defaultSelfTradePreventionMode"`
	AllowedSelfTradePreventionModes []string   `json:"allowedSelfTradePreventionModes"`
}

type Spot_SymbolFilters struct {
	PRICE_FILTER           *Spot_SymbolFilter_PRICE_FILTER
	PERCENT_PRICE          *Spot_SymbolFilter_PERCENT_PRICE
	PERCENT_PRICE_BY_SIDE  *Spot_SymbolFilter_PERCENT_PRICE_BY_SIDE
	LOT_SIZE               *Spot_SymbolFilter_LOT_SIZE
	MIN_NOTIONAL           *Spot_SymbolFilter_MIN_NOTIONAL
	NOTIONAL               *Spot_SymbolFilter_NOTIONAL
	ICEBERG_PARTS          *Spot_SymbolFilter_ICEBERG_PARTS
	MARKET_LOT_SIZE        *Spot_SymbolFilter_MARKET_LOT_SIZE
	MAX_NUM_ORDERS         *Spot_SymbolFilter_MAX_NUM_ORDERS
	MAX_NUM_ALGO_ORDERS    *Spot_SymbolFilter_MAX_NUM_ALGO_ORDERS
	MAX_NUM_ICEBERG_ORDERS *Spot_SymbolFilter_MAX_NUM_ICEBERG_ORDERS
	MAX_POSITION           *Spot_SymbolFilter_MAX_POSITION
	TRAILING_DELTA         *Spot_SymbolFilter_TRAILING_DELTA
}

type Spot_ExchangeInfo_SORS struct {
	BaseAsset string   `json:"baseAsset"`
	Symbols   []string `json:"symbols"`
}

type Spot_SymbolFilter_PRICE_FILTER struct {
	FilterType string `json:"filterType"`
	MinPrice   string `json:"minPrice"`
	MaxPrice   string `json:"maxPrice"`
	TickSize   string `json:"tickSize"`
}

type Spot_SymbolFilter_PERCENT_PRICE struct {
	FilterType     string `json:"filterType"`
	MultiplierUp   string `json:"multiplierUp"`
	MultiplierDown string `json:"multiplierDown"`
	AvgPriceMins   int64  `json:"avgPriceMins"`
}

type Spot_SymbolFilter_PERCENT_PRICE_BY_SIDE struct {
	FilterType        string `json:"filterType"`
	BidMultiplierUp   string `json:"bidMultiplierUp"`
	BidMultiplierDown string `json:"bidMultiplierDown"`
	AskMultiplierUp   string `json:"askMultiplierUp"`
	AskMultiplierDown string `json:"askMultiplierDown"`
	AvgPriceMins      int64  `json:"avgPriceMins"`
}

type Spot_SymbolFilter_LOT_SIZE struct {
	FilterType string `json:"filterType"`
	MinQty     string `json:"minQty"`
	MaxQty     string `json:"maxQty"`
	StepSize   string `json:"stepSize"`
}

type Spot_SymbolFilter_MIN_NOTIONAL struct {
	FilterType    string `json:"filterType"`
	MinNotional   string `json:"minNotional"`
	ApplyToMarket bool   `json:"applyToMarket"`
	AvgPriceMins  int64  `json:"avgPriceMins"`
}

type Spot_SymbolFilter_NOTIONAL struct {
	FilterType       string `json:"filterType"`
	MinNotional      string `json:"minNotional"`
	ApplyMinToMarket bool   `json:"applyMinToMarket"`
	MaxNotional      string `json:"maxNotional"`
	ApplyMaxToMarket bool   `json:"applyMaxToMarket"`
	AvgPriceMins     int64  `json:"avgPriceMins"`
}

type Spot_SymbolFilter_ICEBERG_PARTS struct {
	FilterType string `json:"filterType"`
	Limit      int64  `json:"limit"`
}

type Spot_SymbolFilter_MARKET_LOT_SIZE struct {
	FilterType string `json:"filterType"`
	MinQty     string `json:"minQty"`
	MaxQty     string `json:"maxQty"`
	StepSize   string `json:"stepSize"`
}

type Spot_SymbolFilter_MAX_NUM_ORDERS struct {
	FilterType   string `json:"filterType"`
	MaxNumOrders int64  `json:"maxNumOrders"`
}

type Spot_SymbolFilter_MAX_NUM_ALGO_ORDERS struct {
	FilterType       string `json:"filterType"`
	MaxNumAlgoOrders int64  `json:"maxNumAlgoOrders"`
}

type Spot_SymbolFilter_MAX_NUM_ICEBERG_ORDERS struct {
	FilterType          string `json:"filterType"`
	MaxNumIcebergOrders int64  `json:"maxNumIcebergOrders"`
}

type Spot_SymbolFilter_MAX_POSITION struct {
	FilterType  string `json:"filterType"`
	MaxPosition string `json:"maxPosition"`
}

type Spot_SymbolFilter_TRAILING_DELTA struct {
	FilterType            string `json:"filterType"`
	MinTrailingAboveDelta int64  `json:"minTrailingAboveDelta"`
	MaxTrailingAboveDelta int64  `json:"maxTrailingAboveDelta"`
	MinTrailingBelowDelta int64  `json:"minTrailingBelowDelta"`
	MaxTrailingBelowDelta int64  `json:"maxTrailingBelowDelta"`
}

//

type Spot_Time struct {
	ServerTime int64 `json:"serverTime"`
	Latency    int64
}

type Spot_ExchangeInfo struct {
	Timezone        string                `json:"timezone"`
	ServerTime      int64                 `json:"serverTime"`
	RateLimits      []*Spot_RateLimitType `json:"rateLimits"`
	ExchangeFilters Spot_ExchangeFilters
	Symbols_arr     []*Spot_Symbol `json:"symbols"`
	Symbols         map[string]*Spot_Symbol
	Sors            []*Spot_ExchangeInfo_SORS `json:"sors"`
}

type Spot_OrderBook struct {
	LastUpdateId int64 `json:"lastUpdateId"`
	//"bids": [
	//    [
	//      "4.00000000",     // PRICE
	//      "431.00000000"    // QTY
	//    ]
	//  ]
	Bids [][2]float64 `json:"bids"`
	// 	"asks": [
	//     [
	//       "4.00000200",
	//       "12.00000000"
	//     ]
	//   ]
	Asks [][2]float64 `json:"asks"`
}

type Spot_Trade struct {
	Id           int64  `json:"id"`
	Price        string `json:"price"`
	Qty          string `json:"qty"`
	QuoteQty     string `json:"quoteQty"`
	Time         int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
	IsBestMatch  bool   `json:"isBestMatch"`
}

type Spot_AggTrade struct {
	// Aggregate tradeId
	AggTradeId int64 `json:"a"`
	// Price
	Price string `json:"p"`
	// Quantity
	Quantity string `json:"q"`
	// First tradeId
	FirstTradeId int64 `json:"f"`
	// Last tradeId
	LastTradeId int64 `json:"l"`
	// Timestamp
	Timestamp int64 `json:"T"`
	// Was the buyer the maker?
	IsMaker bool `json:"m"`
	// Was the trade the best price match?
	IsBestMatch bool `json:"M"`
}

type Spot_Candlestick struct {
	// Kline open time
	OpenTime int64 `json:"openTime"`
	// Open price
	OpenPrice string `json:"openPrice"`
	// High price
	HighPrice string `json:"highPrice"`
	// Low price
	LowPrice string `json:"lowPrice"`
	// Close price
	ClosePrice string `json:"closePrice"`
	// Volume
	Volume string `json:"volume"`
	// Kline Close time
	CloseTime int64 `json:"closeTime"`
	// Quote asset volume
	QuoteAssetVolume string `json:"quoteAssetVolume"`
	// Number of trades
	TradeCount int64 `json:"tradeCount"`
	// Taker buy base asset volume
	TakerBuyBaseAssetVolume string `json:"TakerBuyBaseAssetVolume"`
	// Taker buy quote asset volume
	TakerBuyQuoteAssetVolume string `json:"TakerBuyQuoteAssetVolume"`
	// Unused field, ignore.
	Unused string `json:"unused"`
}

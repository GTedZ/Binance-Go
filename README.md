# Binance-Go
 A simple Golang package for the Binance API (No docs until v1.0.0)

## TODO v1.0.0

### SPOT
 - Finish the implementation of all REST and WS endpoints
 - Websocket API

### FUTURES
 - Websocket API

### Helper Functions
 - A function to fetch all candles between X and Y timestamps
 - A function to fetch all candles between X and Y timestamps, parse them to float64
 - A function to fetch all candles between X and Y timestamps, parse them to float64, serialize them and return []byte
 - A function to fetch the serialized data, store them, and the ability to pick up where it left off if the function stopped
 - A function to read and deserialize the data, along with helper functions (fetching a specific timestamp off of a big interval)
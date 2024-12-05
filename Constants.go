package main

var Constants = struct {
	Methods Methods
}{
	Methods: Methods{
		GET:    "GET",
		POST:   "POST",
		PUT:    "PUT",
		PATCH:  "PATCH",
		DELETE: "DELETE",
	},
}

type Methods struct {
	GET    string
	POST   string
	PUT    string
	PATCH  string
	DELETE string
}

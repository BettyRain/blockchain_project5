package p3

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Canonical",
		"GET",
		"/canonical",
		Canonical,
	},
	Route{
		"Show",
		"GET",
		"/show",
		Show,
	},
	Route{
		"Upload",
		"GET",
		"/upload",
		Upload,
	},
	Route{
		"UploadBlock",
		"GET",
		"/block/{height}/{hash}",
		UploadBlock,
	},
	Route{
		"HeartBeatReceive",
		"POST",
		"/heartbeat/receive",
		HeartBeatReceive,
	},
	Route{
		"DataReceive",
		"POST",
		"/data/receive",
		DataReceive,
	},
	Route{
		"Start",
		"GET",
		"/start",
		Start,
	},
	Route{
		"ShowTrie",
		"GET",
		"/trie",
		ShowTrie,
	},
}

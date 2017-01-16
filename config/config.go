package config

var BindAddr = ":20000"

// The URL of the babbage instance to use.
var BabbageURL = "http://localhost:8080"

// The URL of the content resolver.
var ResolverURL = "http://localhost:20020"

// The URL of the content renderer
var RendererURL = "http://localhost:20010"

// The URL of the Data Discovery frontend controller
var DataDiscoveryURL = "http://localhost:20030"

// The percentage of requests to send to the new homepage
var HomepageABPercent = 0

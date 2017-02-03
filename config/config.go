package config

var BindAddr = ":20000"

// The URL of the babbage instance to use.
var BabbageURL = "http://localhost:8080"

// The URL of the content resolver.
var ResolverURL = "http://localhost:20020"

// The URL of the content renderer
var RendererURL = "http://localhost:20010"

// The percentage of requests to send to the new homepage
var HomepageABPercent = 0

// Whether the template rendering engine is in development mode or not
var DebugMode = false

// The CDN assets path
var PatternLibraryAssetsPath = "https://cdn.ons.gov.uk/sixteens/6cc1837"

// The site domain
var SiteDomain = "ons.gov.uk"

// Splash page
var SplashPage = ""

// Redirect secret
var RedirectSecret = "secret"

package config

var BindAddr = ":20000"

// The URL of the babbage instance to use.
var BabbageURL = "http://localhost:8080"

// The URL of the content renderer
var RendererURL = "http://localhost:20010"

// Dataset routes enabled
var DatasetRoutesEnabled = false

// The URL of the dataset controller
var DatasetControllerURL = "http://localhost:20200"

// The URL of the filter dataset controller
var FilterDatasetControllerURL = "http://localhost:20001"

// The URL of the cookies controller
var CookiesControllerURL = "http://localhost:23800"

// The URL of the Geography controller
var GeographyControllerURL = "http://localhost:23700"

// Geography feature is enabled
var GeographyEnabled = false

// The URL of Zebedee API
var ZebedeeURL = "http://localhost:8082"

// The URL of the file downloader service
var DownloaderURL = "http://localhost:23400"

// The CDN assets path
var PatternLibraryAssetsPath = "https://cdn.ons.gov.uk/sixteens/f816ac8"

// The site domain
var SiteDomain = "ons.gov.uk"

// Redirect secret
var RedirectSecret = "secret"

// SQS URL for analytics data
var SQSAnalyticsURL = ""

// ContentTypeByteLimit respresents the response size at which we stop checking content-type to avoid oom errors
var ContentTypeByteLimit = 5000000 // 5mb

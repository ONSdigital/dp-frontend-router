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

// The URL of the Geography controller
var GeographyControllerURL = "http://localhost:23700"

// Geography feature is enabled
var GeographyEnabled = false

// The URL of Zebedee API
var ZebedeeURL = "http://localhost:8082"

// The URL of the file downloader service
var DownloaderURL = "http://localhost:23400"

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

// Disabled page
var DisabledPage = ""

// TaxonomyDomain is link to website. Used for CMD beta so global links go to website rather than beta domain
var TaxonomyDomain = ""

// SQS URL for analytics data
var SQSAnalyticsURL = ""

// ContentTypeByteLimit respresents the response size at which we stop checking content-type to avoid oom errors
var ContentTypeByteLimit = 5000000 // 5mb

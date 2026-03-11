
 = "C:\Users\Administrator\Desktop\Testing\aegisgate\cmd\aegisgate\main.go"
 = Get-Content -Path  -Raw

# Add flag and runtime imports
 =  -replace '("strconv")', '
	"flag"'
 =  -replace '("syscall")', '
	"runtime"'

# Add version constant after imports close
 = @'

// Version is the current AegisGate version
const version = "0.28.2"

// Build information - set via ldflags during build
var (
	buildDate = "unknown"
	gitCommit = "unknown"
	gitBranch = "unknown"
)
'@
 =  -replace '(
\// Security imports
)',  + ''

# Add CLI flags after "var i18nManager"
 = @'

// CLI flags
var (
	showVersion = flag.Bool("version", false, "show version information")
	showHelp    = flag.Bool("help", false, "show this help message")
	configPath  = flag.String("config", "", "path to configuration file")
)
'@
 =  -replace '(
var i18nManager \*i18n\.Manager)',  + ''

# Add early flag parsing at start of main()
 = @'

	// Parse CLI flags first
	flag.Parse()

	// Handle version flag
	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	// Handle help flag
	if *showHelp {
		printHelp()
		os.Exit(0)
	}

'@
 =  -replace '(
\t\/\/ Initialize structured logging
\tslog\.SetDefault)',  + ''

# Add helper functions at the end
 = @'


// printVersion prints version information
func printVersion() {
	println("AegisGate Security Gateway")
	println("Version:", version)
	println("Go Version:", runtime.Version())
	println("Build Date:", buildDate)
	println("Git Commit:", gitCommit)
	println("Git Branch:", gitBranch)
}

// printHelp prints help information
func printHelp() {
	println("AegisGate - Enterprise AI/LLM Security Gateway")
	println()
	println("Usage:")
	println("  aegisgate [options]")
	println()
	println("Options:")
	flag.PrintDefaults()
	println()
	println("Environment Variables:")
	println("  AEGISGATE_CONFIG_PATH      Path to configuration file")
	println("  AEGISGATE_LOG_LEVEL        Log level (debug, info, warn, error)")
	println("  AEGISGATE_LOCALE           Locale for i18n (en, es, fr, de, etc.)")
	println("  AEGISGATE_MITM_ENABLED     Enable MITM proxy mode")
	println("  AEGISGATE_MITM_PORT        MITM proxy port (default: 3128)")
}
'@
 =  + 

Set-Content -Path  -Value  -NoNewline
Write-Output "main.go updated with CLI flags"


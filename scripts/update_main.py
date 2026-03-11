
import re
import sys

main_go_path = r"C:\Users\Administrator\Desktop\Testing\padlock\cmd\padlock\main.go"

try:
    with open(main_go_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    print("Original file loaded, length:", len(content))
    
    # 1. Add flag and runtime to imports
    content = content.replace('"strconv"', '"strconv"\n\t"flag"')
    content = content.replace('"syscall"', '"syscall"\n\t"runtime"')
    print("Imports added")
    
    # 2. Add version constant and CLI flags after imports
    version_block = """\
// Version is the current Padlock version
const version = "0.28.2"

// Build information - set via ldflags during build
var (
	buildDate = "unknown"
	gitCommit = "unknown"
	gitBranch = "unknown"
)

// CLI flags
var (
	showVersion = flag.Bool("version", false, "show version information")
	showHelp    = flag.Bool("help", false, "show this help message")
	configPath  = flag.String("config", "", "path to configuration file")
)
"""
    
    content = content.replace('// AuthAccessControl contains', version_block + '
// AuthAccessControl contains')
    print("Version block added")
    
    # 3. Add flag parsing at start of main()
    flag_parsing = """\
\t// Parse CLI flags first
\tflag.Parse()

\t// Handle version flag
\tif *showVersion {
\t\tprintVersion()
\t\tos.Exit(0)
\t}

\t// Handle help flag
\tif *showHelp {
\t\tprintHelp()
\t\tos.Exit(0)
\t}

"""
    
    content = content.replace('\t// Initialize structured logging\n\t', flag_parsing + '\t// Initialize structured logging\n\t')
    print("Flag parsing added")
    
    # 4. Add helper functions at end
    helper_funcs = """\

// printVersion prints version information
func printVersion() {
\tprintln("Padlock Security Gateway")
\tprintln("Version:", version)
\tprintln("Go Version:", runtime.Version())
\tprintln("Build Date:", buildDate)
\tprintln("Git Commit:", gitCommit)
\tprintln("Git Branch:", gitBranch)
}

// printHelp prints help information
func printHelp() {
\tprintln("Padlock - Enterprise AI/LLM Security Gateway")
\tprintln()
\tprintln("Usage:")
\tprintln("  padlock [options]")
\tprintln()
\tprintln("Options:")
\tflag.PrintDefaults()
\tprintln()
\tprintln("Environment Variables:")
\tprintln("  PADLOCK_CONFIG_PATH    Path to configuration file")
\tprintln("  PADLOCK_LOG_LEVEL     Log level (debug, info, warn, error)")
\tprintln("  PADLOCK_LOCALE        Locale for i18n (en, es, fr, de, etc.)")
\tprintln("  PADLOCK_MITM_ENABLED  Enable MITM proxy mode")
\tprintln("  PADLOCK_MITM_PORT     MITM proxy port (default: 3128)")
}
"""
    
    content = content + helper_funcs
    print("Helper functions added")
    
    with open(main_go_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("File updated successfully!")
    
except Exception as e:
    print(f"Error: {e}")
    sys.exit(1)


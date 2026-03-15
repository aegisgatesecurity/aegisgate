// AegisGate - Enterprise AI API Security Platform
// Main entry point for the AegisGate service
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aegisgatesecurity/aegisgate/pkg/core"
	"github.com/aegisgatesecurity/aegisgate/pkg/metrics"
	"github.com/aegisgatesecurity/aegisgate/pkg/middleware"
	"github.com/aegisgatesecurity/aegisgate/pkg/proxy"
)

// Build info - set during build
const version = "v1.0.8"
const commit  = "dev"
var date = time.Now().Format(time.RFC3339)
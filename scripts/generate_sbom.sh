#!/bin/bash
# SBOM Generation Script for AegisGate
# Generates SPDX JSON SBOM from Go dependencies

set -e

echo "Generating SBOM for AegisGate..."
cd "$(dirname "$0")"

# Create output directory
mkdir -p sbom

# Generate SBOM from Go modules
echo '{"spdxVersion": "SPDX-2.3", "dataLicense": "CC0-1.0", "SPDXID": "SPDXRef-AegisGate", "name": "AegisGate", "documentNamespace": "https://github.com/aegisgatesecurity/aegisgate/sbom", "creationInfo": {"created": "'"$(date -u +%Y-%m-%dT%H:%M:%SZ)"'", "creators": ["Tool: go-mod-sbom"]}, "packages": [' > sbom/spdx.json

# Get all modules
MODULES=$(go list -m all 2>/dev/null)

FIRST=true
while IFS= read -r module; do
    # Skip the main module
    if [[ "$module" == "github.com/aegisgatesecurity/aegisgate" ]]; then
        continue
    fi
    
    # Extract name and version
    NAME=$(echo "$module" | cut -d'@' -f1)
    VERSION=$(echo "$module" | cut -d'@' -f2)
    
    if [ -z "$VERSION" ]; then
        VERSION="v0.0.0"
    fi
    
    # Convert to SPDX package format
    if [ "$FIRST" = true ]; then
        FIRST=false
    else
        echo "," >> sbom/spdx.json
    fi
    
    cat >> sbom/spdx.json << EOF
{
  "SPDXID": "SPDXRef-$NAME",
  "name": "$NAME",
  "versionInfo": "$VERSION",
  "downloadLocation": "NOASSERTION",
  "filesAnalyzed": false
}
EOF
done <<< "$MODULES"

echo ']}' >> sbom/spdx.json

# Also generate CycloneDX JSON using go list
echo "Generating CycloneDX SBOM..."
go list -m -json all | jq -r '
  select(.Path != "github.com/aegisgatesecurity/aegisgate") | 
  {
    "bomFormat": "CycloneDX",
    "specVersion": "1.5",
    "serialNumber": "urn:uuid:" + now | strftime("%Y-%m-%d"),
    "version": 1,
    "components": [
      inputs | {
        "type": "library",
        "name": .Path,
        "version": .Version // "v0.0.0"
      }
    ]
  }
' > sbom/cyclonedx.json 2>/dev/null || echo "CyloneDX generation requires jq"

echo "SBOM generated successfully!"
echo "- SPDX: sbom/spdx.json"
echo "- CycloneDX: sbom/cyclonedx.json (if jq available)"

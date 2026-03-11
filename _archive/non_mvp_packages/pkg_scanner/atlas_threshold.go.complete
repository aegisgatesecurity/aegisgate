// Package atlas provides MITRE ATLAS scanner utilities
package atlas

// Threshold constants for ATLAS prompt injection detection
const (
	// ATLAS_THRESHOLD_MIN is the minimum score threshold for prompt injection detection
	ATLAS_THRESHOLD_MIN = 0.3

	// ATLAS_THRESHOLD_MAX is the maximum score threshold for prompt injection detection
	ATLAS_THRESHOLD_MAX = 0.4

	// Default thresholds for >95% accuracy as recommended by MITRE ATLAS
	DefaultMinThreshold = 0.3
	DefaultMaxThreshold = 0.4
)

// GetThresholdRange returns the current threshold range for ATLAS detection
func GetThresholdRange() (min, max float64) {
	return ATLAS_THRESHOLD_MIN, ATLAS_THRESHOLD_MAX
}

// ScoreThresholded checks if a detection score falls within the threshold range
func ScoreThresholded(score float64) bool {
	return score >= ATLAS_THRESHOLD_MIN && score <= ATLAS_THRESHOLD_MAX
}

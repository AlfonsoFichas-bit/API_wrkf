package tests

import (
	"os"
	"testing"
)

// TestMain serves as the entry point for tests in this package.
// With the new test_helpers.go, global setup is no longer needed here.
// Each test or test file can use SetupTestApp() for a clean, isolated environment.
func TestMain(m *testing.M) {
	// No global setup is needed anymore.
	// The test helper `SetupTestApp` provides an isolated instance for each test.

	// Run tests
	code := m.Run()

	// No global teardown is needed anymore.
	os.Exit(code)
}

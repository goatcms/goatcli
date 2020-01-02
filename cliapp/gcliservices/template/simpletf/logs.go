package simpletf

import "fmt"

// Log print log message
func Log(msg string) error {
	fmt.Printf("\n  %v", msg)
	return nil
}

// LogBuildSuccess print success build log message
func LogBuildSuccess(destPath string) error {
	Log("Build success: " + destPath)
	return nil
}

// LogBuildFail print fail build log message
func LogBuildFail(destPath string) error {
	Log(" !!!Build fail: " + destPath)
	return nil
}

// LogWarning print warning message
func LogWarning(msg string) error {
	Log("Warning: " + msg)
	return nil
}

// LogCriticError print critic error message
func LogCriticError(destPath string) error {
	Log(" !!!Build fail: " + destPath)
	return nil
}

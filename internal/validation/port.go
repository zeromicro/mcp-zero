package validation

import (
	"fmt"
	"net"
)

// ValidatePort validates a port number
// Must be between 1024-65535 and not currently in use
func ValidatePort(port int) error {
	// Check valid range (avoid privileged ports)
	if port < 1024 || port > 65535 {
		return fmt.Errorf("port must be between 1024 and 65535, got %d", port)
	}

	// Check if port is already in use
	if isPortInUse(port) {
		return fmt.Errorf("port %d is already in use", port)
	}

	return nil
}

// isPortInUse checks if a port is currently in use
func isPortInUse(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		// Port is in use
		return true
	}
	// Port is available, close the test listener
	listener.Close()
	return false
}

// SuggestAvailablePort suggests an available port starting from the given port
func SuggestAvailablePort(startPort int) (int, error) {
	// Try ports from startPort to startPort+100
	for port := startPort; port <= startPort+100 && port <= 65535; port++ {
		if !isPortInUse(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available ports found in range %d-%d", startPort, startPort+100)
}

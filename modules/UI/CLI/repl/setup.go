package repl

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func handleModSetupCommand(args []string, timestamp int64) int64 {
	fmt.Println("=== Setup Libr Mod ===")

	osName := runtime.GOOS
	arch := runtime.GOARCH
	fmt.Printf("Detected OS: %s, Architecture: %s\n", osName, arch)

	binaryName := "mod"
	if osName == "windows" {
		binaryName += ".exe"
	}
	binaryURL := fmt.Sprintf("https://example.com/binaries/%s-%s/%s", osName, arch, binaryName)
	fmt.Printf("Downloading binary from: %s\n", binaryURL)

	installDir := getInstallDir()
	err := downloadBinary(binaryURL, filepath.Join(installDir, binaryName))
	if err != nil {
		fmt.Printf("Error downloading binary: %s\n", err)
		return timestamp
	}

	fmt.Print("Enter the absolute path to your key.json file: ")
	reader := bufio.NewReader(os.Stdin)
	keyPath, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %s\n", err)
		return timestamp
	}
	keyPath = strings.TrimSpace(keyPath)

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		fmt.Printf("Error: key.json file not found at %s\n", keyPath)
		return timestamp
	}

	os.Setenv("KEY_JSON_PATH", keyPath)
	fmt.Println("Environment variables set successfully!")

	binaryPath := filepath.Join(installDir, binaryName)
	err = runBinary(binaryPath)
	if err != nil {
		fmt.Printf("Error running mod binary: %s\n", err)
		return timestamp
	}

	fmt.Println("Setup complete. You can now use Libr Mod.")
	return timestamp
}

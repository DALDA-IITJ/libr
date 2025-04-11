package repl

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func divideByHundred(ts int64) int64 {
	return ts / 100
}

func Start() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("=== Welcome to Libr ===")
	fmt.Println("Type 'help' for commands. Type '\\q' to quit.")
	fmt.Println("===========================")

	currentTimestamp := divideByHundred(time.Now().Unix())

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "\\q" {
			fmt.Println("Exiting Libr. Goodbye!")
			break
		}

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]

		switch command {
		case "send":
			handleSendCommand(args)
			currentTimestamp = divideByHundred(time.Now().Unix())
		case "fetch", "f":
			handleFetchCommand(currentTimestamp)
		case "prev", "p":
			currentTimestamp = handlePrevCommand(currentTimestamp)
		case "setup mod":
			handleModSetupCommand()
		case "help":
			printHelp()
		default:
			fmt.Println("Unknown command. Type 'help' for commands.")
		}
	}
}

// handleSetupCommand handles the setup process
func handleModSetupCommand() {
	fmt.Println("=== Setup Libr Mod ===")

	// Detect OS and architecture
	osName := runtime.GOOS
	arch := runtime.GOARCH
	fmt.Printf("Detected OS: %s, Architecture: %s\n", osName, arch)

	// Determine the binary URL
	binaryName := "mod"
	if osName == "windows" {
		binaryName += ".exe"
	}
	binaryURL := fmt.Sprintf("https://example.com/binaries/%s-%s/%s", osName, arch, binaryName)
	fmt.Printf("Downloading binary from: %s\n", binaryURL)

	// Download the binary to ~/.libr/
	installDir := getInstallDir()
	err := downloadBinary(binaryURL, filepath.Join(installDir, binaryName))
	if err != nil {
		fmt.Printf("Error downloading binary: %s\n", err)
		return
	}
	fmt.Println("Binary downloaded successfully!")

	// Prompt for key.json file path
	fmt.Print("Enter the absolute path to your key.json file: ")
	reader := bufio.NewReader(os.Stdin)
	keyPath, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %s\n", err)
		return
	}
	keyPath = strings.TrimSpace(keyPath)

	// Validate the key.json file
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		fmt.Printf("Error: key.json file not found at %s\n", keyPath)
		return
	}

	// Set environment variables
	os.Setenv("KEY_JSON_PATH", keyPath)
	fmt.Println("Environment variables set successfully!")

	// Execute the mod binary
	binaryPath := filepath.Join(installDir, binaryName)
	fmt.Println("Running mod binary...")
	err = runBinary(binaryPath)
	if err != nil {
		fmt.Printf("Error running mod binary: %s\n", err)
		return
	}
	fmt.Println("Setup complete. You can now use Libr Mod.")
}

// getInstallDir returns the installation directory for binaries
func getInstallDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		os.Exit(1)
	}

	installDir := filepath.Join(homeDir, ".libr")
	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		// Create directory if it doesn't exist
		err = os.MkdirAll(installDir, 0700)
		if err != nil {
			fmt.Println("Error creating install directory:", err)
			os.Exit(1)
		}
	}
	return installDir
}

// downloadBinary fetches the binary from the URL
func downloadBinary(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch binary: %s", resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// Make the binary executable
	err = os.Chmod(dest, 0755)
	if err != nil {
		return err
	}

	return nil
}

// runBinary executes a binary with the current environment variables
func runBinary(binaryPath string) error {
	cmd := exec.Command(binaryPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	return cmd.Run()
}

package repl

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/DALDA-IITJ/libr/modules/client"
)

func handleSendCommand(args []string, timestamp int64) int64 {
	if len(args) < 1 {
		fmt.Println("Usage: send <message>")
		return timestamp
	}

	// message := args[0]
	message := strings.Join(args, " ")
	core := client.NewCore()
	err := core.SendMessage(message)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("✔ Message sent successfully")
	}
	return divideByHundred(time.Now().Unix())
}

func handleFetchCommand(args []string, timestamp int64) int64 {
	core := client.NewCore()
	messages, err := core.FetchMessages(fmt.Sprint(timestamp))
	if err != nil {
		fmt.Println("Error fetching messages:", err)
		return timestamp
	}

	fmt.Println("\n=== Messages After Timestamp ===")
	if len(messages) == 0 {
		fmt.Println("No newer messages found.")
	} else {
		for _, msg := range messages {
			fmt.Printf("[%s] [%s] => %s\n", time.Unix(timestamp, 0).Format("2006-01-02 15:04:05"), msg.Sender, msg.Content)
		}
	}
	fmt.Println("==============================\n")
	return timestamp
}

func handlePrevCommand(args []string, currentBucket int64) int64 {
	prevBucket := currentBucket - 1
	core := client.NewCore()
	messages, err := core.FetchMessages(fmt.Sprint(prevBucket))
	if err != nil {
		fmt.Println("Error fetching messages:", err)
		return currentBucket // Don't update timestamp on error
	}

	fmt.Println("\n=== Messages in Previous Bucket ===")
	if len(messages) == 0 {
		fmt.Println("No messages found in previous bucket.")
	} else {
		for _, msg := range messages {
			fmt.Printf("[%s] [%s] => %s\n", time.Unix(currentBucket*100, 0).Format("2006-01-02 15:04:05"), msg.Sender, msg.Content)
		}
	}
	fmt.Println("===============================\n")
	return prevBucket
}

func handleHelpCommand(args []string, timestamp int64) int64 {
	fmt.Println("\n=== CLI Commands ===")
	fmt.Println("send <message>     - Send a message.")
	fmt.Println("f, fetch           - Fetch messages after the current timestamp.")
	fmt.Println("p, prev            - Fetch messages before the current timestamp.")
	fmt.Println("\\q                 - Quit the CLI.")
	fmt.Println("====================\n")
	return timestamp
}

func handleModSetupCommand(args []string, timestamp int64) int64 {
	fmt.Println("=== Setup Libr Mod ===")

	osName := runtime.GOOS
	arch := runtime.GOARCH
	fmt.Printf("Detected OS: %s, Architecture: %s\n", osName, arch)

	binaryName := "mod"
	if osName == "windows" {
		binaryName += ".exe"
	}
	//  https://github.com/DALDA-IITJ/libr/releases/download/v0.0.0/libr-mod.exe
	// binaryURL := fmt.Sprintf("https://example.com/binaries/%s-%s/%s", osName, arch, binaryName)
	binaryURL := fmt.Sprintf("https://github.com/DALDA-IITJ/libr/releases/download/modules/mod/v0.1.0/libr-mod.exe")
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

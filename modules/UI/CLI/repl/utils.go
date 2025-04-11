package repl

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func getInstallDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		os.Exit(1)
	}

	installDir := filepath.Join(homeDir, ".libr")
	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		err = os.MkdirAll(installDir, 0700)
		if err != nil {
			fmt.Println("Error creating install directory:", err)
			os.Exit(1)
		}
	}
	return installDir
}

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

	return os.Chmod(dest, 0755)
}

func runBinary(binaryPath string) error {
	cmd := exec.Command(binaryPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	return cmd.Run()
}

package installer

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kelleyblackmore/opsdev/internal/utils"
)

type ToolInstaller struct {
	OS    string
	Tools []Tool
}

func NewToolInstaller() *ToolInstaller {
	return &ToolInstaller{
		OS:    runtime.GOOS,
		Tools: GetDefaultTools(),
	}
}

func (t *ToolInstaller) StartSetup() error {
	if !t.isSupportedOS() {
		return fmt.Errorf("unsupported operating system: %s", t.OS)
	}

	fmt.Printf("%sDevOps Environment Setup%s\n", utils.ColorGreen, utils.ColorNC)
	fmt.Printf("%sDetected OS: %s%s\n\n", utils.ColorYellow, t.OS, utils.ColorNC)

	selectedTools, err := t.checkAndSelectTools()
	if err != nil {
		return fmt.Errorf("error selecting tools: %v", err)
	}

	return t.installSelectedTools(selectedTools)
}

func (t *ToolInstaller) isSupportedOS() bool {
	return t.OS == "linux" || t.OS == "darwin"
}

func (t *ToolInstaller) isToolInstalled(tool Tool) (bool, string) {
	cmd := exec.Command("sh", "-c", tool.CheckCommand)
	if err := cmd.Run(); err != nil {
		return false, ""
	}

	versionCmd := exec.Command("sh", "-c", tool.InfoCommand)
	output, err := versionCmd.Output()
	if err != nil {
		return true, "version unknown"
	}
	return true, strings.TrimSpace(string(output))
}

func (t *ToolInstaller) checkAndSelectTools() (map[string]string, error) {
	selected := make(map[string]string)
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%sChecking installed tools...%s\n\n", utils.ColorBlue, utils.ColorNC)

	for _, tool := range t.Tools {
		installed, versionInfo := t.isToolInstalled(tool)
		if installed {
			fmt.Printf("%s✓ %s is already installed%s\n", utils.ColorGreen, tool.Name, utils.ColorNC)
			fmt.Printf("  Current version: %s\n", versionInfo)
			fmt.Printf("Would you like to reinstall/update %s? (y/n): ", tool.Name)
		} else {
			fmt.Printf("%s✗ %s is not installed%s\n", utils.ColorYellow, tool.Name, utils.ColorNC)
			fmt.Printf("Would you like to install %s? (y/n): ", tool.Name)
		}

		answer, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading input: %v", err)
		}
		answer = strings.TrimSpace(strings.ToLower(answer))

		if answer == "y" {
			version := "latest"
			if len(tool.Versions) > 0 {
				version = t.selectVersion(tool)
			}
			selected[tool.Name] = version
		}
		fmt.Println()
	}

	return selected, nil
}

func (t *ToolInstaller) selectVersion(tool Tool) string {
	fmt.Printf("%sAvailable versions for %s:%s\n", utils.ColorYellow, tool.Name, utils.ColorNC)
	for i, version := range tool.Versions {
		fmt.Printf("%d) %s\n", i+1, version)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Select %s version (1-%d): ", tool.Name, len(tool.Versions))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		index := 0
		fmt.Sscanf(input, "%d", &index)

		if index > 0 && index <= len(tool.Versions) {
			return tool.Versions[index-1]
		}
	}
}

func (t *ToolInstaller) installSelectedTools(selected map[string]string) error {
	for tool, version := range selected {
		fmt.Printf("\n%sInstalling %s version %s for %s...%s\n",
			utils.ColorGreen, tool, version, t.OS, utils.ColorNC)

		var err error
		switch tool {
		case "aws-cli":
			err = t.installAWSCLI()
		case "azure-cli":
			err = t.installAzureCLI()
		case "terraform":
			err = t.installTerraform(version)
		// Add other tools here
		default:
			fmt.Printf("%sSkipping %s: installation not implemented%s\n",
				utils.ColorYellow, tool, utils.ColorNC)
			continue
		}

		if err != nil {
			return fmt.Errorf("error installing %s: %v", tool, err)
		}
	}
	return nil
}

func (t *ToolInstaller) installAWSCLI() error {
	if t.OS == "linux" {
		return t.downloadAndUnzip(
			"https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip",
			"awscliv2.zip",
			func() error {
				cmd := exec.Command("sudo", "./aws/install")
				return cmd.Run()
			},
		)
	} else if t.OS == "darwin" {
		return t.downloadFile(
			"https://awscli.amazonaws.com/AWSCLIV2.pkg",
			"AWSCLIV2.pkg",
			func() error {
				cmd := exec.Command("sudo", "installer", "-pkg", "AWSCLIV2.pkg", "-target", "/")
				return cmd.Run()
			},
		)
	}
	return fmt.Errorf("unsupported OS for AWS CLI installation")
}

func (t *ToolInstaller) installAzureCLI() error {
	if t.OS == "linux" {
		cmd := exec.Command("curl", "-sL", "https://aka.ms/InstallAzureCLIDeb", "|", "sudo", "bash")
		return cmd.Run()
	} else if t.OS == "darwin" {
		cmd := exec.Command("brew", "install", "azure-cli")
		return cmd.Run()
	}
	return fmt.Errorf("unsupported OS for Azure CLI installation")
}

func (t *ToolInstaller) installTerraform(version string) error {
	osName := "linux"
	if t.OS == "darwin" {
		osName = "darwin"
	}

	url := fmt.Sprintf(
		"https://releases.hashicorp.com/terraform/%s/terraform_%s_%s_amd64.zip",
		version, version, osName,
	)

	return t.downloadAndUnzip(url, "terraform.zip", func() error {
		cmd := exec.Command("sudo", "mv", "terraform", "/usr/local/bin/")
		return cmd.Run()
	})
}

func (t *ToolInstaller) downloadFile(url, filepath string, afterDownload func() error) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading file: %v", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	if afterDownload != nil {
		return afterDownload()
	}
	return nil
}

func (t *ToolInstaller) downloadAndUnzip(url, filepath string, afterUnzip func() error) error {
	err := t.downloadFile(url, filepath, nil)
	if err != nil {
		return err
	}

	zipReader, err := zip.OpenReader(filepath)
	if err != nil {
		return fmt.Errorf("error opening zip file: %v", err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		path := filepath.Join(file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return fmt.Errorf("error creating directory: %v", err)
		}

		dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("error creating file: %v", err)
		}

		fileInArchive, err := file.Open()
		if err != nil {
			dstFile.Close()
			return fmt.Errorf("error opening file in archive: %v", err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			dstFile.Close()
			fileInArchive.Close()
			return fmt.Errorf("error copying file: %v", err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	if afterUnzip != nil {
		return afterUnzip()
	}
	return nil
}
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var errCancelled = errors.New("cancelled")

var rootCmd = &cobra.Command{
	Use:   "wacomctl",
	Short: "Controls the Wacom stylus for detected monitors",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleCommand(args)
	},
}

func main() {
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	if err := rootCmd.Execute(); err != nil {
		if errors.Is(err, errCancelled) || errors.Is(err, promptui.ErrInterrupt) || errors.Is(err, syscall.EINTR) {
			fmt.Println("Selection cancelled.")
			return
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Printf("Usage: %s [interactive|map <monitor>|both|off|on|vga|hdmi]\n", os.Args[0])
	fmt.Println("Maps or controls the stylus for a detected monitor.")
	fmt.Println("  interactive  Opens an interactive selector to choose a monitor")
	fmt.Println("  map          Maps the stylus to the provided monitor")
	fmt.Println("  both         Maps the stylus to all active monitors")
	fmt.Println("  off          Turns the stylus off")
	fmt.Println("  on           Turns the stylus on")
}

func handleCommand(args []string) error {
	if len(args) == 0 {
		return runInteractiveSelection()
	}

	command := strings.ToLower(args[0])
	switch command {
	case "interactive", "select", "i":
		return runInteractiveSelection()
	case "map":
		if len(args) < 2 {
			showHelp()
			return fmt.Errorf("please provide the monitor name")
		}
		return mapToOutput(args[1])
	case "both":
		return mapToBoth()
	case "off":
		return turnOff()
	case "on":
		return turnOn()
	case "vga":
		return mapToMonitorByHint("VGA")
	case "hdmi":
		return mapToMonitorByHint("HDMI")
	default:
		if len(args) == 1 {
			return mapToOutput(args[0])
		}
		showHelp()
		return fmt.Errorf("invalid parameter: %s", args[0])
	}
}

func getStylusDeviceID() (string, error) {
	cmd := exec.Command("xsetwacom", "--list", "devices")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("xsetwacom --list devices failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "stylus") {
			parts := strings.Fields(line)
			if len(parts) >= 7 {
				return parts[6], nil
			}
		}
	}
	return "", fmt.Errorf("stylus device not found")
}

func listMonitors() ([]string, error) {
	cmd := exec.Command("xrandr", "--listmonitors")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("xrandr --listmonitors failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return parseMonitors(string(output))
}

func parseMonitors(output string) ([]string, error) {
	monitors := []string{}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "Monitors:") {
			continue
		}

		parts := strings.Fields(trimmed)
		if len(parts) < 2 {
			continue
		}

		monitorName := strings.TrimLeft(parts[1], "+*")
		monitorName = strings.TrimSpace(monitorName)
		if monitorName == "" {
			continue
		}

		if !contains(monitors, monitorName) {
			monitors = append(monitors, monitorName)
		}
	}

	if len(monitors) == 0 {
		return nil, fmt.Errorf("no monitors found")
	}

	return monitors, nil
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func mapToMonitorByHint(hint string) error {
	monitors, err := listMonitors()
	if err != nil {
		return err
	}

	for _, monitor := range monitors {
		if strings.Contains(strings.ToUpper(monitor), strings.ToUpper(hint)) {
			return mapToOutput(monitor)
		}
	}

	return fmt.Errorf("no monitor %s found", hint)
}

func mapToOutput(target string) error {
	deviceID, err := getStylusDeviceID()
	if err != nil {
		return err
	}

	resolvedTarget, err := resolveWacomTarget(target)
	if err != nil {
		return err
	}

	cmd := exec.Command("xsetwacom", "set", deviceID, "MapToOutput", resolvedTarget)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to map stylus: %w: %s", err, strings.TrimSpace(string(output)))
	}

	if err := savePersistedTarget(resolvedTarget); err != nil {
		fmt.Printf("Warning: could not persist target: %v\n", err)
	}

	fmt.Printf("Stylus mapped to monitor %s.\n", resolvedTarget)
	return nil
}

func mapToBoth() error {
	deviceID, err := getStylusDeviceID()
	if err != nil {
		return err
	}

	cmd := exec.Command("xsetwacom", "set", deviceID, "MapToOutput", "desktop")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to map stylus: %w: %s", err, strings.TrimSpace(string(output)))
	}

	if err := savePersistedTarget("desktop"); err != nil {
		fmt.Printf("Warning: could not persist target: %v\n", err)
	}

	fmt.Println("Stylus mapped to all monitors.")
	return nil
}

func resolveWacomTarget(target string) (string, error) {
	if target == "desktop" {
		return target, nil
	}

	monitors, err := listMonitors()
	if err == nil {
		for index, monitor := range monitors {
			if strings.EqualFold(monitor, target) {
				return resolveWacomTargetName(target, index), nil
			}
		}
	}

	if strings.HasPrefix(strings.ToUpper(target), "HEAD-") {
		return target, nil
	}

	return target, fmt.Errorf("could not resolve a valid monitor target for %s", target)
}

func resolveWacomTargetName(target string, index int) string {
	if target == "desktop" {
		return target
	}
	return fmt.Sprintf("HEAD-%d", index)
}

func turnOff() error {
	deviceID, err := getStylusDeviceID()
	if err != nil {
		return err
	}

	cmd := exec.Command("xinput", "disable", deviceID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to turn off stylus: %w: %s", err, strings.TrimSpace(string(output)))
	}

	fmt.Println("Stylus turned off.")
	return nil
}

func turnOn() error {
	deviceID, err := getStylusDeviceID()
	if err != nil {
		return err
	}

	cmd := exec.Command("xinput", "enable", deviceID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to turn on stylus: %w: %s", err, strings.TrimSpace(string(output)))
	}

	persistedTarget, err := loadPersistedTarget()
	if err != nil {
		return nil
	}

	if persistedTarget == "" || persistedTarget == "desktop" {
		fmt.Println("Stylus turned on.")
		return nil
	}

	if err := mapToOutput(persistedTarget); err != nil {
		fmt.Printf("Warning: could not restore persisted target %q: %v\n", persistedTarget, err)
	}

	return nil
}

func stateFilePath() string {
	if override := os.Getenv("WACOMCTL_STATE_FILE"); override != "" {
		return override
	}
	return filepath.Join(os.Getenv("HOME"), ".wacomctl-state")
}

func savePersistedTarget(target string) error {
	if target == "" {
		return nil
	}
	return os.WriteFile(stateFilePath(), []byte(target), 0o600)
}

func loadPersistedTarget() (string, error) {
	data, err := os.ReadFile(stateFilePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func runInteractiveSelection() error {
	monitors, err := listMonitors()
	if err != nil {
		fmt.Printf("Could not list monitors: %v\n", err)
		monitors = nil
	}

	options := buildSelectionOptions(monitors)
	if len(options) == 0 {
		return fmt.Errorf("no options available")
	}

	selected, err := promptSelection(options)
	if err != nil {
		return err
	}

	switch selected {
	case "both":
		return mapToBoth()
	case "turn off":
		return turnOff()
	case "turn on":
		return turnOn()
	default:
		return mapToOutput(selected)
	}
}

func buildSelectionOptions(monitors []string) []string {
	options := append([]string{}, monitors...)
	options = append(options, "both", "turn off", "turn on")
	return options
}

func promptSelection(options []string) (string, error) {
	if !isInteractiveTerminal() {
		index, err := promptSelectionByNumber(options)
		if err != nil {
			return "", err
		}
		return options[index], nil
	}

	prompt := promptui.Select{
		Label:    "Stylus destination",
		Items:    options,
		Size:     8,
		HideHelp: true,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "▶ {{ . | cyan }}",
			Inactive: "  {{ . | faint }}",
			Selected: "✓ {{ . | green }}",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) || errors.Is(err, syscall.EINTR) {
			return "", errCancelled
		}
		return "", err
	}
	return result, nil
}

func promptSelectionByNumber(options []string) (int, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Select an option:")
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option)
	}
	fmt.Print("Enter the option number: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	choice, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || choice < 1 || choice > len(options) {
		return 0, fmt.Errorf("invalid option")
	}

	return choice - 1, nil
}

func isInteractiveTerminal() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

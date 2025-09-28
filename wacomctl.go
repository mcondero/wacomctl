package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		showHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "vga":
		mapToVGA()
	case "hdmi":
		mapToHDMI()
	case "both":
		mapToBoth()
	case "off":
		turnOff()
	case "on":
		turnOn()
	default:
		fmt.Println("Parâmetro inválido.")
		showHelp()
	}
}

func showHelp() {
	fmt.Printf("Uso: %s [vga|hdmi|both|off|on]\n", os.Args[0])
	fmt.Println("Mapeia ou controla a stylus para o monitor especificado.")
	fmt.Println("  vga    Mapeia a stylus para o monitor VGA (ex: VGA-1)")
	fmt.Println("  hdmi   Mapeia a stylus para o monitor HDMI (ex: HDMI-1)")
	fmt.Println("  both   Mapeia a stylus para todos os monitores ativos")
	fmt.Println("  off    Desliga a stylus")
	fmt.Println("  on     Liga a stylus")
	os.Exit(1)
}

func getStylusDeviceID() (string, error) {
	cmd := exec.Command("xsetwacom", "--list", "devices")
	output, err := cmd.Output()
	if err != nil {
		return "", err
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
	return "", fmt.Errorf("dispositivo de stylus não encontrado")
}

func getVGAMonitor() (string, error) {
	cmd := exec.Command("xrandr", "--listmonitors")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToUpper(line), "VGA") {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				return parts[3], nil
			}
		}
	}
	return "", fmt.Errorf("nenhum monitor VGA encontrado")
}

func getHDMIMonitor() (string, error) {
	cmd := exec.Command("xrandr", "--listmonitors")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToUpper(line), "HDMI") {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				return parts[3], nil
			}
		}
	}
	return "", fmt.Errorf("nenhum monitor HDMI encontrado")
}

func mapToVGA() {
	deviceID, err := getStylusDeviceID()
	if err != nil {
		fmt.Println("Dispositivo de stylus não encontrado.")
		os.Exit(1)
	}

	vgaMonitor, err := getVGAMonitor()
	if err != nil {
		fmt.Println("Nenhum monitor VGA encontrado.")
		return
	}

	cmd := exec.Command("xsetwacom", "set", deviceID, "MapToOutput", vgaMonitor)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Erro ao mapear stylus: %v\n", err)
		return
	}

	fmt.Printf("Stylus mapeada para o monitor %s.\n", vgaMonitor)
}

func mapToHDMI() {
	deviceID, err := getStylusDeviceID()
	if err != nil {
		fmt.Println("Dispositivo de stylus não encontrado.")
		os.Exit(1)
	}

	hdmiMonitor, err := getHDMIMonitor()
	if err != nil {
		fmt.Println("Nenhum monitor HDMI encontrado.")
		return
	}

	cmd := exec.Command("xsetwacom", "set", deviceID, "MapToOutput", hdmiMonitor)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Erro ao mapear stylus: %v\n", err)
		return
	}

	fmt.Printf("Stylus mapeada para o monitor %s.\n", hdmiMonitor)
}

func mapToBoth() {
	deviceID, err := getStylusDeviceID()
	if err != nil {
		fmt.Println("Dispositivo de stylus não encontrado.")
		os.Exit(1)
	}

	cmd := exec.Command("xsetwacom", "set", deviceID, "MapToOutput", "desktop")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Erro ao mapear stylus: %v\n", err)
		return
	}

	fmt.Println("Stylus mapeada para todos os monitores.")
}

func turnOff() {
	deviceID, err := getStylusDeviceID()
	if err != nil {
		fmt.Println("Dispositivo de stylus não encontrado.")
		os.Exit(1)
	}

	cmd := exec.Command("xinput", "disable", deviceID)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Erro ao desligar stylus: %v\n", err)
		return
	}

	fmt.Println("Stylus desligada.")
}

func turnOn() {
	deviceID, err := getStylusDeviceID()
	if err != nil {
		fmt.Println("Dispositivo de stylus não encontrado.")
		os.Exit(1)
	}

	cmd := exec.Command("xinput", "enable", deviceID)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Erro ao ligar stylus: %v\n", err)
		return
	}

	fmt.Println("Stylus ligada.")
}
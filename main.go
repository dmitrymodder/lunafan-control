package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Point struct {
	Temp    float64 `json:"temp"`
	Percent float64 `json:"percent"`
}

type Fan struct {
	Name      string  `json:"name"`
	PwmPath   string  `json:"pwm_path"`
	InputPath string  `json:"input_path"`
	Curve     []Point `json:"curve"`
}

type Config struct {
	TempSensor       string `json:"temp_sensor"`
	UpdateIntervalMs int    `json:"update_interval_ms"`
	Fans             []Fan  `json:"fans"`
}

func loadConfig() Config {
	data, err := os.ReadFile("/etc/lunafan-control/config.json")
	if err != nil {
		log.Fatal("Failed to read config: ", err)
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Failed to parse config: ", err)
	}
	return config
}

func computePercent(temp float64, curve []Point) float64 {
	if len(curve) == 0 {
		return 0
	}
	if temp <= curve[0].Temp {
		return curve[0].Percent
	}
	n := len(curve)
	if temp >= curve[n-1].Temp {
		return curve[n-1].Percent
	}
	for i := 0; i < n-1; i++ {
		if temp > curve[i].Temp && temp <= curve[i+1].Temp {
			dt := curve[i+1].Temp - curve[i].Temp
			dp := curve[i+1].Percent - curve[i].Percent
			frac := (temp - curve[i].Temp) / dt
			return curve[i].Percent + frac*dp
		}
	}
	return 0
}

func getEnablePath(pwmPath string) string {
	dir := filepath.Dir(pwmPath)
	base := filepath.Base(pwmPath)
	return filepath.Join(dir, base+"_enable")
}

func runLoop() {
	config := loadConfig()
	originalEnables := make(map[string]string)
	for _, fan := range config.Fans {
		enablePath := getEnablePath(fan.PwmPath)
		data, err := os.ReadFile(enablePath)
		if err == nil {
			originalEnables[enablePath] = strings.TrimSpace(string(data))
		}
		os.WriteFile(enablePath, []byte("1\n"), 0644)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		for path, val := range originalEnables {
			os.WriteFile(path, []byte(val+"\n"), 0644)
		}
		os.Exit(0)
	}()

	for {
		tempStr, err := os.ReadFile(config.TempSensor)
		if err != nil {
			log.Println("Failed to read temp: ", err)
			time.Sleep(time.Second)
			continue
		}
		tempMilli, err := strconv.Atoi(strings.TrimSpace(string(tempStr)))
		if err != nil {
			log.Println("Failed to parse temp: ", err)
			time.Sleep(time.Second)
			continue
		}
		temp := float64(tempMilli) / 1000.0

		for _, fan := range config.Fans {
			percent := computePercent(temp, fan.Curve)
			pwm := int(percent/100*255 + 0.5)
			if pwm > 255 {
				pwm = 255
			}
			if pwm < 0 {
				pwm = 0
			}
			err := os.WriteFile(fan.PwmPath, []byte(strconv.Itoa(pwm)+"\n"), 0644)
			if err != nil {
				log.Println("Failed to write pwm: ", err)
			}
		}
		time.Sleep(time.Duration(config.UpdateIntervalMs) * time.Millisecond)
	}
}

func printStats() {
	config := loadConfig()
	for _, fan := range config.Fans {
		if fan.InputPath != "" {
			rpmStr, _ := os.ReadFile(fan.InputPath)
			rpm, _ := strconv.Atoi(strings.TrimSpace(string(rpmStr)))
			fmt.Printf("%s: %d RPM\n", fan.Name, rpm)
		}
	}
}

func setConfig(name string) {
	if os.Geteuid() != 0 {
		log.Fatal("This operation requires root privileges")
	}
	configDir := "/etc/lunafan-control/configs"
	target := "/etc/lunafan-control/config.json"
	src := filepath.Join(configDir, name+".json")
	if _, err := os.Stat(src); err != nil {
		log.Fatal("Config not found: ", src)
	}
	os.Remove(target)
	err := os.Symlink(src, target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Config set to ", name)
}

func manageService(action string) {
	if os.Geteuid() != 0 {
		log.Fatal("This operation requires root privileges")
	}
	cmd := exec.Command("systemctl", action, "lunafan-control.service")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: lunafan-control [start|stop|enable|disable|stats|config <name>|run]")
		os.Exit(1)
	}
	cmd := os.Args[1]
	switch cmd {
		case "run":
			runLoop()
		case "start":
			manageService("start")
		case "stop":
			manageService("stop")
		case "enable":
			manageService("enable")
		case "disable":
			manageService("disable")
		case "stats":
			printStats()
		case "config":
			if len(os.Args) < 3 {
				fmt.Println("Usage: lunafan-control config <name>")
				os.Exit(1)
			}
			setConfig(os.Args[2])
		default:
			fmt.Println("Unknown command: ", cmd)
			os.Exit(1)
	}
}

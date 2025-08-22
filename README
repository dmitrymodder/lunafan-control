# LunaFan Control

LunaFan Control is a lightweight daemon for managing fan speeds in Linux based on CPU temperature or other sensors. Written in Go, it uses JSON configuration for flexible fan speed curve customization. It serves as an alternative to `fancontrol` with an emphasis on simplicity and precision.

## Features
- Simple setup via JSON config.
- Linear interpolation for smooth fan control.
- Support for multiple fans with individual curves.
- Systemd integration for running in the background.
- Commands to start, stop, view status, and switch configs.
- Easily adaptable to different hardware (requires hwmon path configuration).

## Requirements
- Linux with access to `/sys/class/hwmon` (typically provided by the kernel).
- Go 1.16+ (for building from source).
- For Arch Linux: `makepkg` (optional, for creating a package).
- Root privileges to write to `/sys/class/hwmon` and install the service.

## Installation

### Option 1: Building ELF Binary (Any Linux)
1. Ensure Go is installed: `go version`.
2. Clone the repository:
   ```bash
   git clone https://github.com/dmitrymodder/lunafan-control.git
   cd lunafan-control
   ```
3. Build the binary:
   ```bash
   make build
   ```
4. Install the program and service:
   ```bash
   sudo make install
   ```

### Option 2: Building Arch Package
1. Clone the repository:
   ```bash
   git clone https://github.com/dmitrymodder/lunafan-control.git
   cd lunafan-control
   ```
2. Build and install the package:
   ```bash
   make package
   sudo pacman -U lunafan-control-*.pkg.tar.zst
   ```

### Post-Build Cleanup
To remove temporary files and folders  (`pkg`, `src`, packages):
```bash
make clean
```

## Configuration
1. Edit `/etc/lunafan-control/config.json` for your hardware:
   - `temp_sensor`: path to the temperature file (example, `/sys/class/hwmon/hwmon1/temp1_input`).
   - `fans`: list of fans with paths to PWM (`pwm_path`), RPM (`input_path`) and speed curves (`curve`).
   - Exaple config:
     ```json
     {
       "temp_sensor": "/sys/class/hwmon/hwmon1/temp1_input",
       "update_interval_ms": 5000,
       "fans": [
         {
           "name": "cpu",
           "pwm_path": "/sys/class/hwmon/hwmon3/pwm2",
           "input_path": "/sys/class/hwmon/hwmon3/fan2_input",
           "curve": [
             {"temp": 35, "percent": 20},
             {"temp": 50, "percent": 35},
             {"temp": 75, "percent": 100}
           ]
         }
       ]
     }
     ```
2. To switch between configs, create files in `/etc/lunafan-control/configs/` and use:
   ```bash
   lunafan-control config <config_name>
   ```

## Usage
- Start the service:
  ```bash
  sudo lunafan-control start
  ```
- Stop the service:
  ```bash
  sudo lunafan-control stop
  ```
- Enable auto-start:
  ```bash
  sudo lunafan-control enable
  ```
- Check status (temperature and RPM):
  ```bash
  lunafan-control stats
  ```
- Run in foreground for debugging:
  ```bash
  sudo lunafan-control run
  ```

## Finding hwmon Paths
To find the correct paths for `temp_sensor`, `pwm_path` and `input_path`:
1. Explore `/sys/class/hwmon`:
   ```bash
   ls /sys/class/hwmon
   ```
2. Look for the temperature sensor (usually `tempX_input`) and fans (`pwmX`, `fanX_input`).
3. Example: CPU temperature might be in  `/sys/class/hwmon/hwmon1/temp1_input`, and a fan in `/sys/class/hwmon/hwmon3/pwm1`.

## License
MIT License. См. файл `LICENSE`.

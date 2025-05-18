package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/jedrw/gowake/pkg/magicpacket"
	"gopkg.in/yaml.v3"
)

type relayTarget struct {
	Mac  string `yaml:"mac"`
	Ip   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

func main() {
	configBytes, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config := map[string]relayTarget{}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ip := os.Getenv("MP_LISTEN_IP")
	if ip == "" {
		ip = "0.0.0.0"
	}

	portString := os.Getenv("MP_LISTEN_PORT")
	if portString == "" {
		portString = "9"
	}
	port, err := strconv.Atoi(portString)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		fmt.Printf("listening for magic packets on %s:%d \n", ip, port)
		_, mac, err := magicpacket.Listen(ip, port)
		if err != nil {
			var errno syscall.Errno
			if errors.As(err, &errno) {
				if errno == syscall.EACCES {
					fmt.Printf("%s: please run as elevated user\n", err)
					os.Exit(1)
				}
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		fmt.Printf("heard magic packet for %s, checking for match in config\n", mac)
		target, exists := config[mac]
		if !exists {
			fmt.Println("no match in config")
		} else {
			fmt.Printf("found match for %s in config, relaying to %s\n", mac, target.Mac)
			relayPacket, err := magicpacket.New(target.Mac)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if target.Port == 0 {
				target.Port = 9
			}

			err = magicpacket.Send(relayPacket, target.Ip, target.Port)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("sent magic packet to %s at %s:%d\n", target.Mac, target.Ip, target.Port)
		}
	}
}

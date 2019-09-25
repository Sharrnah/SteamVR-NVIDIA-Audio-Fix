package main

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"log"
	"os"
)

func main() {
	var names []string

	// Open Registry Key with Class = MEDIA (is it always {4d36e96c-e325-11ce-bfc1-08002be10318} ?)
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Class\{4d36e96c-e325-11ce-bfc1-08002be10318}`, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	names, err = k.ReadSubKeyNames(-1)
	if err != nil {
		log.Fatal(err)
	}

	// Search for Subkey with DriverDescription = "NVIDIA High Definition Audio"
	for _, k := range names {
		if k != "Configuration" && k != "Properties" { // skip Configuration and Properties SubKeys
			subKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Class\{4d36e96c-e325-11ce-bfc1-08002be10318}\`+k, registry.READ)
			if err != nil {
				fmt.Println("could not read key "+k+": ", err)
			}
			DriverDescValue, _, err := subKey.GetStringValue("DriverDesc")
			subKey.Close()
			if DriverDescValue == "NVIDIA High Definition Audio" {
				fmt.Println("Found "+DriverDescValue+" device. setting Power Settings...")
				changeKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Class\{4d36e96c-e325-11ce-bfc1-08002be10318}\`+k+`\PowerSettings`, registry.ALL_ACCESS)

				// Set PowerSetting Values
				if err != nil {
					log.Fatal(`could not open key `+k+`\PowerSettings: `, err)
				}
				err = changeKey.SetBinaryValue("ConservationIdleTime", []byte{0xff,0xff,0xff,0xff})
				if err != nil {
					log.Fatal("could not set value of ConservationIdleTime: ", err)
				}
				err = changeKey.SetBinaryValue("IdlePowerState", []byte{0x00,0x00,0x00,0x00})
				if err != nil {
					log.Fatal("could not set value of IdlePowerState: ", err)
				}
				err = changeKey.SetBinaryValue("PerformanceIdleTime", []byte{0xff,0xff,0xff,0xff})
				if err != nil {
					log.Fatal("could not set value of PerformanceIdleTime: ", err)
				}
				changeKey.Close()
			}
		}
	}
	fmt.Println("Finished!")

	if (len(os.Args) < 2) || (len(os.Args) >= 2 && os.Args[1] != "--close") {
		fmt.Println("Press the Enter Key to close.")
		var input string
		fmt.Scanln(&input)
	}
}
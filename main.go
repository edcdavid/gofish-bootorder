// SPDX-License-Identifier: BSD-3-Clause
package main

import (
	"fmt"
	"os"

	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/redfish"
)

func main() {
	// Create a new instance of gofish client, ignoring self-signed certs
	config := gofish.ClientConfig{
		Endpoint: "https://xxx.xxx.xxx.xxx",
		Username: "login",
		Password: "password",
		Insecure: true,
	}
	c, err := gofish.Connect(config)
	if err != nil {
		panic(err)
	}
	defer c.Logout()

	// Retrieve the service root
	service := c.Service

	// Query the computer systems
	ss, err := service.Systems()
	if err != nil {
		panic(err)
	}

	for _, system := range ss {
		// Get Manufacturer
		fmt.Printf("model: %#v\n\n", system.Manufacturer)

		// Get Old boot order
		fmt.Printf("Old boot Order: %#v\n\n", system.Boot.BootOrder)

		//Get Boot options (details of each boot target in boot order)
		myBootOptions, err := system.BootOptions()
		if err != nil {
			fmt.Printf("error getting  boot options\n")
			os.Exit(1)
		}
		fmt.Printf("Boot Options: \n")
		for i := 0; i < len(myBootOptions); i++ {
			fmt.Printf("%#v\n\n", myBootOptions[i])
		}

		// Creates a boot override to uefi once
		bootOverride := redfish.Boot{}

		// copy old boot order to new one
		bootOverride.BootOrder = system.Boot.BootOrder

		//swapping the 2 first boot targets
		backup := system.Boot.BootOrder[0]
		bootOverride.BootOrder[0] = system.Boot.BootOrder[1]
		bootOverride.BootOrder[1] = backup
		fmt.Printf("New boot Order: %#v\n\n", system.Boot.BootOrder)

		// set new boot order
		err = system.SetBoot(bootOverride)
		if err != nil {
			fmt.Printf("error setting boot, err=%s", err)
			os.Exit(1)
		}

		// restart the system to apply the changes
		/*err = system.Reset(redfish.ForceRestartResetType)
		if err != nil {
			panic(err)
		}*/
	}
}

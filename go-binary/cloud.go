package main

import (
	"os"
	"log"
	"os/exec"
)

// This Go application parses input from the upstream MSI installer/command line, and calls VNC Server to perform the required action
func main() {
	if len(os.Args) < 2 {
		// Check that the correct number of arguments have been provided
		log.Fatal("Usage: cloud.exe <join/leave> [<token>]")
	} else if len(os.Args) < 3 && os.Args[1] == "join" {
		// If join argument is provided, check for required token. If missing, exit.
		log.Print("No token providing, skipping cloud join")
		os.Exit(0)
	} else if len(os.Args) < 3 && os.Args[1] == "leave" {
		// If leave argument provided, tell VNC Server to remove itself from the cloud.
		log.Print("Removing VNC Server from the cloud...")
		svr := os.ExpandEnv("$ProgramFiles\\RealVNC\\VNC Server\\vncserver.exe")
		cmd := exec.Command(svr, "-service", "-leaveCloud")
		err := cmd.Start()
		cmd.Wait() //wait for VNC Server command to complete.
		if err != nil {
			//Something went wrong, log a fatal error and exit.
			log.Fatal(err)
		} else {
			//Success, log success message.
			log.Print("VNC Server has been removed from the cloud")
			os.Exit(0)
		}
	} else if os.Args[1] == "join"{
		// If join argument provided, use token to join VNC Server to the cloud
		log.Print("Joining VNC Server to the cloud...")
		svr := os.ExpandEnv("$ProgramFiles\\RealVNC\\VNC Server\\vncserver.exe")
		cmd := exec.Command(svr, "-service", "-joinCloud", os.Args[2])
		err := cmd.Start()
		cmd.Wait() //wait for VNC Server command to complete.
		if err != nil {
			//Something went wrong, log a fatal error and exit.
			log.Fatal(err)
		} else {
			//Success, log success message.
			log.Print("VNC Server has been joined to the cloud")
			os.Exit(0)
		}
	}
}

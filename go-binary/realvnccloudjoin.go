package main

import (
	"log"
	"os"
	"os/exec"
)

// This Go application parses input from the upstream MSI installer/command line, and calls RealVNC Server to perform the required action
func main() {
	if len(os.Args) < 2 {
		// Check that the correct number of arguments have been provided
		log.Fatal("Usage: realvnccloudjoin.exe <join/leave> [<token>]")
	} else if len(os.Args) < 3 && os.Args[1] == "join" {
		// If join argument is provided, check for required token. If missing, exit.
		log.Print("No cloud connectivity token provided, skipping cloud join")
		os.Exit(0)
	} else if len(os.Args) < 3 && os.Args[1] == "leave" {
		// If leave argument provided, tell RealVNC Server to remove itself from the cloud.
		log.Print("Removing RealVNC Server from the cloud...")
		svr := os.ExpandEnv("$ProgramFiles\\RealVNC\\VNC Server\\vncserver.exe")
		cmd := exec.Command(svr, "-service", "-leaveCloud")
		err := cmd.Start()
		cmd.Wait() //wait for RealVNC Server command to complete.
		if err != nil {
			//Something went wrong, log a fatal error and exit.
			log.Fatal(err)
		} else {
			//Success, log success message.
			log.Print("RealVNC Server has been removed from the cloud")
			os.Exit(0)
		}
	} else if os.Args[1] == "join" {
		// If join argument provided, use token to join VNC Server to the cloud
		log.Print("Joining RealVNC Server to the cloud...")
		svr := os.ExpandEnv("$ProgramFiles\\RealVNC\\VNC Server\\vncserver.exe")
		cmd := exec.Command(svr, "-service", "-joinCloud", os.Args[2])
		err := cmd.Start()
		cmd.Wait() //wait for RealVNC Server command to complete.
		if err != nil {
			//Something went wrong, log a fatal error and exit.
			log.Fatal(err)
		} else {
			//Success, log success message.
			log.Print("RealVNC Server has been joined to the cloud")
			os.Exit(0)
		}
	}
}
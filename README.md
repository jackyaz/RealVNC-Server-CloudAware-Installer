# VNC Server CloudAware Installer
How-to guide to generate an MSI for VNC Server that accepts a cloud token to enable cloud connectivity at install time, as well as self-removing from the cloud when uninstalled. This was developed for use with the MSI installer for [VNC Connect](https://www.realvnc.com/en/connect/) by RealVNC.

This is an unofficial application and is not developed or supported by RealVNC.

## How do I download it?
Ready to use MSIs for VNC Server CloudAware Installer are found under [Releases](https://github.com/jackyaz/VNC-Server-CloudAware-Installer/releases)

## How do I use it?
VNC Server CloudAware Installer is used at the command line. It is called using the Windows msiexec application. A sample command (Administrator Command Prompt required) is below:
```
msiexec /i "VNC-Server-6.7.2-CloudAware.msi" CLOUDTOKEN=MYTOKEN /qn
```
where MYTOKEN is replaced by your VNC Connect cloud deployment token.

If you don't have a token, you need VNC Connect Device Access Enterprise to [generate a token](https://help.realvnc.com/hc/en-us/articles/360005474138#generating-a-cloud-join-token-0-0)

When uninstalling VNC Server it will automatically remove itself from the cloud.

VNC Server CloudAware Installer will function as a normal installer if run interactively.

## Which versions and operating systems are supported?
This application is used to install VNC Server on the computer(s) that you want to connect to and control.
*   VNC Server - this has been tested on VNC Server (RealVNC) on version 6.2.0 and above.
*   Operating system - this was developed and tested on Windows 10, but should run on the Windows versions supported by VNC Connect (Windows 7 and later).

## How does it work?
This application was created out of the need to simplify the installation of VNC Server when wanting to join it to the cloud without needing to sign in or manually run a command post-install.

There are 2 components used to create the VNC Server CloudAware Installer from the standard MSI installer from RealVNC - a Go binary, and an MSI Transform file.

### Go binary
The Go binary is embedded into the MSI. It parses input from the upstream MSI installer, and calls VNC Server to perform the required action.

It can be invoked as a standalone application (Administrator required) to modify the cloud status of the VNC Server installed on the same computer.
```
cloud.exe <join/leave> [token]
```
The source for the Go binary is available in [go-binary](https://github.com/jackyaz/VNC-Server-CloudAware-Installer/blob/master/go-binary) along with a pre-compiled binary.

You can compile the Go source yourself by installing [Go](https://golang.org/doc/install) and running the below command in the same directory as the source file
```
go build -o cloud.exe
```

### MSI Transform
The MSI Transform is used to customise the MSI installer and add custom actions to use the Go binary for install and uninstall. Transform files are created using Microsoft's [Orca](https://docs.microsoft.com/en-us/windows/win32/msi/orca-exe) application. A ready-to-use Transform is available in [msi-transform](https://github.com/jackyaz/VNC-Server-CloudAware-Installer/blob/master/msi-transform). This can be applied to a generic VNC Server MSI installer (example command below, Administrator Command Prompt required), or applied using Orca to create a transformed MSI.
```
msiexec /i "VNC Server.msi" transforms="CloudJoin.mst" /qn
```
The transform modifies 3 tables: Binary, CustomAction and InstallExecuteSequence.

#### Create the transform
Open Orca, then click File, Open, and browse to the MSI you downloaded from the VNC Connect [website](https://www.realvnc.com/download)

Click Transform, New Transform.

#### Binary
The Binary table is used to embed the Go binary in the transform. To add a binary, click Tables, Add Row. Enter Name as cloudexe. For Data, click Browse and select the compiled Go binary (cloud.exe).

![Binary table](https://github.com/jackyaz/VNC-Server-CloudAware-Installer/raw/master/msi-transform/screenshots/Binary.png)

#### CustomAction
The CustomAction table is used to add the joinCloud and leaveCloud actions. Create 2 new rows as shown below:

| Action | Type | Source | Target |
| ------ | ---- | ------ | ------ |
| leaveCloud | 3074 | cloudexe | leave |
| joinCloud | 3074 | cloudexe | join \[CLOUDTOKEN\] |

![CustomAction table](https://github.com/jackyaz/VNC-Server-CloudAware-Installer/raw/master/msi-transform/screenshots/CustomAction.png)

#### InstallExecuteSequence
The InstallExecuteSequence table is used to tell the MSI when to run our actions from CustomAction. Create 2 new rows as shown below:

| Action | Condition | Sequence |
| ------ | ---- | ------ |
| leaveCloud | REMOVE~="ALL" | 1599 |
| joinCloud | NOT Installed | 5898 |

![InstallExecuteSequence table](https://github.com/jackyaz/VNC-Server-CloudAware-Installer/raw/master/msi-transform/screenshots/InstallExecuteSequence.png)

#### Saving the Transform
Click Transform, Generate Transform, and save the MST file.

#### Creating a transformed MSI
Instead of generating a transform file, you can save the MSI with the transform already applied. To do this, open the MSI in Orca and apply the transform.

Click Tools, Options and select the Database tab. Enable "Copy embedded streams during Save As" and click OK.

Click File, Save Transformed As and save the MSI file.

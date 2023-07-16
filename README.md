# RealVNC Server CloudAware Installer
How-to guide to generate an MSI for RealVNC Server that accepts a cloud token to enable cloud connectivity at install time, as well as self-removing from the cloud when uninstalled. This was developed for use with the MSI installer for [RealVNC Connect](https://www.realvnc.com/en/connect/) by RealVNC.

This is an unofficial application and is not developed or supported by RealVNC.

## How do I download it?
Ready to use MSIs for RealVNC Server CloudAware Installer are found under [Releases](https://github.com/jackyaz/RealVNC-Server-CloudAware-Installer/releases)

## How do I use it?
RealVNC Server CloudAware Installer is used at the command line. It is called using the Windows msiexec application. A sample command (Administrator Command Prompt required) is below:
```
msiexec /i "RealVNC-Server-7.5.1-Windows-en-64bit-CloudAware.msi" CLOUDTOKEN=MYTOKEN /qn
```
where MYTOKEN is replaced by your RealVNC Connect cloud deployment token.

If you don't have a token, you can obtain one from the Deployment page of the RealVNC Connect Portal. [Click here for more information](https://help.realvnc.com/hc/en-us/articles/360005474138#generating-a-cloud-connectivity-token-0-0).

When uninstalling RealVNC Server it will automatically remove itself from the cloud.

RealVNC Server CloudAware Installer will function as a normal installer if run interactively.

## Which versions and operating systems are supported?
This application is used to install VNC Server on the computer(s) that you want to connect to and control.
*   RealVNC Server - this has been tested on RealVNC Server on version 7.5.1 and above.
*   Operating system - this was developed and tested on Windows 10, but should run on the Windows versions supported by RealVNC Connect (Windows 10 and later).

## How does it work?
This application was created out of the need to simplify the installation of RealVNC Server when wanting to join it to the cloud without needing to sign in or manually run a command post-install.

There are 2 components used to create the RealVNC Server CloudAware Installer from the standard MSI installer from RealVNC - a Go binary, and an MSI Transform file.

### Go binary
The Go binary is embedded into the MSI. It parses input from the upstream MSI installer, and calls RealVNC Server to perform the required action.

It can be invoked as a standalone application (Administrator required) to modify the cloud status of the RealVNC Server installed on the same computer.
```
realvnccloudjoin.exe <join/leave> [token]
```
The source for the Go binary is available in [go-binary](https://github.com/jackyaz/RealVNC-Server-CloudAware-Installer/blob/master/go-binary) along with a pre-compiled binary.

You can compile the Go source yourself by installing [Go](https://golang.org/doc/install) and running the below command in the same directory as the source file
```
GOOS=windows GOARCH=amd64 go build realvnccloudjoin
```

### MSI Transform
The MSI Transform is used to customise the MSI installer and add custom actions to use the Go binary for install and uninstall. Transform files are created using Microsoft's [Orca](https://docs.microsoft.com/en-us/windows/win32/msi/orca-exe) application. A ready-to-use Transform is available in [msi-transform](https://github.com/jackyaz/RealVNC-Server-CloudAware-Installer/blob/master/msi-transform). This can be applied to a generic RealVNC Server MSI installer (example command below, Administrator Command Prompt required), or applied using Orca to create a transformed MSI.
```
msiexec /i "VNC-Server-7.5.1-Windows-en-64bit.msi" transforms="CloudJoin.mst" /qn
```
The transform modifies 3 tables: Binary, CustomAction and InstallExecuteSequence.

#### Create the transform
Open Orca, then click File, Open, and browse to the MSI you downloaded from the [RealVNC website](https://downloads.realvnc.com/download/file/vnc.files/VNC-Server-Latest-Windows-msi.zip).

Click Transform, New Transform.

#### Binary
The Binary table is used to embed the Go binary in the transform. To add a binary, click Tables, Add Row. Enter Name as cloudexe. For Data, click Browse and select the compiled Go binary (realvnccloudjoin.exe).

![Binary table](https://github.com/jackyaz/RealVNC-Server-CloudAware-Installer/raw/master/msi-transform/screenshots/Binary.png)

#### CustomAction
The CustomAction table is used to add the joinCloud and leaveCloud actions. Create 2 new rows as shown below:

| Action | Type | Source | Target |
| ------ | ---- | ------ | ------ |
| leaveCloud | 3074 | cloudexe | leave |
| joinCloud | 3074 | cloudexe | join \[CLOUDTOKEN\] |

![CustomAction table](https://github.com/jackyaz/RealVNC-Server-CloudAware-Installer/raw/master/msi-transform/screenshots/CustomAction.png)

#### InstallExecuteSequence
The InstallExecuteSequence table is used to tell the MSI when to run our actions from CustomAction. Create 2 new rows as shown below:

| Action | Condition | Sequence |
| ------ | ---- | ------ |
| leaveCloud | REMOVE~="ALL" | 1599 |
| joinCloud | NOT Installed | 5898 |

![InstallExecuteSequence table](https://github.com/jackyaz/RealVNC-Server-CloudAware-Installer/raw/master/msi-transform/screenshots/InstallExecuteSequence.png)

#### Saving the Transform
Click Transform, Generate Transform, and save the MST file.

#### Creating a transformed MSI
Instead of generating a transform file, you can save the MSI with the transform already applied. To do this, open the MSI in Orca and apply the transform.

Click Tools, Options and select the Database tab. Enable "Copy embedded streams during Save As" and click OK.

Click File, Save Transformed As and save the MSI file.

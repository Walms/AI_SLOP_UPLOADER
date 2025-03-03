# NOSTROMO FILE TRANSFER SYSTEM

<div align="center">

```
 _       __  ______  __  __  _        ___    _   _  ______ 
| |     / / / ____/ / / / / | |      /   |  / \ / / /_  _/ 
| | /| / / / __/   / / / /  | |     / /| | /  _  /   / /   
| |/ |/ / / /___  / /_/ /   | |___ / ___ |/ /| \ \  / /    
|__/|__/ /_____/  \____/    |_____/_/  |_/_/ |_\_/ /_/     
                                                          
                "Building Better Worlds"
```

**A retro-futuristic file transfer utility inspired by the computers from the film "Alien"**

</div>

![Screenshot of the Nostromo File Transfer System](screenshot.jpg)

## OVERVIEW

The Nostromo File Transfer System is a standalone Go application that provides a simple, browser-based file uploading utility with a retro-futuristic interface inspired by the computer terminals from the USCSS Nostromo spacecraft in the "Alien" movie.

This lightweight tool allows quick file transfers over a local network with minimal setup, making it perfect for temporary file sharing situations when more complex solutions would be overkill.

## FEATURES

- **Retro-Futuristic UI**: Green monochrome CRT-style interface with classic computer terminal aesthetics
- **Drag & Drop Uploads**: Simple and intuitive file selection
- **Real-time Progress**: Visual progress bars and console feedback for all operations  
- **Standalone Binary**: Runs as a single executable with no dependencies
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **Zero Configuration**: Just run it and start uploading files

## USE CASES

- **Quick File Transfers**: Need to get files onto a Linux server quickly without setting up SFTP or other services
- **LAN File Sharing**: Share files between devices on a local network during a LAN party or meeting
- **Temporary Upload Solution**: When you need a temporary drop point for files in a controlled environment
- **Development Testing**: For quickly uploading test files during development
- **Air-Gapped Systems**: Transfer files to systems that aren't connected to the internet or shared networks

## ⚠️ SECURITY NOTICE

**This tool is designed for convenience, not security.**

The Nostromo File Transfer System has **NO AUTHENTICATION** and **NO ENCRYPTION**. Anyone who can access the server's IP address and port can upload files to your system.

**DO NOT USE IN PRODUCTION ENVIRONMENTS OR ON PUBLIC NETWORKS.**

Appropriate use cases are limited to:
- Temporary file transfers on trusted networks
- Development and testing environments
- Personal use in controlled settings

## BUILDING FROM SOURCE

### Prerequisites

- Go 1.16 or newer

### Standard Build

```bash
# Clone the repository
git clone https://github.com/yourusername/nostromo-transfer.git
cd nostromo-transfer

# Build the binary
go build -o nostromo-transfer
```

### Cross-Platform Builds

```bash
# For Windows
GOOS=windows GOARCH=amd64 go build -o nostromo-transfer.exe

# For macOS
GOOS=darwin GOARCH=amd64 go build -o nostromo-transfer-mac

# For Linux
GOOS=linux GOARCH=amd64 go build -o nostromo-transfer-linux
```

## CREATING A PORTABLE VERSION

To create a fully portable version that can be run from a USB drive or shared folder:

```bash
# For Windows
go build -ldflags="-H windowsgui" -o nostromo-portable.exe

# For Linux/macOS (static linking)
go build -ldflags="-linkmode external -extldflags -static" -o nostromo-portable
```

## USAGE

### Basic Usage

```bash
# Run with default settings (port 8080, saving files to current directory)
./nostromo-transfer
```

### Command Line Options

```bash
# Run on a specific port
./nostromo-transfer --port 9000

# Save files to a specific directory
./nostromo-transfer --dir ./uploads

# Combine options
./nostromo-transfer --port 7777 --dir /path/to/upload/folder
```

### Accessing the Interface

Once running, access the interface by opening a web browser and navigating to:

- Local access: http://localhost:8080 (or your specified port)
- Network access: http://YOUR_IP:8080 (as displayed in the terminal)

## TROUBLESHOOTING

### Unable to Access Server from Other Devices

- Check if a firewall is blocking the port
- Ensure you're using the correct IP address (the one displayed in the terminal)
- Try a different port with the `--port` option

### Upload Permission Errors

- Ensure the directory specified with `--dir` has write permissions
- Try running with elevated privileges if necessary (but avoid this when possible)

### Large File Upload Issues

The default maximum file size is 32MB. For larger files, consider modifying the source code:

```go
// Change the 32 << 20 (32MB) to your desired size, e.g., 128 << 20 (128MB)
if err := r.ParseMultipartForm(128 << 20); err != nil {
    // ...
}
```

## FUTURE ENHANCEMENTS

While staying true to the retro aesthetic, some potential improvements could include:

- Basic authentication option
- Download functionality
- File listing capabilities
- Configurable upload size limits
- System sounds and audio feedback

## LEGAL DISCLAIMER

This project is a fan creation and is not affiliated with, endorsed by, or in any way officially connected with the "Alien" franchise, 20th Century Studios, or any of its subsidiaries or affiliates.

The MU/TH/UR 6000, Nostromo, and Weyland-Yutani Corporation are fictional entities from the "Alien" universe created by Dan O'Bannon, Ronald Shusett, and others.

This project is intended as a creative homage and is for personal and educational use only.

## LICENSE

The Unlicense

package main


import (
	"flag"
	"fmt"
	"io"

	"log"
	"net"
	"net/http"
	"os"

	"path/filepath"
	"time"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NOSTROMO MU/TH/UR 6000 FILE SYSTEM</title>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Share+Tech+Mono&display=swap');
        
        :root {
            --bg-color: #000000;

            --terminal-color: #001100;
            --text-color: #5cdb5c;
            --accent-color: #93e293;
            --warning-color: #ff6b6b;
            --highlight-color: #98fb98;
            --grid-color: rgba(0, 59, 0, 0.3);
        }
        
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }
        
        body {
            background-color: var(--bg-color);
            color: var(--text-color);
            font-family: 'Share Tech Mono', monospace;
            font-size: 16px;
            line-height: 1.4;

            padding: 20px;
            position: relative;
            overflow-x: hidden;
            min-height: 100vh;
        }
        
        /* CRT screen effect */
        body::before {
            content: "";
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: linear-gradient(
                rgba(18, 16, 16, 0) 50%,
                rgba(0, 0, 0, 0.25) 50%
            );
            background-size: 100% 4px;
            pointer-events: none;
            z-index: 10;
        }
        
        /* Vignette effect */
        body::after {
            content: "";
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: radial-gradient(
                circle at center,
                transparent 50%,
                rgba(0, 10, 0, 0.4) 100%
            );
            pointer-events: none;
            z-index: 11;
        }
        
        .container {
            max-width: 900px;
            margin: 0 auto;
        }

        
        .screen {

            border: 8px solid #222;

            border-radius: 2px;
            background-color: var(--terminal-color);
            padding: 30px;
            box-shadow: 
                0 0 20px rgba(0, 100, 0, 0.5),
                inset 0 0 30px rgba(0, 30, 0, 0.5);
            margin-bottom: 20px;
            position: relative;
            overflow: hidden;
        }
        
        .scanline {
            width: 100%;
            height: 4px;
            background-color: rgba(0, 255, 0, 0.07);
            position: absolute;
            top: 0;
            left: 0;
            animation: scanline 8s linear infinite;
            z-index: 8;

            pointer-events: none;
        }
        
        .header {
            text-align: center;
            margin-bottom: 30px;
            border-bottom: 1px solid var(--accent-color);

            padding-bottom: 15px;
        }
        
        .company-logo {
            font-size: 14px;
            color: var(--accent-color);
            margin-bottom: 10px;
            letter-spacing: 1px;
        }
        
        h1 {
            font-size: 26px;
            font-weight: normal;
            letter-spacing: 2px;
            margin-bottom: 5px;

        }
        

        .console-line {
            opacity: 0.8;

            font-size: 14px;

            margin-bottom: 5px;
        }

        
        .system-info {
            font-size: 14px;
            margin-bottom: 20px;
            text-align: left;
            padding: 10px;

            border: 1px solid var(--accent-color);

            background-color: rgba(0, 20, 0, 0.4);

        }
        

        .cursor {
            display: inline-block;
            width: 8px;
            height: 15px;
            background: var(--accent-color);
            margin-left: 5px;
            animation: blink 1s step-end infinite;
        }
        
        .grid-container {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 20px;

            margin-bottom: 20px;
        }
        
        .drop-zone {
            border: 2px solid var(--accent-color);
            padding: 30px;
            text-align: center;
            font-size: 18px;
            margin: 20px 0;
            position: relative;
            background: rgba(0, 20, 0, 0.4);
            transition: all 0.3s;
            cursor: pointer;
            min-height: 200px;

            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
        }
        
        .drop-zone:hover, .drop-zone.highlight {
            background: rgba(0, 40, 0, 0.4);

            box-shadow: 0 0 15px rgba(92, 219, 92, 0.3);
        }
        
        .drop-zone .arrow {
            font-size: 30px;
            opacity: 0.8;
            display: block;
            margin: 10px 0;

            animation: pulse 2s infinite;
        }
        
        .file-input {
            margin: 20px 0;
        }
        
        .console-box {
            font-family: 'Share Tech Mono', monospace;
            background-color: rgba(0, 15, 0, 0.6);
            border: 1px solid var(--accent-color);
            padding: 15px;
            margin-bottom: 20px;
            height: 200px;
            overflow-y: auto;
            font-size: 14px;
        }
        
        .console-box p {
            margin: 3px 0;
            word-break: break-all;
        }
        
        .btn {
            background: rgba(0, 30, 0, 0.6);
            color: var(--accent-color);
            font-family: 'Share Tech Mono', monospace;
            font-size: 16px;
            padding: 10px 20px;
            border: 1px solid var(--accent-color);
            cursor: pointer;
            transition: all 0.3s;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        
        .btn:hover {
            background: rgba(0, 60, 0, 0.6);
            box-shadow: 0 0 10px rgba(92, 219, 92, 0.3);
        }
        
        .file-list {
            margin-top: 20px;
        }

        
        .file-item {
            background: rgba(0, 20, 0, 0.4);
            border: 1px solid var(--accent-color);
            padding: 15px;
            margin-bottom: 10px;
            display: flex;
            justify-content: space-between;
            align-items: center;

        }
        
        .file-info {

            flex-grow: 1;
        }
        
        .file-name {
            color: var(--highlight-color);
            font-size: 16px;
        }
        
        .file-size {
            opacity: 0.8;
            font-size: 14px;
        }
        
        .progress-container {
            height: 15px;

            background: rgba(0, 30, 0, 0.6);
            border: 1px solid var(--accent-color);
            width: 100%;
            margin-top: 10px;

            position: relative;
            overflow: hidden;
        }
        
        .progress-bar {
            height: 100%;
            background: linear-gradient(
                to right,
                var(--text-color),
                var(--highlight-color)

            );
            width: 0%;
            transition: width 0.2s;
            position: relative;
        }
        
        .status {
            margin-left: 20px;
            font-weight: normal;
            text-transform: uppercase;
            font-size: 14px;
        }
        
        .success {
            color: var(--highlight-color);
            animation: blink 1s infinite;
        }
        
        .error {
            color: var(--warning-color);
            animation: blink 0.5s infinite;

        }
        

        .footer {
            text-align: center;
            margin-top: 30px;
            color: var(--accent-color);
            font-size: 14px;
            padding: 10px;
            border-top: 1px solid var(--accent-color);
            opacity: 0.8;
        }
        
        @keyframes scanline {
            0% {
                top: -5%;
            }
            100% {
                top: 105%;
            }
        }
        
        @keyframes blink {
            0%, 49% {
                opacity: 1;
            }
            50%, 100% {
                opacity: 0;

            }
        }
        
        @keyframes pulse {

            0%, 100% {
                opacity: 0.5;
            }
            50% {
                opacity: 1;
            }
        }
        
        /* Boot sequence effect */
        .boot-sequence {

            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: var(--terminal-color);
            padding: 40px;
            z-index: 20;
            overflow: hidden;
            font-family: 'Share Tech Mono', monospace;
            color: var(--text-color);
            display: flex;
            flex-direction: column;
            justify-content: flex-start;
            animation: fadeOut 4s forwards;
            animation-delay: 5s;
        }
        
        @keyframes fadeOut {
            0% {

                opacity: 1;
                visibility: visible;
            }

            99% {
                opacity: 0;
                visibility: visible;
            }
            100% {
                opacity: 0;
                visibility: hidden;
            }
        }
        
        .boot-line {
            margin: 5px 0;
            white-space: nowrap;
            overflow: hidden;
            animation: typing 0.5s steps(30, end);
            animation-fill-mode: both;
        }
        
        @keyframes typing {
            from { width: 0 }
            to { width: 100% }
        }
        
        .boot-line:nth-child(1) { animation-delay: 0.2s; }
        .boot-line:nth-child(2) { animation-delay: 0.8s; }
        .boot-line:nth-child(3) { animation-delay: 1.4s; }
        .boot-line:nth-child(4) { animation-delay: 2.0s; }
        .boot-line:nth-child(5) { animation-delay: 2.6s; }

        .boot-line:nth-child(6) { animation-delay: 3.2s; }
        .boot-line:nth-child(7) { animation-delay: 3.8s; }
        .boot-line:nth-child(8) { animation-delay: 4.4s; }
        
        .wy-logo {
            text-align: center;
            margin: 20px 0;
            opacity: 0;
            animation: fadeIn 1s forwards;
            animation-delay: 4.8s;
        }
        
        @keyframes fadeIn {
            from { opacity: 0; }
            to { opacity: 1; }
        }
        
        .hexgrid {
            position: absolute;

            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background-image: 
                linear-gradient(var(--grid-color) 1px, transparent 1px),
                linear-gradient(90deg, var(--grid-color) 1px, transparent 1px);
            background-size: 30px 30px;
            opacity: 0.2;
            pointer-events: none;
        }
        
        .system-status {
            display: flex;

            justify-content: space-between;
            background: rgba(0, 20, 0, 0.4);
            border: 1px solid var(--accent-color);
            padding: 10px;
            margin-bottom: 20px;
            font-size: 14px;
        }

        
        .status-item {
            display: flex;
            align-items: center;
        }
        
        .status-indicator {
            width: 10px;
            height: 10px;
            background-color: var(--highlight-color);
            border-radius: 50%;
            margin-right: 8px;

            animation: pulse 2s infinite;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="screen">
            <div class="scanline"></div>

            <div class="hexgrid"></div>

            
            <!-- Boot sequence animation -->
            <div class="boot-sequence">

                <div class="boot-line">WEYLAND-YUTANI CORPORATION</div>

                <div class="boot-line">INITIALIZING MU/TH/UR 6000 INTERFACE...</div>

                <div class="boot-line">CHECKING SYSTEM INTEGRITY... OK</div>
                <div class="boot-line">LOADING FILE TRANSFER PROTOCOLS... OK</div>
                <div class="boot-line">ESTABLISHING DATA LINK... OK</div>
                <div class="boot-line">SECURITY CLEARANCE: LEVEL C</div>

                <div class="boot-line">SYSTEM STATUS: OPERATIONAL</div>

                <div class="boot-line">INITIALIZING NOSTROMO DATA MANAGEMENT...</div>

                
                <div class="wy-logo">
                    <pre>
 _       __  ______  __  __  _        ___    _   _  ______ 
| |     / / / ____/ / / / / | |      /   |  / \ / / /_  _/ 
| | /| / / / __/   / / / /  | |     / /| | /  _  /   / /   
| |/ |/ / / /___  / /_/ /   | |___ / ___ |/ /| \ \  / /    
|__/|__/ /_____/  \____/    |_____/_/  |_/_/ |_\_/ /_/     
                                                          
                "Building Better Worlds"
                    </pre>
                </div>
            </div>
            
            <div class="header">
                <div class="company-logo">WEYLAND-YUTANI CORPORATION</div>
                <h1>NOSTROMO DATA TRANSFER MODULE</h1>
                <div class="console-line">MU/TH/UR 6000 INTERFACE VERSION 2.1.0</div>
            </div>

            
            <div class="system-status">
                <div class="status-item">
                    <div class="status-indicator"></div>
                    SYSTEM: OPERATIONAL
                </div>
                <div class="status-item">
                    <div class="status-indicator"></div>
                    CONNECTION: SECURED
                </div>
                <div class="status-item">
                    <div class="status-indicator"></div>
                    STORAGE: 89.2TB AVAILABLE
                </div>
            </div>
            
            <div class="system-info">
                >_ TERMINAL SESSION: USR.RIPLEY.3829<br>
                >_ LOCATION: DECK C - SCIENCE DIVISION<br>
                >_ DATE: <span id="currentDate">--.--.----</span> | TIME: <span id="currentTime">--:--:--</span><br>
                >_ WARNING: ALL TRANSFERS LOGGED AND MONITORED<span class="cursor"></span>
            </div>
            
            <div class="grid-container">
                <div class="console-box" id="consoleBox">
                    <p>>_ SESSION INITIALIZED</p>
                    <p>>_ READY FOR FILE UPLOAD/DOWNLOAD</p>

                    <p>>_ AWAITING USER INPUT...</p>
                </div>
                
                <div class="drop-zone" id="dropZone">
                    <div class="arrow">↓↓↓</div>
                    TRANSFER FILES TO NOSTROMO DATABASE
                    <div class="arrow">↓↓↓</div>
                    <div>SELECT FILES OR DROP HERE</div>
                    <input type="file" id="fileInput" multiple class="file-input" />
                </div>
            </div>
            
            <div class="file-list" id="fileList">
                <!-- File items will be added dynamically -->
            </div>
        </div>
        
        <div class="footer">
            © WEYLAND-YUTANI CORP 2122 • NOSTROMO MU/TH/UR 6000 FILE SYSTEM • UNAUTHORIZED ACCESS PROHIBITED
        </div>
    </div>

    
    <script>

        document.addEventListener('DOMContentLoaded', () => {
            const dropZone = document.getElementById('dropZone');
            const fileInput = document.getElementById('fileInput');
            const fileList = document.getElementById('fileList');
            const consoleBox = document.getElementById('consoleBox');
            const currentDate = document.getElementById('currentDate');
            const currentTime = document.getElementById('currentTime');

            
            // Update time and date in futuristic format
            function updateDateTime() {
                const now = new Date();
                const day = String(now.getDate()).padStart(2, '0');
                const month = String(now.getMonth() + 1).padStart(2, '0');
                const year = now.getFullYear();
                
                const hours = String(now.getHours()).padStart(2, '0');

                const minutes = String(now.getMinutes()).padStart(2, '0');
                const seconds = String(now.getSeconds()).padStart(2, '0');
                
                currentDate.textContent = day + '.' + month + '.' + year;
                currentTime.textContent = hours + ':' + minutes + ':' + seconds;
            }
            
            setInterval(updateDateTime, 1000);
            updateDateTime();
            
            function addConsoleMessage(message) {
                const p = document.createElement('p');
                p.textContent = '>_ ' + message;
                consoleBox.appendChild(p);
                consoleBox.scrollTop = consoleBox.scrollHeight;

            }
            

            // Add boot-up messages after animation
            setTimeout(() => {
                addConsoleMessage('SYSTEM READY FOR DATA TRANSFER');
                addConsoleMessage('AWAITING FILE SELECTION...');
            }, 9000);
            

            // Prevent default drag behaviors
            ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
                dropZone.addEventListener(eventName, preventDefaults, false);
                document.body.addEventListener(eventName, preventDefaults, false);
            });

            
            // Highlight drop zone when item is dragged over it
            ['dragenter', 'dragover'].forEach(eventName => {
                dropZone.addEventListener(eventName, highlight, false);
            });
            
            ['dragleave', 'drop'].forEach(eventName => {
                dropZone.addEventListener(eventName, unhighlight, false);
            });
            
            // Handle dropped files
            dropZone.addEventListener('drop', handleDrop, false);
            
            // Handle files from input element

            fileInput.addEventListener('change', handleFiles, false);
            
            function preventDefaults(e) {
                e.preventDefault();
                e.stopPropagation();

            }
            

            function highlight() {

                dropZone.classList.add('highlight');
            }
            
            function unhighlight() {
                dropZone.classList.remove('highlight');
            }
            
            function handleDrop(e) {
                const dt = e.dataTransfer;
                const files = dt.files;
                handleFiles({ target: { files } });
            }
            
            function handleFiles(e) {

                let files = [...e.target.files];
                addConsoleMessage(files.length + ' FILE(S) SELECTED FOR TRANSFER');

                files.forEach(uploadFile);
            }
            
            function uploadFile(file) {
                addConsoleMessage('UPLOADING: ' + file.name + ' (' + formatBytes(file.size) + ')');
                
                // Create file entry in the list
                const fileItem = document.createElement('div');
                fileItem.className = 'file-item';
                
                const fileInfo = document.createElement('div');

                fileInfo.className = 'file-info';
                
                const fileName = document.createElement('div');
                fileName.className = 'file-name';
                fileName.textContent = file.name;
                
                const fileSize = document.createElement('div');
                fileSize.className = 'file-size';
                fileSize.textContent = formatBytes(file.size);

                
                const progressContainer = document.createElement('div');
                progressContainer.className = 'progress-container';
                
                const progressBar = document.createElement('div');
                progressBar.className = 'progress-bar';
                
                const statusElement = document.createElement('div');

                statusElement.className = 'status';

                statusElement.textContent = 'PROCESSING';

                
                progressContainer.appendChild(progressBar);
                fileInfo.appendChild(fileName);
                fileInfo.appendChild(fileSize);
                fileInfo.appendChild(progressContainer);
                fileItem.appendChild(fileInfo);
                fileItem.appendChild(statusElement);
                fileList.appendChild(fileItem);
                
                // Create FormData and upload the file
                const formData = new FormData();
                formData.append('file', file);
                
                const xhr = new XMLHttpRequest();
                
                // Update progress bar
                xhr.upload.addEventListener('progress', (e) => {
                    if (e.lengthComputable) {
                        const percentComplete = (e.loaded / e.total) * 100;
                        progressBar.style.width = percentComplete + '%';
                    }
                });
                
                // Handle response
                xhr.onload = function() {
                    if (xhr.status === 200) {
                        statusElement.textContent = 'COMPLETE';
                        statusElement.className = 'status success';
                        addConsoleMessage('TRANSFER COMPLETE: ' + file.name);
                    } else {
                        statusElement.textContent = 'ERROR';

                        statusElement.className = 'status error';
                        addConsoleMessage('ERROR UPLOADING: ' + file.name + ' - ' + xhr.statusText);
                    }
                };
                
                xhr.onerror = function() {
                    statusElement.textContent = 'ERROR';

                    statusElement.className = 'status error';
                    addConsoleMessage('CONNECTION FAILURE: ' + file.name);
                };
                
                xhr.open('POST', '/upload', true);
                xhr.send(formData);
            }
            
            function formatBytes(bytes) {
                if (bytes === 0) return '0 Bytes';
                const k = 1024;
                const sizes = ['Bytes', 'KB', 'MB', 'GB'];
                const i = Math.floor(Math.log(bytes) / Math.log(k));
                return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
            }
            
            // Add noise effect to the screen
            function addNoise() {
                const noise = document.createElement('div');
                noise.style.position = 'absolute';
                noise.style.top = Math.random() * 100 + '%';
                noise.style.left = Math.random() * 100 + '%';
                noise.style.width = Math.random() * 5 + 'px';
                noise.style.height = Math.random() * 1 + 'px';
                noise.style.backgroundColor = 'rgba(92, 219, 92, 0.5)';
                noise.style.zIndex = '9';
                document.querySelector('.screen').appendChild(noise);
                
                setTimeout(() => {
                    noise.remove();
                }, 100);
            }
            
            setInterval(addNoise, 200);

        });

    </script>
</body>
</html>`

func main() {
	// Parse command line flags
	port := flag.Int("port", 8080, "Port to run the server on")
	uploadDir := flag.String("dir", ".", "Directory to save uploaded files")
	flag.Parse()

	// Ensure upload directory exists
	if err := os.MkdirAll(*uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Define handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, htmlTemplate)
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {

			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

			return
		}

		// Parse multipart form data (32MB max memory)
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Get the file from the form data
		file, header, err := r.FormFile("file")
		if err != nil {

			http.Error(w, "Failed to get file: "+err.Error(), http.StatusBadRequest)
			return

		}
		defer file.Close()

		// Create a new file in the upload directory
		filename := header.Filename

		out, err := os.Create(filepath.Join(*uploadDir, filename))

		if err != nil {
			http.Error(w, "Failed to create file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer out.Close()

		// Add a small delay to simulate processing for very small files
		if header.Size < 10000 { // If less than 10KB
			time.Sleep(time.Millisecond * 500)
		}


		// Copy the file data
		_, err = io.Copy(out, file)
		if err != nil {

			http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Saved file: %s (%d bytes) to %s", filename, header.Size, *uploadDir)
		fmt.Fprintf(w, "File uploaded successfully")

	})


	// Print server information
	printServerInfo(*port, *uploadDir)


	// Start the server
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func printServerInfo(port int, uploadDir string) {
	fmt.Println("\n========================================")
	fmt.Println("NOSTROMO FILE TRANSFER SYSTEM")

	fmt.Println("WEYLAND-YUTANI CORPORATION")
	fmt.Println("----------------------------------------")
	fmt.Printf("Local access: http://localhost:%d\n", port)

	// Get network interfaces

	ifaces, err := net.Interfaces()
	if err == nil {
		for _, i := range ifaces {
			addrs, err := i.Addrs()
			if err != nil {
				continue
			}

			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}

				// Skip loopback and IPv6 addresses
				if ip.IsLoopback() || ip.To4() == nil {
					continue
				}

				fmt.Printf("Network access: http://%s:%d\n", ip.String(), port)
			}
		}
	}

	absPath, err := filepath.Abs(uploadDir)
	if err == nil {
		fmt.Printf("Files will be saved to: %s\n", absPath)
	} else {
		fmt.Printf("Files will be saved to: %s\n", uploadDir)
	}

	fmt.Println("Press Ctrl+C to stop the server")
	fmt.Println("========================================\n")
}

# Building the Application

If you wish to build the application yourself or contribute to the project, you are more than welcome! Pull requests are highly appreciated. Below are the instructions to build the application.

## 🛠 Build Prerequisites

Make sure you have:

- **Python 3.10 to 3.13**: Required for the `rembg` library.
- **`rembg`**:
  
  ```bash
  pip install "rembg[cli]"
  ```
- **Go**: Install from [here](https://golang.org/dl/).
- **Wails CLI**:
  
  ```bash
  go install github.com/wailsapp/wails/v3/cmd/wails3@latest
  ```
- **garble**: Obfuscates the release build.
  
  ```bash
  go install mvdan.cc/garble@latest
  ```
- **Bun**: Install from [here](https://bun.sh/).
- **iscc (Inno Setup Compiler)**: Download and install from [here](https://jrsoftware.org/isinfo.php) to create the `.exe` installer.
    - Add the Inno Setup installation folder to your **PATH** environment variable.

## 📝 Installation Steps

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/VdanielVPL/RMBG.git
   cd RMBG
   ```

2. **Install Dependencies**:
   - **Python Dependencies**:
     ```bash
     pip install "rembg[cli]"
     ```
   - **Go Dependencies**:
     ```bash
     go mod tidy
     ```
   - **Bun Dependencies**:
     ```bash
     cd frontend
     bun i
     ```

3. **Build the Application**:
   
   Return to the main directory:
   ```bash
   cd ..
   ```
   
   Run the build command:
   ```bash
   wails3 build; if($?) { iscc setup.iss }
   ```
   
   This command will build the application and then use `iscc` to create the `.exe` installer.

   For development:
   ```bash
   wails3 dev
   ```

5. **Run the Installer (If Release Build)**:
   After a successful build and packaging, you will find the `.exe` installer in `./bin`. Run it to install and use the application.

## 🔑 Key Files

- **`main.go`**: The Go backend entry point.
- **`frontend/`**: Contains the React-based frontend.
- **`Taskfile.yml`**: Build steps used by `wails3 build`.


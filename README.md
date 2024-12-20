<h1 align="center"><img src="./app-icon.png" width="100px" align="center">RMBG</h1>
<div align="center">
  <img alt="RMBG" src="https://github.com/user-attachments/assets/your-banner-image.png">
</div>

**RMBG** is a simple tool designed to remove backgrounds from images effortlessly on Windows 10/11. It utilizes **Wails** (Go backend) and the Python-based `rembg` CLI library for background removal. You will need **Python version 3.8 to 3.12**, along with the `rembg` library. The frontend is built with **React**, using **Bun** for building, ensuring a smooth and responsive user interface.

<div align="center">
  
[![Downloads](https://img.shields.io/github/downloads/VdanielVPL/RMBG/total.svg?style=for-the-badge)](https://github.com/VdanielVPL/RMBG/releases)  [![GitHub License](https://img.shields.io/github/license/VdanielVPL/RMBG?style=for-the-badge)](https://www.gnu.org/licenses/gpl-3.0.en.html)  [![Current version](https://img.shields.io/github/v/release/VdanielVPL/RMBG?label=version&style=for-the-badge)](https://github.com/VdanielVPL/RMBG/releases)

</div>

## Downloading the Installer

  <a href="https://github.com/VdanielVPL/RMBG/releases/latest/download/RMBG-setup.exe">
    <img height='40px' src="https://img.shields.io/badge/download-red?style=for-the-badge" alt="Download Installer">
  </a>
</br>

1. **Ensure Prerequisites are Met**:
   - Make sure you have **Python version 3.8 to 3.12** and the **`rembg` library** installed (see [Prerequisites](#prerequisites)).

2. **Run the Installer**:
   - Double-click the downloaded `.exe` file to install the application.
   - After installation, the application is ready to use.

## Prerequisites

Before using RMBG, ensure you have the following installed:

- **Python 3.8 to 3.12**: Download from [here](https://www.python.org/downloads/).
- **`rembg` Library**: Install with:
  
  ```bash
  pip install "rembg[cli]"
  ```
  Check that it’s installed correctly:
  ```bash
  rembg
  ```
  This should show the `rembg` help information.

## Features

- Remove backgrounds from images quickly and easily.
- **Crop images** before removing the background.
- Switch between **cropping view** and **background removal view**.
- Supports various image formats (PNG, JPEG, WebP).
- User-friendly interface with drag-and-drop functionality.
- Lightweight and fast.

## Usage

Once the application is installed and started:

1. **Add an Image**:
   - Click on the designated area or drag and drop an image into the window.

2. **Crop the Image**:
   - Use the cropping feature to select the desired portion of the image before removal.

3. **Switch Views**:
   - Easily toggle between cropping view and background removal view.

4. **Remove the Background**:
   - Let the application process the image with `rembg` to remove the background.

5. **Save the Processed Image**:
   - Save the result with a transparent background to your chosen location.

## How It Works

RMBG uses the `rembg` CLI tool for background removal:

- **`rembg` Integration**: The Go backend invokes the Python-based `rembg` CLI.
- **Go Backend with Wails**: Ensures efficient performance and a native desktop experience.
- **React Frontend**: Provides a responsive and modern interface.
- **Supported Formats**: Handles PNG, JPEG, and WebP.

## Building the Application

If you wish to build the application yourself, you’ll need additional tools.

### Build Prerequisites

Make sure you have:

- **Python 3.8 to 3.12**: Required for the `rembg` library.
- **`rembg`**:
  
  ```bash
  pip install "rembg[cli]"
  ```
- **Go**: Install from [here](https://golang.org/dl/).
- **Wails CLI**:
  
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```
- **Bun**: Install from [here](https://bun.sh/).
- **iscc (Inno Setup Compiler)**: Download and install from [here](https://jrsoftware.org/isinfo.php) to create the `.exe` installer.
    - Add the Inno Setup installation folder to your PATH environment variable.

### Installation Steps

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
   wails build -platform windows -ldflags "-s -w" -trimpath -tags "-B"; if($?) { iscc setup.iss }
   ```
   
   This command will build the application and then use `iscc` to create the `.exe` installer.

   For development:
   ```bash
   wails dev
   ```

4. **Run the Installer (If Release Build)**:
   After a successful build and packaging, you will find the `.exe` installer in `./build/bin`. Run it to install and use the application.

## Key Files

- **`main.go`**: The Go backend, which calls `rembg`.
- **`frontend/`**: Contains the React-based frontend.
- **`wails.json`**: Wails configuration file.

## License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your improvements.

## Acknowledgments

- [Wails](https://wails.io/) for enabling cross-platform desktop applications using Go and web technologies.
- [`rembg`](https://github.com/danielgatis/rembg) for the Python-based background removal library.

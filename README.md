<div align="center"><img src="./build/appicon.png" width="100px"></div>
<h1 align="center">RMBG</h1>

**RMBG** is a simple tool designed to remove backgrounds from images effortlessly on Windows 10/11. It utilizes **Wails** (Go backend) and the Python-based rembg CLI library for background removal. You will need **Python version 3.8 to 3.12**, along with the rembg library. The frontend is built with **React**, using **Bun** for building, ensuring a smooth and responsive user interface.

<p align="center">
  <img src="https://github.com/user-attachments/assets/151b6e26-b891-479c-b13e-1cfa915913e2" alt="Main background removal view" width="49.5%">
  <img src="https://github.com/user-attachments/assets/d1876e1b-7607-4522-b9fa-0e394dc84215" alt="Image cropping view" width="49.5%">
</p>
<p align="center">
  <img src="https://github.com/user-attachments/assets/fb241b5f-e30b-4a88-99e2-66b2c23a4554" alt="Main background removal view 2" width="49.5%">
  <img src="https://github.com/user-attachments/assets/e2badf40-4038-4398-b48a-f32b9cd74f5b" alt="Image cropping view 2" width="49.5%">
</p>

<div align="center">

[![Downloads](https://img.shields.io/github/downloads/VdanielVPL/RMBG/total.svg?style=for-the-badge)](https://github.com/VdanielVPL/RMBG/releases)  [![GitHub License](https://img.shields.io/github/license/VdanielVPL/RMBG?style=for-the-badge)](https://www.gnu.org/licenses/gpl-3.0.en.html)  [![Current version](https://img.shields.io/github/v/release/VdanielVPL/RMBG?label=version&style=for-the-badge)](https://github.com/VdanielVPL/RMBG/releases)

</div>

## 📥 Downloading the Installer

<table border="10">
  <tr>
    <td align="center" valign="middle">
      <a href="https://github.com/VdanielVPL/RMBG/releases/latest/download/RMBG-setup.exe">
        <img height="40"
             src="https://img.shields.io/badge/download-blue?style=for-the-badge"
             alt="Download Installer">
      </a>
    </td>
    <td align="center" valign="middle">
      <a href="https://github.com/VdanielVPL/RMBG/releases/latest/download/RMBG-setup.exe">
        <img height="40"
             src="https://github.com/user-attachments/assets/1ccd8012-3ec5-4458-8e63-9357d0cc7fb5"
             alt="Download Installer">
      </a>
    </td>
  </tr>
</table>

1. **Ensure Prerequisites are Met**:
   - Make sure you have **Python version 3.8 to 3.12** and the **[`rembg`](https://github.com/danielgatis/rembg) library** installed (see [Prerequisites](#-prerequisites)).

2. **Run the Installer**:
   - Double-click the downloaded `.exe` file to install the application.
   - After installation, the application is ready to use.

## 📝 Prerequisites

Before using RMBG, ensure you have the following installed:

- **Python 3.8 to 3.12**: Download from [here](https://www.python.org/downloads/).
- **[`rembg`](https://github.com/danielgatis/rembg) Library**: Install with:
  
  ```bash
  pip install "rembg[cli]"
  ```
  Check that it’s installed correctly:
  ```bash
  rembg
  ```
  This should show the rembg help information.

## ✨ Features

- **Remove backgrounds** from images quickly and easily.
- **Crop images** before removing the background.
- Switch between **cropping view** and **background removal view**.
- Supports various image formats (PNG, JPEG, WebP).
- User-friendly interface with drag-and-drop functionality.
- Windows 10/11 theme colors in app (Accent color, Mica backdrop).
- Lightweight and fast.

## 🚀 Usage

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

## 🛠️ How It Works

- **rembg Integration**: The Go backend invokes the Python-based [`rembg`](https://github.com/danielgatis/rembg) CLI for background removal.
- **Go Backend with Wails**: Ensures efficient performance and a native desktop experience.
- **React Frontend**: Provides a responsive and modern interface.
- **Supported Formats**: Handles PNG, JPEG, and WebP.

## 📜 License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your improvements. For information on building the application, see [CONTRIBUTING](CONTRIBUTING.md).

## 🙏 Acknowledgments

- [Wails](https://wails.io/) for enabling cross-platform desktop applications using Go and web technologies.
- [rembg](https://github.com/danielgatis/rembg) for the Python-based background removal library.

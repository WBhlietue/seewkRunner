This project provides Windows 10/11 environment setup and runtime tools for the [Seewk] project.


This project automatically installs development tools such as CMake, Ninja, and vcpkg, while also configuring the vcvarsall.bat script and the VCPKG_ROOT environment variable.

# Getting Started
1. Clone this repository
bat
git clone https://github.com/WBhlietue/seewkRunner.git

2. From the project root directory, run the following command in your terminal to automatically install the toolchain.
bat
seewkInit

3. If MSVC 2022 or later is not installed on your system, please install it manually before proceeding.
- Open the official [Visual Studio Downloads](https://visualstudio.microsoft.com/downloads/) page and download Visual Studio.
![MSVC](msvcInstall.png)
- Scroll down to the All Downloads section and locate Visual Studio Tools (highlighted in red in the image above).
- Locate Build Tools for Visual Studio and download it. It is highlighted in red in the image above.
- After installation, open it, select Desktop development with C++, and then install.
- During installation, make sure not to install vcpkg through the Visual Studio Installer. If it is selected by default, uncheck it.


4. After navigating to the root directory of the Seewk project, use the following command to build the project with CMake.
bat
seewk make

5. Use the following command to compile and run the project.
bat
seewk start
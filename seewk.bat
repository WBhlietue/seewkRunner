@echo off



if "%VCPKG_ROOT%"=="" (
    set "vcpkg_path="
    for /f "delims=" %%i in ('where vcpkg 2^>nul') do (
        set "vcpkg_path=%%i"
        goto :found
    )
    :found
    if defined vcpkg_path (
        for %%a in ("%vcpkg_path%") do set "VCPKG_ROOT=%%~dpa"
    ) else (
        set VCPKG_ROOT=%~dp0vcpkg
    )
)

if "%IS_INIT%"=="" (
    setlocal enabledelayedexpansion
    chcp 65001
    set "initURL=C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Auxiliary\Build\vcvarsall.bat"
    call "!initURL!" x64
    
    if "%1"=="start" goto startProj
    if "%1"=="make" goto make
    if "%1"=="run" goto run
    if "%1"=="compile" goto compile
    cmd /k
    endlocal
)

if "%1"=="start" goto startProj
if "%1"=="run" goto run
if "%1"=="make" goto make
if "%1"=="compile" goto compile

echo command not found
goto :EOF

:init 
if "%IS_INIT%" == "" (
    setlocal enabledelayedexpansion
    set IS_INIT=1
    cmd /k
    endlocal
)
goto :EOF 



:startProj
ninja -C build && build\Test.exe
goto init

:compile
ninja -C build
goto init



:run
goto init


:make
call cmake -B build -S . -G Ninja
@REM call cmake -B build -S . -DVCPKG_TARGET_TRIPLET=x64-windows-static -DCMAKE_BUILD_TYPE=Release -G Ninja
goto init



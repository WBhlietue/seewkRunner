@echo off

@REM set "SEEWK_DIR=%~dp0"
@REM cd /d "%SEEWK_DIR%"
@REM set VCPKG_ROOT="%~dp0vcpkg\scripts\buildsystems\vcpkg.cmake"
if "%VCPKG_ROOT%"=="" set VCPKG_ROOT=%~dp0vcpkg

if "%IS_INIT%"=="" (
    setlocal enabledelayedexpansion
    set "initURL=C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Auxiliary\Build\vcvarsall.bat"
    call "!initURL!" x64
    set IS_INIT=1
    if "%1"=="start" goto startProjInit
    if "%1"=="make" goto makeInit
    if "%1"=="run" goto runInit
    if "%1"=="compile" goto compileInit
    cmd /k
    endlocal
)
chcp 65001
if "%1"=="start" goto startProj
if "%1"=="run" goto run
if "%1"=="make" goto make
if "%1"=="compile" goto compile
goto :EOF

:startProjInit
ninja -C build && build\Test.exe
cmd /k
goto :EOF

:makeInit
call cmake -B build -S . -DVCPKG_TARGET_TRIPLET=x64-windows-static -DCMAKE_BUILD_TYPE=Release -G Ninja
cmd /k
goto :EOF

:startProj
ninja -C build && build\Test.exe
goto :EOF

:compile
ninja -C build
goto :EOF
:compileInit
ninja -C build
cmd /k
goto :EOF


:run
build\Test.exe 
goto :EOF
:runInit
build\Test.exe 
cmd /k
goto :EOF

:make

call cmake -B build -S . -DVCPKG_TARGET_TRIPLET=x64-windows-static -DCMAKE_BUILD_TYPE=Release -G Ninja
goto :EOF


@echo off

@REM set "SEEWK_DIR=%~dp0"
@REM cd /d "%SEEWK_DIR%"
set vclocation="%~dp0vcpkg\scripts\buildsystems\vcpkg.cmake"
if "%IS_INIT%"=="" (
    setlocal enabledelayedexpansion
    set "initURL=C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Auxiliary\Build\vcvarsall.bat"
    call "!initURL!" x64
    set IS_INIT=1
    if "%1"=="start" goto startProjInit
    if "%1"=="make" goto makeInit
    cmd /k
    endlocal
)
chcp 65001
if "%1"=="start" goto startProj
if "%1"=="make" goto make
goto :EOF

:startProjInit
ninja -C build && build\Test.exe
cmd /k
goto :EOF

:makeInit
call cmake -B build -S . -DCMAKE_TOOLCHAIN_FILE=%vclocation% -DVCPKG_TARGET_TRIPLET=x64-windows-static -DCMAKE_BUILD_TYPE=Release -G Ninja
cmd /k
goto :EOF

:startProj
ninja -C build && build\Test.exe
goto :EOF

:make

call cmake -B build -S . -DCMAKE_TOOLCHAIN_FILE=%vclocation% -DVCPKG_TARGET_TRIPLET=x64-windows-static -DCMAKE_BUILD_TYPE=Release -G Ninja
goto :EOF


@echo off

@REM set "SEEWK_DIR=%~dp0"
@REM cd /d "%SEEWK_DIR%"
set vclocation="%~dp0vcpkg\scripts\buildsystems\vcpkg.cmake"
if "%IS_INIT%"=="" (
    setlocal enabledelayedexpansion
    set "initURL=C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\VC\Auxiliary\Build\vcvarsall.bat"
    call "!initURL!" x64
    set IS_INIT=1
    cmd /k
    endlocal
    goto :EOF
)

if "%1"=="start" goto startProj
if "%1"=="make" goto make

:startProj
ninja -C build && build\Test.exe
goto :EOF

:make

call cmake -B build -S . -DCMAKE_TOOLCHAIN_FILE=%vclocation% -DVCPKG_TARGET_TRIPLET=x64-windows -DCMAKE_BUILD_TYPE=Release -G Ninja
goto :EOF


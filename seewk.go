package main

import (
    "archive/zip"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "strings"	
)

func Unzip(zipPath, destDir string) error {
    destAbs, err := filepath.Abs(destDir)
    if err != nil {
        return fmt.Errorf("get absolute path failed: %w", err)
    }

    r, err := zip.OpenReader(zipPath)
    if err != nil {
        return fmt.Errorf("open zip error: %w", err)
    }
    defer r.Close()

    if err := os.MkdirAll(destAbs, 0755); err != nil {
        return fmt.Errorf("create folder error: %w", err)
    }

    for _, f := range r.File {
        destPath := filepath.Join(destAbs, f.Name)

        if !strings.HasPrefix(destPath, destAbs+string(os.PathSeparator)) {
            return fmt.Errorf("invalid file path: %s", f.Name)
        }

        if f.FileInfo().IsDir() {
            os.MkdirAll(destPath, f.Mode())
            continue
        }

        if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
            return fmt.Errorf("create folder error: %w", err)
        }

        srcFile, err := f.Open()
        if err != nil {
            return fmt.Errorf("open extract failed: %w", err)
        }

        dstFile, err := os.Create(destPath)
        if err != nil {
            srcFile.Close()
            return fmt.Errorf("create destination failed: %w", err)
        }

        _, err = io.Copy(dstFile, srcFile)
        srcFile.Close()
        dstFile.Close()

        if err != nil {
            return fmt.Errorf("write error: %w", err)
        }
    }

    return nil
}


func commandExists(cmdName string) bool {
    var cmd *exec.Cmd

    if runtime.GOOS == "windows" {
        cmd = exec.Command("where", cmdName)
    } else {
        cmd = exec.Command("which", cmdName)
    }

    err := cmd.Run()
    return err == nil
}

func isDirExists(path string) bool {
    info, err := os.Stat(path)
    if os.IsNotExist(err) {
        return false
    }
    return err == nil && info.IsDir()
}

func addToCurrentSession(dirToAdd string) {
    currentPath := os.Getenv("PATH")
    if strings.Contains(currentPath, dirToAdd) {
        return
    }
    newPath := currentPath + ";" + dirToAdd
    os.Setenv("PATH", newPath)
}
func AddToUserPath(dirToAdd string) error {
    ps := fmt.Sprintf(`$p=[Environment]::GetEnvironmentVariable("PATH","User"); if($p -notlike "*%s*"){[Environment]::SetEnvironmentVariable("PATH",($p+";%s"),"User")}`, dirToAdd, dirToAdd)
    return exec.Command("powershell", "-Command", ps).Run()
}


func installReq() {
	absPath, err := filepath.Abs("./")
	if err != nil {
		fmt.Printf("      * Failed to get absolute path: %v\n", err)
		return
	}
	AddToUserPath(absPath)
    fmt.Println("Initializing project...")

    
    if(commandExists("vcpkg")){
        fmt.Println("  > vcpkg found")
    }else{
        fmt.Println("  * vcpkg not found")
        pathExists := isDirExists("./vcpkg")
        if pathExists {
            fmt.Println("      > CMake found in ./vcpkg")
            absPath, err := filepath.Abs("./vcpkg")
            if err != nil {
                fmt.Printf("      * Failed to get absolute path: %v\n", err)
                return
            }
            AddToUserPath(absPath)
        } else {
            fmt.Println("      > vcpkg installing")
            exec.Command("cmd", "/c", 
                "git clone https://github.com/microsoft/vcpkg.git && " +
                "cd vcpkg && " +
                "bootstrap-vcpkg.bat").Run()
            
            absPath, err := filepath.Abs("./vcpkg")
            if err != nil {
                fmt.Printf("      * Failed to get absolute path: %v\n", err)
                return
            }
            AddToUserPath(absPath)
            fmt.Println("      > vcpkg installed")
        }
    }
    if commandExists("ninja") {
        fmt.Println("  > Ninja found")
    } else {
        fmt.Println("  * Ninja not found")
        fmt.Println("      Checking Ninja...")
        lastVersion, err := getLatestVersion("ninja-build", "ninja")
        if err != nil {
            fmt.Printf("      * Failed to get version: %v\n", err)
            return
        }
        fmt.Printf("      > Downloading Ninja version %s\n", lastVersion)
        err = downloadFile(fmt.Sprintf("https://github.com/ninja-build/ninja/releases/download/%s/ninja-win.zip", lastVersion), "ninja.zip")
        if err != nil {
            fmt.Printf("      * Download failed: %v\n", err)
            return
        }
        fmt.Println("      > Extracting Ninja")
        err = Unzip("ninja.zip", ".")
        if err != nil {
            fmt.Printf("      * Extract failed: %v\n", err)
            return
        }
        err = os.Remove("ninja.zip")
        if err != nil {
            fmt.Printf("       failed to delete ninja.zip: %v\n", err)
        } else {
            fmt.Println("      > Cleaned up ninja.zip")
        }
        fmt.Println("  > Ninja installed")
    }

    // CMake installation
    if commandExists("cmake") {
        fmt.Println("  > CMake found")
    } else {
        fmt.Println("  * CMake not found")
        pathExists := isDirExists("./cmake/bin")
        if pathExists {
            fmt.Println("      > CMake found in ./cmake/bin")
            absPath, err := filepath.Abs("./cmake/bin")
            if err != nil {
                fmt.Printf("      * Failed to get absolute path: %v\n", err)
                return
            }
            AddToUserPath(absPath)
        } else {
            fmt.Println("      Checking CMake...")
            lastVersion, err := getLatestVersion("Kitware", "CMake")
            if err != nil {
                fmt.Printf("      * Failed to get version: %v\n", err)
                return
            }
            fmt.Printf("      > Downloading CMake version %s\n", lastVersion)
            versionNum := strings.TrimPrefix(lastVersion, "v")
            err = downloadFile(fmt.Sprintf("https://github.com/Kitware/CMake/releases/download/%s/cmake-%s-windows-x86_64.zip", lastVersion, versionNum), "cmake.zip")
            if err != nil {
                fmt.Printf("      * Download failed: %v\n", err)
                return
            }
            fmt.Println("      > Extracting CMake")
            err = Unzip("cmake.zip", ".")
            if err != nil {
                fmt.Printf("      * Extract failed: %v\n", err)
                return
            }
            os.Rename("cmake-"+versionNum+"-windows-x86_64", "cmake")
            err = os.Remove("cmake.zip")
            if err != nil {
                fmt.Printf("      failed to delete cmake.zip: %v\n", err)
            } else {
                fmt.Println("      > Cleaned up cmake.zip")
            }
            fmt.Println("  > CMake installed")

            absPath, err := filepath.Abs("./cmake/bin")
            if err != nil {
                fmt.Printf("      * Failed to get absolute path: %v\n", err)
                return
            }
            AddToUserPath(absPath)
        }
        fmt.Println("~ Restart your terminal for changes to take effect")
    }
}



func getLatestVersion(owner, repo string) (string, error) {
    url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

    resp, err := http.Get(url)
    if err != nil {
        return "", fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("API returned status: %s", resp.Status)
    }

    var result struct {
        TagName string `json:"tag_name"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", fmt.Errorf("JSON parse failed: %w", err)
    }

    return result.TagName, nil
}

func downloadFile(url, filepath string) error {
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("bad status: %s", resp.Status)
    }

    out, err := os.Create(filepath)
    if err != nil {
        return fmt.Errorf("create file failed: %w", err)
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return fmt.Errorf("write file failed: %w", err)
    }

    return nil
}

type Config struct {
    VCVarsURL string `json:"vcvars_url"`
}
func initENV() string {
    configPath := filepath.Join(".", "conf.json")
    data, err := os.ReadFile(configPath)
    if err != nil {
        fmt.Printf("Warning: Failed to read config file: %v\n", err)
        return "-"
    }
    
    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        fmt.Printf("Warning: Failed to parse config: %v\n", err)
        return "-"
    }
    
    return config.VCVarsURL

}
func setInit(url string) {
    config := struct {
        VCVarsURL string `json:"vcvars_url"`
    }{VCVarsURL: url}
    
    data, _ := json.MarshalIndent(config, "", "  ")
    os.WriteFile("conf.json", data, 0644)
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Run 'seewk' for usage")
        return
    }
    switch os.Args[1] {
    case "install":
        installReq()
    case "init":
        url := initENV();
        fmt.Println(url)
    case "setInit":
        setInit(os.Args[2]);
    default:
        fmt.Printf("Unknown command: %s\n", os.Args[1])
        fmt.Println("Run 'seewk' for usage")
    }
}
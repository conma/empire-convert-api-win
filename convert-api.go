package main

import (
    "fmt"
    "io"
    "log"
    "mime/multipart"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
    "regexp"
)

func main() {
    http.HandleFunc("/convert", convertHandler)
    fmt.Println("Server started at http://10.10.100.3:8080")
    log.Fatal(http.ListenAndServe("10.10.100.3:8080", nil))
}

// convertHandler handles file upload and conversion
func convertHandler(w http.ResponseWriter, r *http.Request) {
    // Limit the size of the request to prevent large files
    r.ParseMultipartForm(10 << 20) // 10 MB max file size

    // Get the file from the request
    file, handler, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Failed to upload file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    target := r.FormValue("path")
    re := regexp.MustCompile(`\d{4}[\\/]\d{2}[\\/]\d{2}[\\/][a-zA-Z0-9]+`)

    // Find the matching substring for each path
    target = re.FindString(target)

    folder, err := os.Getwd()
    fmt.Println("current-", folder)

    // Save the uploaded file to a temporary location
    saveTo := filepath.Join(folder, "uploads", target)
    inputPath := saveTo

    outputName := strings.Replace(handler.Filename, ".wmf", ".png", 1)
    outputName = strings.Replace(outputName, ".emf", ".png", 1)
    outputPath := filepath.Join(saveTo, outputName)
    if err := saveUploadedFile(file, inputPath, handler.Filename); err != nil {
        http.Error(w, "Failed to save file", http.StatusInternalServerError)
        return
    }

    // Call magick.exe to convert the file to PNG
    if err := convertToPNG(filepath.Join(saveTo, handler.Filename), outputPath); err != nil {
        http.Error(w, "Failed to convert file"+err.Error(), http.StatusInternalServerError)
        return
    }

    // Open the converted PNG file
    convertedFile, err := os.Open(outputPath)
    if err != nil {
        http.Error(w, "Failed to open converted file", http.StatusInternalServerError)
        return
    }
    defer convertedFile.Close()

    // Set response headers for file download
    w.Header().Set("Content-Disposition", "attachment; filename="+outputName)
    w.Header().Set("Content-Type", "image/png")

    // Write the PNG file to the response
    if _, err := io.Copy(w, convertedFile); err != nil {
        http.Error(w, "Failed to send converted file", http.StatusInternalServerError)
    }

    // Clean up temporary files
    // os.Remove(inputPath)
    // os.Remove(outputPath)
}

// saveUploadedFile saves the uploaded file to a temporary location
func saveUploadedFile(file multipart.File, path, filename string) error {
    err := os.MkdirAll(path, 0755)

    out, err := os.Create(filepath.Join(path, filename))
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, file)
    return err
}

// convertToPNG runs magick.exe to convert the file to PNG format
func convertToPNG(inputPath, outputPath string) error {
    // Run the command: magick.exe inputPath outputPath
    cmd := exec.Command("magick", inputPath, outputPath)
    fmt.Println(time.Now(), "run: ", cmd.String())
    return cmd.Run()
}

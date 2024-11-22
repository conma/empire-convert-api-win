package docx

import (
	"bytes"
	"docxprocessing/constant"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const (
	url = constant.WMF_API
)

func Post(folder, inputname, outputname string) error {
	err := postFile(folder, inputname, outputname)
	return err
}

// postFile sends a file as a multipart form POST request
func postFile(folder, inputname, outputname string) error {
	filePath := filepath.Join(folder, inputname)
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Create a new buffer to hold the form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	writer.WriteField("path", folder)

	// Create a form field for the file
	fileWriter, err := writer.CreateFormFile("file", filepath.Base(inputname))
	if err != nil {
		return fmt.Errorf("could not create form file: %v", err)
	}

	// Copy the file contents into the form field
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return fmt.Errorf("could not copy file contents: %v", err)
	}

	// Close the writer to set the Content-Type with boundary
	writer.Close()

	// Create a new HTTP POST request with the form data
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return fmt.Errorf("could not create POST request: %v", err)
	}

	// Set the Content-Type header to multipart/form-data with the boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Read the response body (for demonstration, write it to a file)
	outputPath := filepath.Join(folder, outputname)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("could not create output file: %v", err)
	}
	defer outputFile.Close()

	// Save the response body to the output file
	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		return fmt.Errorf("could not save output file: %v", err)
	}

	return nil
}

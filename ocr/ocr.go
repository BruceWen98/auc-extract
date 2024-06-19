package ocr

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type Ocr struct {
	Frame      int
	File       string
	InputPath  string
	OutputPath string
}

func NewOCR(frame int, file, inputPath, outputPath string) *Ocr {
	return &Ocr{
		Frame:      frame,
		File:       file,
		InputPath:  inputPath,
		OutputPath: outputPath,
	}
}

// // extract from bottom height 50
// func (o *Ocr) Crop(height int) (string, error) {
// 	file, err := os.Open(o.InputPath + "/" + o.File)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer file.Close()
// 	img, _, err := image.Decode(file)
// 	if err != nil {
// 		return "", err
// 	}
// 	rect := img.Bounds()
// 	crop := image.Rect(rect.Min.X, rect.Max.Y-height, rect.Max.X, rect.Max.Y)
// 	subImg := img.(*image.YCbCr).SubImage(crop) //jpeg

// 	base := strings.Split(filepath.Base(o.File), ".")[0]
// 	ext := filepath.Ext(o.File)
// 	cropFilename := fmt.Sprintf(o.OutputPath+"/%s_crop%s", base, ext)
// 	croppedFile, err := os.Create(cropFilename)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer croppedFile.Close()
// 	err = jpeg.Encode(croppedFile, subImg, nil)
// 	if err != nil {
// 		return "", err
// 	}
// 	return cropFilename, nil
// }

// func (o *Ocr) Extract(file string) (string, error) {
// 	// return o.ExtractWithAPI(file)
// 	return o.ExtractWithCommand(file)
// }

// extract text from image through gosseract API.
// func (o *Ocr) ExtractWithAPI(file string) (string, error) {
// 	client := gosseract.NewClient()
// 	defer client.Close()
// 	client.SetImage(file)
// 	text, err := client.Text()
// 	if err != nil {
// 		return "", err
// 	}
// 	return text, nil
// }

func runCmd(cmd string, args ...string) (string, error) {
	// log.Printf("running command: %s %s", cmd, strings.Join(args, " "))
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run command: %w, out=%s", err, out)
	}
	if exitError, ok := err.(*exec.ExitError); ok {
		exitCode := exitError.ExitCode()
		if exitCode != 0 {
			return "", fmt.Errorf("command exited with non-zero exit code: %d, output=%s", exitCode, out)
		}
	}

	return strings.TrimSpace(string(out)), nil
}

func (o *Ocr) Extract(file string, psm string) (string, error) {
	//grayscale
	// convert frame98040.jpg -set colorspace Gray -separate -average frame98040.graytscale.jpg
	grayScaleFile := strings.Replace(file, ".jpg", ".grayscale.jpg", 1)
	// -set colorspace Gray -separate -average
	_, err := runCmd("convert", file, "-set", "colorspace", "Gray", "-separate", "-average", grayScaleFile)
	if err != nil {
		return "", fmt.Errorf("failed to run convert to grayscale: %w", err)
	}
	dpi300File := strings.Replace(grayScaleFile, ".jpg", ".dpi300.jpg", 1)
	// -units PixelsPerInch -density 300
	_, err = runCmd("convert", grayScaleFile, "-units", "PixelsPerInch", "-density", "300", dpi300File)
	if err != nil {
		return "", fmt.Errorf("failed to run convert to increase dpi: %w", err)
	}
	// tesseract x.jpg - --oem 1 quiet --psm 6 //psm to keep as a whole line
	// remove "quiet" otherwise the psm 6 no effect
	args := []string{dpi300File, "-", "--oem", "1", "-l", "eng+chi_sim+chi_tra"}
	if psm != "-1" {
		args = append(args, "--psm", psm)
	}
	out, err := runCmd("tesseract", args...)
	if err != nil {
		return "", fmt.Errorf("failed to run tesseract: %w", err)
	}

	return strings.TrimSpace(string(out)), nil
}

// crop and extract in one call
func (o *Ocr) Process(cropSizeOffset string, psm string) (string, error) {
	base := strings.Split(filepath.Base(o.File), ".")[0]
	ext := filepath.Ext(o.File)
	cropFilename := fmt.Sprintf(o.OutputPath+"/%s.crop%s", base, ext)

	_, err := runCmd("convert", o.InputPath+"/"+o.File, "-gravity", "South", "-crop", cropSizeOffset, cropFilename)
	if err != nil {
		return "", fmt.Errorf("failed to crop the image: %w", err)
	}
	text, err := o.Extract(cropFilename, psm)
	if err != nil {
		return "", fmt.Errorf("failed to extract text from the image: %w", err)
	}

	return text, nil
}

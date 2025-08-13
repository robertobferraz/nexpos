package utils

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/nfnt/resize"
)

func ValidationBearerToken(token *string) (*string, error) {
	if len(*token) <= 0 {
		return nil, errors.New("invalid Token")
	}

	x := strings.Split(*token, "Bearer ")
	if len(x) < 2 {
		return nil, errors.New("invalid Token")
	}

	return &x[1], nil
}

func CreateRandomUsername(name *string) *string {
	var usernameBuffer bytes.Buffer
	usernameBuffer.WriteString(strings.ToLower(strings.ReplaceAll(*name, " ", "_")))
	usernameBuffer.WriteString(fmt.Sprintf(":%v", rand.Int31()))

	return PString(usernameBuffer.String())
}

func UploadImage(width, height, finalSize *int, initialSize *int64, file *multipart.FileHeader) ([]byte, *string, *string, error) {
	var imageName *string
	if file.Size > *initialSize {
		return nil, nil, nil, errors.New("file size too big")
	}

	fileHeader, err := file.Open()
	if err != nil {
		return nil, nil, nil, err
	}

	defer func(fileHeader multipart.File) {
		err := fileHeader.Close()
		if err != nil {

		}
	}(fileHeader)

	imageData, err := io.ReadAll(fileHeader)
	if err != nil {
		return nil, nil, nil, err
	}

	contentType := http.DetectContentType(imageData)
	if !strings.HasPrefix(contentType, "image/") {
		return nil, nil, nil, errors.New("invalid content type")
	}

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, nil, nil, err
	}

	bounds := img.Bounds()
	if bounds.Dx() != *width || bounds.Dy() != *height {
		resizedImg := resize.Resize(uint(*width), uint(*height), img, resize.Lanczos3)

		var buf bytes.Buffer
		err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: 85})
		if err != nil {
			return nil, nil, nil, err
		}
		imageData = buf.Bytes()

		if len(imageData) > *finalSize {
			return nil, nil, nil, errors.New("image too big")
		}

		imageName = PString(file.Filename)
		if !strings.HasSuffix(strings.ToLower(*imageName), ".jpg") && !strings.HasSuffix(strings.ToLower(*imageName), ".jpeg") {
			imageName = PString(strings.TrimSuffix(*imageName, strings.ToLower(filepath.Ext(*imageName))) + ".jpg")
		}
	}

	return imageData, imageName, &contentType, nil
}

var (
	rateData = make(map[string]rateEntry)
	mu       sync.Mutex
)

type rateEntry struct {
	count     *int
	expiresAt *time.Time
}

func SimpleRateLimit(keyParts ...string) *bool {
	key := ""
	for _, part := range keyParts {
		key += part + "|"
	}

	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	if entry, exists := rateData[key]; exists {
		if now.After(*entry.expiresAt) {
			rateData[key] = rateEntry{count: PInt(1), expiresAt: PTime(now.Add(10 * time.Second))}
			return PBool(false)
		}

		*entry.count++
		rateData[key] = entry

		if *entry.count > 5 {
			return PBool(true)
		}
		return PBool(false)
	}

	rateData[key] = rateEntry{count: PInt(1), expiresAt: PTime(now.Add(10 * time.Second))}
	return PBool(false)
}

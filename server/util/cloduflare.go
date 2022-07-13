package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	c "github.com/angelorc/sinfonia-go/config"
	"github.com/angelorc/sinfonia-go/mongo/model"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

func UploadImage(cfg c.CloudflareConfig, data model.MerkledropUpdateReq) (*string, error) {
	var imageUrl *string

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	body := &bytes.Buffer{}
	bodywriter := multipart.NewWriter(body)

	writer, err := bodywriter.CreateFormFile("file", data.Name)
	if err != nil {
		return imageUrl, err
	}

	_, err = io.Copy(writer, data.Image.File)
	if err != nil {
		return imageUrl, err
	}

	err = bodywriter.Close()
	if err != nil {
		return imageUrl, err
	}

	cloudFlareImagesUrl := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/images/v1", cfg.Account)
	req, err := http.NewRequest("POST", cloudFlareImagesUrl, bytes.NewReader(body.Bytes()))
	if err != nil {
		return imageUrl, err
	}

	req.Header.Set("Content-Type", bodywriter.FormDataContentType())
	req.Header.Add("Authorization", "Bearer "+cfg.Images)
	rsp, _ := client.Do(req)
	if rsp.StatusCode != http.StatusOK {
		return imageUrl, fmt.Errorf("request failed with response code: %d", rsp.StatusCode)
	}
	defer rsp.Body.Close()
	rspBz, _ := ioutil.ReadAll(rsp.Body)

	var cloudlfareResp model.MerkledropUpdateImageResponse
	if err := json.Unmarshal(rspBz, &cloudlfareResp); err != nil {
		return imageUrl, err
	}

	if cloudlfareResp.Success {
		if len(cloudlfareResp.Result.Variants) > 0 {
			imageUrl = &cloudlfareResp.Result.Variants[0]
		}
	}

	return imageUrl, nil
}

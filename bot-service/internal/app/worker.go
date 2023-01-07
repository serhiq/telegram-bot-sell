package app

import (
	"bot/config"
	"bot/internal/entity"
	"fmt"
	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func SyncMenu(a *An) error {

	var response entity.Menu

	endpoint := a.Cfg.BaseUrl + "/product/" + a.Cfg.Store
	resp, err := a.Client.R().
		SetHeader("Authorization", a.Cfg.Auth).
		SetResult(&response).
		Get(endpoint)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		log.Println("sync: error" + resp.Status())
		return nil
	}

	result := []*entity.MenuItemDatabase{}

	for _, item := range response {

		if !item.CanAddToOrder() {
			continue
		}

		previewImage := ""
		if item.Image != "" {
			tmpImage := filepath.Join(config.TempPatch, item.Image)
			previewImage = filepath.Join(config.PreviewCachePatch, item.Image)
			err = Download(tmpImage, item.ImageUrl)
			if err != nil {
				log.Errorf("download: %s", err)
			}

			err = resizeTmpFile(tmpImage, previewImage)
			if err != nil {
				log.Errorf("resize image: %s", err)
			}
			_ = os.Remove(tmpImage)
		}

		result = append(result, &entity.MenuItemDatabase{
			Name: item.Name,
			//StoreID:    item.StoreID,
			UserID:      "",
			UUID:        item.UUID,
			ParentUUID:  item.ParentUUID,
			Group:       item.Group,
			Image:       previewImage,
			MeasureName: item.MeasureName,
			Price:       item.Price,
		})
	}

	err = a.Db.ImportMenu(result)
	return nil
}

func resizeTmpFile(tmpDest string, filename string) error {
	file, err := os.Open(tmpDest)

	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}
	file.Close()

	m := resize.Resize(0, 300, img, resize.Lanczos3)

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	jpeg.Encode(out, m, nil)

	return nil
}

func Download(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func PostOrder(a *An, order *entity.OrderRequest) (*entity.PostOrderResponse, error) {
	order.State = "new"
	order.ID = entity.GetRandomOrderNumber()

	endpoint := a.Cfg.BaseUrl + "/order/" + a.Cfg.Store

	result := entity.PostOrderResponse{}

	log.Print("заказ")
	log.Print(order)

	response, err := a.Client.R().
		EnableTrace().
		SetBody(order).
		SetResult(&result).
		SetHeader("Authorization", a.Cfg.Auth).
		Post(endpoint)

	if response.StatusCode() != 200 {
		log.Print("Ответ")
		log.Print(response)
		return nil, fmt.Errorf("failed post order: %d", response.StatusCode())
	}

	return &result, err
}

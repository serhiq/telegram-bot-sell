package app

import (
	"bot/config"
	"bot/internal/entity"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"
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
		log.Println("ошибка: " + resp.Status())
		return nil
	}

	result := []*entity.MenuItemDatabase{}

	for _, item := range response {

		if !item.CanAddToOrder() {
			continue
		}

		previewImage := ""
		if item.Image != "" {
			fakeImageUrl := "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f7/Lemon_-_whole_and_split.jpg/1280px-Lemon_-_whole_and_split.jpg"
			// вообще есть прекрасные методы `os.TempDir()` и `os.CreateTemp()`, чтобы не изобретать свой tmp
			tmpImage := filepath.Join(config.TempPatch, item.Image)
			previewImage = filepath.Join(config.PreviewCachePatch, item.Image)
			err = Download(tmpImage, fakeImageUrl)
			if err != nil {
				log.Errorf("download: %s", err)
				// опять же, если не удалось скачать изображение, то зачем продолжать?
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

// я бы вместо `filepath` передавал `io.Writer` - решение гибче получится
func Download(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer out.Close()

	// никогда не используйте дефолтный HttpClient - если заглянете в код, то обнаружите, что у него бесконечные таймауты
	// то есть все может прекрасно так зависнуть
	// и в целом странно, используя `resty` еще использовать и `http`
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	// просто закрыть тело недостаточно - его обязательно надо считать до конца даже в случае ошибки, иначе начнутся утечки памяти
	// можно считывать в `io.Discard`
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

	// не совсем верно проверять только на 200, есть прекрасный метод `response.IsSuccess()`
	// фактически успеход является любой код от 200 до 299, а если возможны перенаправления, то и до 399
	if response.StatusCode() != 200 {
		log.Print("Ответ")
		log.Print(response)
		return nil, fmt.Errorf("failed post order: %d", response.StatusCode())

	}

	return &result, err
}

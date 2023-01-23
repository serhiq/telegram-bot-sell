package worker

import (
	"bot/config"
	c "bot/pkg/restoClient"
	repositoryProduct "bot/services/bot/pkg/repository/product"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/go-resty/resty/v2"
	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"
)

type SyncWorker struct {
	r         *repositoryProduct.Repository
	evoClient *c.RestoClient
	client    *resty.Client
}

func (s SyncWorker) Start() error {

	resp, err := s.evoClient.GetMenu()
	if err != nil {
		return err
	}

	clearOldPreviews()

	result := []*repositoryProduct.Product{}

	for _, product := range resp {

		if !product.CanAddToOrder() {
			continue
		}

		menuItem := &repositoryProduct.Product{
			Name: product.Name,
			//StoreID:    item.StoreID,
			UUID:        product.UUID,
			ParentUUID:  product.ParentUUID,
			Group:       product.Group,
			Image:       "",
			MeasureName: product.MeasureName,
			Price:       product.Price,
		}

		if product.ImageNotEmpty() {
			thumbnail, e := createThumbnail(s.client, product.Image, product.ImageUrl)
			if e != nil {
				log.Printf("sync: err create trumbnail, %s", err)
				continue
			}
			menuItem.Image = thumbnail
		}
		result = append(result, menuItem)
	}

	err = s.r.ImportMenu(result)
	if err != nil {
		return fmt.Errorf("sync: err, %s", err)
	}
	return nil
}

func createThumbnail(r *resty.Client, fileName string, url string) (string, error) {
	tmpImage := filepath.Join(config.TempPatch, fileName)
	defer os.Remove(tmpImage)

	previewImage := filepath.Join(config.PreviewCachePatch, fileName)

	resp, err := r.R().
		SetOutput(tmpImage).
		Get(url)

	if err != nil {
		return "", fmt.Errorf("sync: error get image, %s", err)

	}

	if !resp.IsSuccess() {
		return "", fmt.Errorf("sync: error get image, code == %d", resp.StatusCode())
	}

	err = resizeTmpFile(tmpImage, previewImage)
	if err != nil {
		return "", fmt.Errorf("sync: error resize image, %s", err)
	}

	return previewImage, nil
}

func clearOldPreviews() {
	//	todo
}

func (s SyncWorker) Stop() {

}
func (s SyncWorker) EnqueueUniquePeriodicWork() {
	ss := gocron.NewScheduler(time.UTC)
	ss.Every(int(10)).Minutes().Do(updateMenu, &s)
	ss.StartAsync()
}

var updateMenu = func(s *SyncWorker) {
	err := s.Start()
	if err != nil {
		log.Printf("sync: err, %s", err)
	}
}

func New(r *repositoryProduct.Repository, client *c.RestoClient, resty *resty.Client) *SyncWorker {
	return &SyncWorker{
		r:         r,
		evoClient: client,
		client:    resty,
	}
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

	return jpeg.Encode(out, m, nil)
}

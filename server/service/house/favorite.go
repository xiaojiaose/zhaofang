package house

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
)

type FavoriteServer struct {
}

func (service *FavoriteServer) CreateOrUpdate(favorite *house.Favorite) (err error) {
	err = global.GVA_DB.Where("id = ?", favorite.ID).First(&house.Favorite{}).Updates(&favorite).Error
	if err != nil && err.Error() == "record not found" {
		err = global.GVA_DB.Create(favorite).Error
	}

	return
}

func (service *FavoriteServer) Delete(ID uint) (err error) {
	err = global.GVA_DB.Delete(&house.Favorite{}, "id = ?", ID).Error
	return
}

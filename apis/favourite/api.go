package favourite

import (
	"github.com/opentreehole/go-common"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	. "treehole_next/models"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
)

// ListFavorites
//
// @Summary List User's Favorites
// @Tags Favorite
// @Produce application/json
// @Router /user/favorites [get]
// @Param object query ListFavoriteModel false "query"
// @Success 200 {object} models.Map
// @Success 200 {array} models.Hole
func ListFavorites(c *fiber.Ctx) error {
	// get userID
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	var query ListFavoriteModel
	err = common.ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	if query.Plain {
		// get favorite ids
		data, err := UserGetFavoriteDataByFavoriteGroup(DB, userID, query.FavoriteGroupID)
		if err != nil {
			return err
		}
		return c.JSON(Map{"data": data})
	} else {
		// get order
		var order string
		switch query.Order {
		case "id":
			order = "hole.id desc"
		case "time_created":
			order = "user_favorites.created_at desc, hole.id desc"
		case "hole_time_updated":
			order = "hole.updated_at desc"
		}

		// get favorites
		holes := make(Holes, 0)
		err = DB.
			Joins("JOIN user_favorites ON user_favorites.hole_id = hole.id AND user_favorites.user_id = ? AND user_favorites.favorite_group_id = ?", userID, query.FavoriteGroupID).
			Order(order).Find(&holes).Error
		if err != nil {
			return err
		}
		return Serialize(c, &holes)
	}
}

// AddFavorite
//
// @Summary Add A Favorite
// @Tags Favorite
// @Accept application/json
// @Produce application/json
// @Router /user/favorites [post]
// @Param json body AddModel true "json"
// @Success 201 {object} Response
// @Success 200 {object} Response
func AddFavorite(c *fiber.Ctx) error {
	// validate body
	var body AddModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// get userID
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	var data []int

	err = DB.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		// add favorite
		err = AddUserFavorite(tx, userID, body.HoleID, body.FavoriteGroupID)
		if err != nil {
			return err
		}

		// create response
		data, err = UserGetFavoriteData(tx, userID)
		return err
	})
	if err != nil {
		return err
	}

	return c.Status(201).JSON(&Response{
		Message: "收藏成功",
		Data:    data,
	})
}

// ModifyFavorite
//
// @Summary Modify User's Favorites
// @Tags Favorite
// @Produce application/json
// @Router /user/favorites [put]
// @Param json body ModifyModel true "json"
// @Success 200 {object} Response
// @Failure 404 {object} Response
func ModifyFavorite(c *fiber.Ctx) error {
	// validate body
	var body ModifyModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// get userID
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	// modify favorite
	err = ModifyUserFavorite(DB, userID, body.HoleIDs, body.FavoriteGroupID)
	if err != nil {
		return err
	}

	// create response
	data, err := UserGetFavoriteData(DB, userID)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(&Response{
		Message: "修改成功",
		Data:    data,
	})
}

// DeleteFavorite
//
// @Summary Delete A Favorite
// @Tags Favorite
// @Produce application/json
// @Router /user/favorites [delete]
// @Param json body DeleteModel true "json"
// @Success 200 {object} Response
// @Failure 404 {object} Response
func DeleteFavorite(c *fiber.Ctx) error {
	// validate body
	var body DeleteModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// get userID
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	// delete favorite
	err = DeleteUserFavorite(DB, userID, body.HoleID, body.FavoriteGroupID)
	if err != nil {
		return err
	}

	// create response
	data, err := UserGetFavoriteData(DB, userID)
	if err != nil {
		return err
	}

	return c.JSON(&Response{
		Message: "删除成功",
		Data:    data,
	})
}

// ListFavoriteGroups
//
// @Summary List User's Favorite Groups
// @Tags Favorite
// @Produce application/json
// @Router /user/favorite_group [get]
// @Param object query ListFavoriteGroupModel false "query"
// @Success 200 {object} models.Map
// @Success 200 {array} models.FavoriteGroups
func ListFavoriteGroups(c *fiber.Ctx) error {
	// get userID
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	var query ListFavoriteGroupModel
	err = common.ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	if query.Plain {
		// get favoriteGroups
		data, err := UserGetFavoriteGroups(DB, userID)
		if err != nil {
			return err
		}
		return c.JSON(Map{"data": data})
	} else {
		// get order
		var order string
		switch query.Order {
		case "id":
			order = "id desc"
		case "time_created":
			order = "created_at desc, id desc"
		case "time_updated":
			order = "updated_at desc, id desc"
		}

		// get favoriteGroups
		err = DB.Where("user_id = ? AND deleted = false", userID).Order(order).Find(&FavoriteGroups{}).Error
		if err != nil {
			return err
		}
		return c.JSON(Map{"data": FavoriteGroups{}})
	}
}

// AddFavoriteGroup
//
// @Summary Add A Favorite Group
// @Tags Favorite
// @Accept application/json
// @Produce application/json
// @Router /user/favorite_group [post]
// @Param json body AddFavoriteGroupModel true "json"
// @Success 201 {object} Response
// @Success 200 {object} Response
func AddFavoriteGroup(c *fiber.Ctx) error {
	// validate body
	var body AddFavoriteGroupModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// get userID
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	// add favorite group
	err = AddUserFavoriteGroup(DB, userID, body.Name)
	if err != nil {
		return err
	}

	// create response
	data, err := UserGetFavoriteData(DB, userID)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(&Response{
		Message: "添加成功",
		Data:    data,
	})
}

// ModifyFavoriteGroup
//
// @Summary Modify User's Favorite Group
// @Tags Favorite
// @Produce application/json
// @Router /user/favorite_group [put]
// @Param json body ModifyFavoriteGroupModel true "json"
// @Success 200 {object} Response
// @Failure 404 {object} Response
func ModifyFavoriteGroup(c *fiber.Ctx) error {
	// validate body
	var body ModifyFavoriteGroupModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// get userID
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	// modify favorite group
	err = ModifyUserFavoriteGroup(DB, userID, body.FavoriteGroupID, body.Name)
	if err != nil {
		return err
	}

	// create response
	data, err := UserGetFavoriteData(DB, userID)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(&Response{
		Message: "修改成功",
		Data:    data,
	})
}

// DeleteFavoriteGroup
//
// @Summary Delete A Favorite Group
// @Tags Favorite
// @Produce application/json
// @Router /user/favorite_group [delete]
// @Param json body DeleteModel true "json"
// @Success 200 {object} Response
// @Failure 404 {object} Response
func DeleteFavoriteGroup(c *fiber.Ctx) error {
	// validate body
	var body DeleteModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// get userID
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	// delete favorite group
	err = DeleteUserFavoriteGroup(DB, userID, body.FavoriteGroupID)
	if err != nil {
		return err
	}

	//create response
	data, err := UserGetFavoriteData(DB, userID)
	if err != nil {
		return err
	}

	return c.JSON(&Response{
		Message: "删除成功",
		Data:    data,
	})
}

// MoveFavorite
//
// @Summary Move User's Favorite
// @Tags Favorite
// @Produce application/json
// @Router /user/favorite_group/move [put]
// @Param json body MoveModel true "json"
// @Success 200 {object} Response
// @Failure 404 {object} Response
func MoveFavorite(c *fiber.Ctx) error {
	// validate body
	var body MoveModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// get userID
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	// move favorite
	err = MoveUserFavorite(DB, userID, body.HoleIDs, body.FromFavoriteGroupID, body.ToFavoriteGroupID)
	if err != nil {
		return err
	}

	// create response
	data, err := UserGetFavoriteData(DB, userID)
	if err != nil {
		return err
	}

	return c.JSON(&Response{
		Message: "移动成功",
		Data:    data,
	})
}

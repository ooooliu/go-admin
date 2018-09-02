package menu

import (
	"github.com/chenhg5/go-admin/modules/connections"
	"strconv"
	"sync"
)

type Menu struct {
	GlobalMenuList   []MenuItem
	GlobalMenuOption []map[string]string
	MaxOrder         int64
}

var GlobalMenu = &Menu{
	GlobalMenuList:   []MenuItem{},
	GlobalMenuOption: []map[string]string{},
	MaxOrder:         0,
}

var InitMenuOnce sync.Once

func InitMenu() {
	InitMenuOnce.Do(func() {
		SetGlobalMenu()
	})
}

func SetGlobalMenu() {
	menus, _ := connections.GetConnection().Query("select * from goadmin_menu where id > 0 order by `order` asc")

	menuOption := make([]map[string]string, 0)

	for i := 0; i < len(menus); i++ {
		menuOption = append(menuOption, map[string]string{
			"id":    strconv.FormatInt(menus[i]["id"].(int64), 10),
			"title": menus[i]["title"].(string),
		})
	}

	menulist := ConstructMenuTree(menus, 0)

	(*GlobalMenu).GlobalMenuOption = menuOption
	(*GlobalMenu).GlobalMenuList = menulist
	(*GlobalMenu).MaxOrder = menus[len(menus)-1]["parent_id"].(int64)
}

func (menu *Menu) SexMaxOrder(order int64) {
	menu.MaxOrder = order
}

func ConstructMenuTree(menus []map[string]interface{}, parentId int64) []MenuItem {

	branch := make([]MenuItem, 0)

	for j := 0; j < len(menus); j++ {
		if parentId == menus[j]["parent_id"].(int64) {

			childList := ConstructMenuTree(menus, menus[j]["id"].(int64))

			child := MenuItem{
				Name:         menus[j]["title"].(string),
				ID:           strconv.FormatInt(menus[j]["id"].(int64), 10),
				Url:          menus[j]["uri"].(string),
				Icon:         menus[j]["icon"].(string),
				Active:       "",
				ChildrenList: childList,
			}

			branch = append(branch, child)
		}
	}

	return branch
}

func GetMenuItemById(id string) MenuItem {
	menu, _ := connections.GetConnection().Query("select * from goadmin_menu where id = ?", id)

	return MenuItem{
		Name:         menu[0]["title"].(string),
		ID:           strconv.FormatInt(menu[0]["id"].(int64), 10),
		Url:          menu[0]["uri"].(string),
		Icon:         menu[0]["icon"].(string),
		Active:       "",
		ChildrenList: []MenuItem{},
	}
}

func (menu *Menu) SetActiveClass(path string) {

	for i := 0; i < len((*menu).GlobalMenuList); i++ {
		(*menu).GlobalMenuList[i].Active = ""
	}

Loop:
	for i := 0; i < len((*menu).GlobalMenuList); i++ {
		if (*menu).GlobalMenuList[i].Url == path {
			(*menu).GlobalMenuList[i].Active = "active"
			break Loop
		} else {
			for j := 0; j < len((*menu).GlobalMenuList[i].ChildrenList); j++ {
				if (*menu).GlobalMenuList[i].ChildrenList[j].Url == path {
					(*menu).GlobalMenuList[i].Active = "active"
					break Loop
				} else {
					(*menu).GlobalMenuList[i].Active = ""
				}
			}
		}
	}
}

type MenuItem struct {
	Name         string
	ID           string
	Url          string
	Icon         string
	Active       string
	ChildrenList []MenuItem
}

func (menu *Menu) GetEditMenuList() []MenuItem {
	return (*menu).GlobalMenuList
}
package mp

import "github.com/wizjin/weixin"

func CreateMenu()(err error) {
	wx := weixin.New("sadfasdf21312312asdfasdf", "wxf976c7be9eb785ba", "189d6d22df9877bd695901cde9f1994d")
	menu:=&weixin.Menu{Buttons: make([]weixin.MenuButton, 2)}
	err = wx.CreateMenu(menu)
	if err != nil {
		return
	}
	return
}

func DeleteMenu()(err error) {
	wx := weixin.New("sadfasdf21312312asdfasdf", "wxf976c7be9eb785ba", "189d6d22df9877bd695901cde9f1994d")
	err = wx.DeleteMenu()
	if err != nil {
		return
	}
	return
}

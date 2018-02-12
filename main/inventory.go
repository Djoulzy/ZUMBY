package main

func (W *world) inventoryAdd(infos inventory) {
	itemList, _ := W.UserList.Get(infos.Owner)
	user := itemList.(*user)
	W.Map.Items[infos.X][infos.Y].Owner = infos.Owner
	user.Inventory[infos.ToPocket] = W.Map.Items[infos.X][infos.Y]
}

func (W *world) inventoryDrop(infos inventory) {
	W.Map.Items[infos.X][infos.Y] = item{
		ID: infos.ID,
	}
}

func (W *world) inventoryUpdate(infos inventory) {
	itemList, _ := W.UserList.Get(infos.Owner)
	user := itemList.(*user)
	user.Inventory[infos.FromPocket] = user.Inventory[infos.ToPocket]
	user.Inventory[infos.ToPocket] = item{
		ID: infos.ID,
	}
}

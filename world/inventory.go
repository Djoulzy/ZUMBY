package world

func (W *WORLD) inventoryAdd(infos INVENTORY) {
	item, _ := W.UserList.Get(infos.Owner)
	user := item.(*USER)
	W.Map.Items[infos.X][infos.Y].Owner = infos.Owner
	user.Inventory[infos.ToPocket] = W.Map.Items[infos.X][infos.Y]
}

func (W *WORLD) inventoryDrop(infos INVENTORY) {
	W.Map.Items[infos.X][infos.Y] = ITEM{
		ID: infos.ID,
	}
}

func (W *WORLD) inventoryUpdate(infos INVENTORY) {
	item, _ := W.UserList.Get(infos.Owner)
	user := item.(*USER)
	user.Inventory[infos.FromPocket] = user.Inventory[infos.ToPocket]
	user.Inventory[infos.ToPocket] = ITEM{
		ID: infos.ID,
	}
}

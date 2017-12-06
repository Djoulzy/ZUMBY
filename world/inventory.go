package world

func (W *WORLD) inventoryAdd(infos INVENTORY) {
	user := W.UserList[infos.Owner]
	W.Map.Items[infos.X][infos.Y].Owner = infos.Owner
	user.Inventory[infos.ToPocket] = W.Map.Items[infos.X][infos.Y]
}

func (W *WORLD) inventoryDrop(infos INVENTORY) {
	W.Map.Items[infos.X][infos.Y] = ITEM{
		ID: infos.ID,
	}
}

func (W *WORLD) inventoryUpdate(infos INVENTORY) {
	user := W.UserList[infos.Owner]
	user.Inventory[infos.FromPocket] = user.Inventory[infos.ToPocket]
	user.Inventory[infos.ToPocket] = ITEM{
		ID: infos.ID,
	}
}

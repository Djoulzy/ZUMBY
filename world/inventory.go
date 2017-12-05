package world

func (W *WORLD) inventoryAdd(infos INVENTORY) {
	user := W.UserList[infos.Owner]
	W.Map.Items[infos.X][infos.Y].Owner = infos.Owner
	user.Inventory[infos.ToPocket] = W.Map.Items[infos.X][infos.Y]
}

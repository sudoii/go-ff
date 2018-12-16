package incoming

import (
	"flyff/common/service/cache"
	"flyff/common/service/database"
	"flyff/common/service/external"
	"flyff/common/service/messaging"
	"flyff/entity/packets/outgoing"
	"log"
)

// Disconnect from World
func Disconnect(p *external.PacketHandler) {
	player := cache.FindByNetID(p.ClientID)
	if player == nil {
		return
	}

	cache.Connection.Where("net_client_id = ?", p.ClientID).Delete(player)

	var dbPlayer database.Player
	err := database.Connection.Model(&dbPlayer).Preload("Items").Where("id = ?", player.PlayerID).First(&dbPlayer).Error
	if err != nil {
		log.Print(err)
		return
	}

	// Clear previous inventory (Pretty bad yeah...)
	database.Connection.Model(&database.Item{}).Where("player_id = ?", player.PlayerID).Delete(database.Item{})

	dbPlayer.FaceID = player.FaceID
	dbPlayer.Gender = player.Gender
	dbPlayer.HairColor = player.HairColor
	dbPlayer.Items = player.Inventory.ConvertToDatabaseSlice()
	dbPlayer.JobID = player.JobID
	dbPlayer.Level = player.Level
	dbPlayer.Name = player.Name
	dbPlayer.Position = player.Position
	dbPlayer.SkinSetID = player.SkinSetID
	dbPlayer.Statistics = player.Statistics

	database.Connection.Save(dbPlayer)

	// Make Visible ----
	messaging.Publish(messaging.ConnectionTopic, &external.PacketEmitter{
		Packet: outgoing.DeleteObj(player),
		To:     cache.FindIDAround(player),
	})
}

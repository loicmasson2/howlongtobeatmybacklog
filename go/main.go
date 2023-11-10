package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/games", func(c *gin.Context) {
		const STEAM_ID string = "76561198000800114"
		const STEAM_KEY string = "8FEF865E63A65A65E8C79C69CCDC1034"
		url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=true", STEAM_KEY, STEAM_ID)

		resp, err := http.Get(url)
		if err != nil {
			// handle error
			log.Fatal(err)
		}
		defer resp.Body.Close()
		var generic map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&generic)
		if err != nil {
			log.Fatalf("Parse response failed, reason: %v \n", err)
		}
		//fmt.Printf("%s", generic)
		c.JSON(http.StatusOK, gin.H{"data": generic})

		// format data
		//{
		//	"appid": 1900,
		//	"img_icon_url": "0a721907de90582dd4a53b55a2f260df19e2c72b",
		//	"name": "Earth 2160",
		//	"playtime_disconnected": 0,
		//	"playtime_forever": 2,
		//	"playtime_linux_forever": 0,
		//	"playtime_mac_forever": 0,
		//	"playtime_windows_forever": 0,
		//	"rtime_last_played": 1354780800
		//},
	})
	r.GET("/test-db", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@mongo:27017"))
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()

		collection := client.Database("howlongtobeatmybacklog").Collection("users")
		filter := bson.D{{"steam_id", "76561198000800114"}}
		var generic map[string]interface{}
		err = collection.FindOne(context.TODO(), filter).Decode(&generic)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(collection)
		//fmt.Println(collection)

		//err = client.Disconnect(context.TODO())
		//
		//if err != nil {
		//	log.Fatal(err)
		//}
		//fmt.Println("Connection to MongoDB closed.")

		fmt.Printf("%s", generic)
		c.JSON(http.StatusOK, gin.H{"data": generic})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	//	http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=8FEF865E63A65A65E8C79C69CCDC1034&steamid=76561198000800114&format=json&include_appinfo=true
}

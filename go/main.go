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

type SteamResponse struct {
	Response struct {
		GameCount int `json:"game_count"`
		Games     []struct {
			Appid                    float64   `json:"appid"`
			ImgIconURL               string    `json:"img_icon_url"`
			Name                     string    `json:"name"`
			PlaytimeDisconnected     float64   `json:"playtime_disconnected"`
			PlaytimeForever          float64   `json:"playtime_forever"`
			PlaytimeLinuxForever     float64   `json:"playtime_linux_forever"`
			PlaytimeMacForever       float64   `json:"playtime_mac_forever"`
			PlaytimeWindowsForever   float64   `json:"playtime_windows_forever"`
			RtimeLastPlayed          float64   `json:"rtime_last_played"`
			HasCommunityVisibleStats bool      `json:"has_community_visible_stats,omitempty"`
			ContentDescriptorids     []float64 `json:"content_descriptorids,omitempty"`
			HasLeaderboards          bool      `json:"has_leaderboards,omitempty"`
			Playtime2Weeks           float64   `json:"playtime_2weeks,omitempty"`
		} `json:"games"`
	} `json:"response"`
}

type Game struct {
	Appid                  float64 `json:"appid,omitempty" bson:"appid,omitempty"`
	ImgIconUrl             string  `json:"img_icon_url" bson:"img_icon_url"`
	Name                   string  `json:"name" bson:"name"`
	PlaytimeDisconnected   float64 `json:"playtime_disconnected" bson:"playtime_disconnected"`
	PlaytimeForever        float64 `json:"playtime_forever" bson:"playtime_forever"`
	PlaytimeLinuxForever   float64 `json:"playtime_linux_forever" bson:"playtime_linux_forever"`
	PlaytimeMacForever     float64 `json:"playtime_mac_forever" bson:"playtime_mac_forever"`
	PlaytimeWindowsForever float64 `json:"playtime_windows_forever" bson:"playtime_windows_forever"`
	RtimeLastPlayed        float64 `json:"rtime_last_played" bson:"rtime_last_played"`
}

type Document struct {
	SteamId string `json:"steam_id" bson:"steam_id"`
	Games   []struct {
		Appid                    float64   `json:"appid"`
		ImgIconURL               string    `json:"img_icon_url"`
		Name                     string    `json:"name"`
		PlaytimeDisconnected     float64   `json:"playtime_disconnected"`
		PlaytimeForever          float64   `json:"playtime_forever"`
		PlaytimeLinuxForever     float64   `json:"playtime_linux_forever"`
		PlaytimeMacForever       float64   `json:"playtime_mac_forever"`
		PlaytimeWindowsForever   float64   `json:"playtime_windows_forever"`
		RtimeLastPlayed          float64   `json:"rtime_last_played"`
		HasCommunityVisibleStats bool      `json:"has_community_visible_stats,omitempty"`
		ContentDescriptorids     []float64 `json:"content_descriptorids,omitempty"`
		HasLeaderboards          bool      `json:"has_leaderboards,omitempty"`
		Playtime2Weeks           float64   `json:"playtime_2weeks,omitempty"`
	} `json:"games"`
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/games", func(c *gin.Context) {
		const SteamId string = "76561198000800114"
		const SteamKey string = "8FEF865E63A65A65E8C79C69CCDC1034"
		url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json", SteamKey, SteamId)

		resp, err := http.Get(url)
		if err != nil {
			// handle error
			log.Fatal(err)
		}
		defer resp.Body.Close()
		fmt.Printf("%+v\n", resp.Body)
		var generic map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&generic)
		if err != nil {
			log.Fatalf("Parse response failed, reason: %v \n", err)
		}
		//fmt.Printf("%s", generic)
		c.JSON(http.StatusOK, gin.H{"data": generic})

	})
	r.GET("/insert-games/:steamid", func(c *gin.Context) {
		steamid := c.Param("steamid")
		const SteamKey string = "8FEF865E63A65A65E8C79C69CCDC1034"
		url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=true", SteamKey, steamid)
		r, err := http.Get(url)
		if err != nil {
			log.Println("Cannot get from URL", err)
		}
		defer r.Body.Close()

		var test SteamResponse
		err = json.NewDecoder(r.Body).Decode(&test)
		if err != nil {
			log.Fatal(err)
			log.Println("Error unmarshalling json data:", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@mongo:27017"))
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()

		coll := client.Database("howlongtobeatmybacklog").Collection("documents")
		newDocument := []interface{}{
			Document{SteamId: steamid, Games: test.Response.Games},
		}
		result, err := coll.InsertMany(context.TODO(), newDocument)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%d documents inserted with IDs:\n", len(result.InsertedIDs))
		for _, id := range result.InsertedIDs {
			fmt.Printf("\t%s\n", id)
		}
		c.JSON(http.StatusOK, gin.H{"data": len(result.InsertedIDs)})

	})
	r.GET("/get-games-from-db/:steamid", func(c *gin.Context) {
		steamid := c.Param("steamid")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@mongo:27017"))
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()

		collection := client.Database("howlongtobeatmybacklog").Collection("documents")
		filter := bson.D{{"steam_id", steamid}}
		var document Document
		err = collection.FindOne(context.TODO(), filter).Decode(&document)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s", collection)
		c.JSON(http.StatusOK, gin.H{"data": document})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	//	http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=8FEF865E63A65A65E8C79C69CCDC1034&steamid=76561198000800114&format=json&include_appinfo=true
}

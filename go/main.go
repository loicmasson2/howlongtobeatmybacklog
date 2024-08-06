package main

import (
	"encoding/json"
	"fmt"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"howlongtobeatmybacklog/components"
	"log"
	"net/http"
	"strings"
)

type Game struct {
	Appid                    float64
	ImgIconURL               string
	Name                     string
	PlaytimeDisconnected     float64
	PlaytimeForever          float64
	PlaytimeLinuxForever     float64
	PlaytimeMacForever       float64
	PlaytimeWindowsForever   float64
	RtimeLastPlayed          float64
	HasCommunityVisibleStats bool
	ContentDescriptorids     []float64
	HasLeaderboards          bool
	Playtime2Weeks           float64
}

type SteamResponse struct {
	Response struct {
		GameCount int
		Games     []components.Game
	}
}

//type Document struct {
//	SteamId string `json:"steam_id" bson:"steam_id"`
//	Games   []struct {
//		Appid                    float64   `json:"appid"`
//		ImgIconURL               string    `json:"img_icon_url"`
//		Name                     string    `json:"name"`
//		PlaytimeDisconnected     float64   `json:"playtime_disconnected"`
//		PlaytimeForever          float64   `json:"playtime_forever"`
//		PlaytimeLinuxForever     float64   `json:"playtime_linux_forever"`
//		PlaytimeMacForever       float64   `json:"playtime_mac_forever"`
//		PlaytimeWindowsForever   float64   `json:"playtime_windows_forever"`
//		RtimeLastPlayed          float64   `json:"rtime_last_played"`
//		HasCommunityVisibleStats bool      `json:"has_community_visible_stats,omitempty"`
//		ContentDescriptorids     []float64 `json:"content_descriptorids,omitempty"`
//		HasLeaderboards          bool      `json:"has_leaderboards,omitempty"`
//		Playtime2Weeks           float64   `json:"playtime_2weeks,omitempty"`
//	} `json:"games"`
//}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

type SearchRequest struct {
	Search string
}

func main() {

	// TODO create a layout
	//r.Use(CORSMiddleware())
	//r.HTMLRender = &TemplRender{}
	//r.GET("/home", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "", hello("LOL"))
	//})

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	//serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	//opts := options.Client().ApplyURI("mongodb+srv://massonloic023:2zxjCi39j0dPo2jA@howlongtobeatmybacklogd.gbv6bax.mongodb.net/?retryWrites=true&w=majority&appName=Howlongtobeatmybacklogdb").SetServerAPIOptions(serverAPI)
	//
	//// Create a new client and connect to the server
	//client, err := mongo.Connect(context.TODO(), opts)
	//if err != nil {
	//	panic(err)
	//}
	//
	//defer func() {
	//	if err = client.Disconnect(context.TODO()); err != nil {
	//		panic(err)
	//	}
	//}()
	//
	//// Send a ping to confirm a successful connection
	//if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
	//	panic(err)
	//}
	//fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	r := chi.NewRouter()
	p := components.Paragraph("Dynamic contents")
	r.Get("/home", templ.Handler(components.Root(p)).ServeHTTP)
	r.Get("/games", func(w http.ResponseWriter, r *http.Request) {
		const SteamId string = "76561198000800114"
		const SteamKey string = "8FEF865E63A65A65E8C79C69CCDC1034"
		url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=true", SteamKey, SteamId)
		fmt.Println(url)
		resp, err := http.Get(url)
		if err != nil {
			// handle error
			log.Fatal(err)
		}
		defer resp.Body.Close()
		fmt.Printf("%+v\n", resp.Body)
		var generic SteamResponse
		err = json.NewDecoder(resp.Body).Decode(&generic)
		if err != nil {
			log.Fatalf("Parse response failed, reason: %v \n", err)
		}
		list := components.NameList(generic.Response.Games)
		templ.Handler(components.Root(list)).ServeHTTP(w, r)
		//w.Write([]byte("hi"))
	})
	r.Post("/search", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			// if there was an error here, we return the error
			// as response along with a 400 http response code
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// once the form has been parsed correctly
		// we can obtain the form fields  using
		// req.FormValue passing the field name as the key
		// see ./post-method.html for the field names

		// it's worth noting here that FormValue() will return
		// the first value if there are multiple values specified
		// in the request
		// if you want to access the multiple values specified,
		// use req.Form directly (see https://pkg.go.dev/net/http#Request)
		search := r.FormValue("search")
		fmt.Printf("Search is %s\n", search)
		const SteamId string = "76561198000800114"
		const SteamKey string = "8FEF865E63A65A65E8C79C69CCDC1034"
		url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=true", SteamKey, SteamId)
		fmt.Println(url)
		resp, err := http.Get(url)
		if err != nil {
			// handle error
			log.Fatal(err)
		}
		defer resp.Body.Close()
		fmt.Printf("%+v\n", resp.Body)
		var generic SteamResponse
		err = json.NewDecoder(resp.Body).Decode(&generic)
		if err != nil {
			log.Fatalf("Parse response failed, reason: %v \n", err)
		}
		positive := []components.Game{}

		for i := range generic.Response.Games {
			if strings.Contains(strings.ToLower(generic.Response.Games[i].Name), strings.ToLower(search)) {
				fmt.Printf("Name is %s\n", generic.Response.Games[i].Name)

				positive = append(positive, generic.Response.Games[i])
			}
		}
		//list := components.NameList(positive)

		templ.Handler(components.NameList(positive)).ServeHTTP(w, r)
		//w.Write([]byte("hi"))
	})
	//r.GET("/games", func(c *gin.Context) {
	//	const SteamId string = "76561198000800114"
	//	const SteamKey string = "8FEF865E63A65A65E8C79C69CCDC1034"
	//	url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=true", SteamKey, SteamId)
	//
	//	resp, err := http.Get(url)
	//	if err != nil {
	//		// handle error
	//		log.Fatal(err)
	//	}
	//	defer resp.Body.Close()
	//	fmt.Printf("%+v\n", resp.Body)
	//	var generic SteamResponse
	//	err = json.NewDecoder(resp.Body).Decode(&generic)
	//	if err != nil {
	//		log.Fatalf("Parse response failed, reason: %v \n", err)
	//	}
	//	fmt.Println(generic.Response.GameCount)
	//	c.HTML(http.StatusOK, "", nameList(generic.Response.Games))
	//
	//})
	//r.POST("/games/:steamid", func(c *gin.Context) {
	//	steamid := c.Param("steamid")
	//	const SteamKey string = "8FEF865E63A65A65E8C79C69CCDC1034"
	//	url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=true", SteamKey, steamid)
	//	r, err := http.Get(url)
	//	if err != nil {
	//		log.Println("Cannot get from URL", err)
	//	}
	//	defer r.Body.Close()
	//
	//	var test SteamResponse
	//	err = json.NewDecoder(r.Body).Decode(&test)
	//	if err != nil {
	//		log.Fatal(err)
	//		log.Println("Error unmarshalling json data:", err)
	//	}
	//
	//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//	defer cancel()
	//	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@mongo:27017"))
	//	defer func() {
	//		if err = client.Disconnect(ctx); err != nil {
	//			panic(err)
	//		}
	//	}()
	//
	//	coll := client.Database("howlongtobeatmybacklog").Collection("documents")
	//	newDocument := []interface{}{
	//		Document{SteamId: steamid, Games: test.Response.Games},
	//	}
	//	result, err := coll.InsertMany(context.TODO(), newDocument)
	//	if err != nil {
	//		panic(err)
	//	}
	//	model := mongo.IndexModel{Keys: bson.D{{"games.name", "text"}}}
	//	name, err := coll.Indexes().CreateOne(context.TODO(), model)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println("Name of index created: " + name)
	//
	//	fmt.Printf("%d documents inserted with IDs:\n", len(result.InsertedIDs))
	//	for _, id := range result.InsertedIDs {
	//		fmt.Printf("\t%s\n", id)
	//	}
	//	c.JSON(http.StatusOK, gin.H{"data": len(result.InsertedIDs)})
	//
	//})
	//r.GET("/games/:steamid", func(c *gin.Context) {
	//	steamid := c.Param("steamid")
	//	search := c.DefaultQuery("search", "")
	//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//	defer cancel()
	//	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@mongo:27017"))
	//	defer func() {
	//		if err = client.Disconnect(ctx); err != nil {
	//			panic(err)
	//		}
	//	}()
	//
	//	collection := client.Database("howlongtobeatmybacklog").Collection("documents")
	//	matchStage := bson.D{{"$match", bson.D{{"steam_id", steamid}}}}
	//	unwindStage := bson.D{
	//		{
	//			"$unwind", bson.D{
	//				{
	//					"path", "$games",
	//				},
	//			},
	//		},
	//	}
	//	filterStage := bson.D{{"$match", bson.D{{"games.name", search}}}}
	//
	//	cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, unwindStage, filterStage})
	//
	//	var results []map[string]interface{}
	//	if err = cursor.All(context.TODO(), &results); err != nil {
	//		fmt.Println("ERROR")
	//		log.Fatal(err)
	//	}
	//	fmt.Println(results)
	//
	//	c.JSON(http.StatusOK, gin.H{"data": results})
	//})
	http.ListenAndServe(":4000", r)

	//r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

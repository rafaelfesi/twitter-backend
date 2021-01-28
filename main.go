package main

import (
	"github.com/arman-aminian/twitter-backend/db"
	"os"

	//_ "github.com/arman-aminian/twitter-backend/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/arman-aminian/twitter-backend/handler"
	"github.com/arman-aminian/twitter-backend/model"
	"github.com/arman-aminian/twitter-backend/router"
	"github.com/arman-aminian/twitter-backend/store"
	//echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	testHeap := false

	if !testHeap {
		r := router.New()
		//r.GET("/swagger/*", echoSwagger.WrapHandler)
		mongoClient, err := db.GetMongoClient()
		if err != nil {
			log.Fatal(err)
		}
		usersDb := db.SetupUsersDb(mongoClient)
		tweetsDb := db.SetupTweetsDb(mongoClient)
		hashtagsDb := db.SetupHashtagsDb(mongoClient)
		g := r.Group("")
		us := store.NewUserStore(usersDb)
		ts := store.NewTweetStore(tweetsDb)
		hs := store.NewHashtagStore(hashtagsDb)
		h := handler.NewHandler(us, ts, hs)
		h.Register(g)

		// Fire up the trends beforehand
		_ = hs.Update()

		// RUN THIS IF YOUR HASHTAG DATABASE IS EMPTY
		// StartUpTrends(ts, h)
		
		r.Logger.Fatal(r.Start(":" + port))
	}

	// if testHeap {
	// 	// Test HashtagStore
	// 	hs := store.NewHashtagStore()
	// 	heap.Init(hs)
	// 	heap.Push(hs, &model.Hashtag{Name: "salam", Tweets: nil, Count: 10})
	// 	heap.Push(hs, &model.Hashtag{Name: "salam2", Tweets: nil, Count: 9})
	// 	heap.Push(hs, &model.Hashtag{Name: "salam3", Tweets: nil, Count: 8})
	// 	heap.Push(hs, &model.Hashtag{Name: "salam4", Tweets: nil, Count: 7})
	// 	heap.Push(hs, &model.Hashtag{Name: "salam5", Tweets: nil, Count: 6})
	// 	heap.Push(hs, &model.Hashtag{Name: "salam6", Tweets: nil, Count: 5})
	// 	heap.Push(hs, &model.Hashtag{Name: "salam7", Tweets: nil, Count: 4})
	// 	heap.Push(hs, &model.Hashtag{Name: "salam8", Tweets: nil, Count: 882})
	// 	heap.Push(hs, &model.Hashtag{Name: "salam9", Tweets: nil, Count: 100})
	// 	for hs.Len() > 0 {
	// 		fmt.Printf("%v ", heap.Pop(hs))
	// 	}
	// 	fmt.Println()
	// 	// Changing one hashtag's size
	// 	heap.Push(hs, &model.Hashtag{Name: "salam", Tweets: nil, Count: 10})
	// 	heap.Push(hs, &model.Hashtag{Name: "salam2", Tweets: nil, Count: 9})
	// 	heap.Push(hs, &model.Hashtag{Name: "salam3", Tweets: nil, Count: 8})
	// 	heap.Push(hs, &model.Hashtag{Name: "salam4", Tweets: nil, Count: 7})
	// 	heap.Push(hs, &model.Hashtag{Name: "salam5", Tweets: nil, Count: 6})
	// 	heap.Push(hs, &model.Hashtag{Name: "salam6", Tweets: nil, Count: 5})
	// 	heap.Push(hs, &model.Hashtag{Name: "salam7", Tweets: nil, Count: 4})
	// 	heap.Push(hs, &model.Hashtag{Name: "salam8", Tweets: nil, Count: 882})
	// 	heap.Push(hs, &model.Hashtag{Name: "salam9", Tweets: nil, Count: 100})
	// 	var temp []*model.Hashtag
	// 	oldVal := 8
	// 	newVal := 125
	// 	for hs.Len() > 0 {
	// 		h := heap.Pop(hs).(*model.Hashtag)
	// 		prevCount := h.Count
	// 		if prevCount == oldVal {
	// 			tempH := model.Hashtag{
	// 				Name:   h.Name,
	// 				Tweets: h.Tweets,
	// 				Count:  newVal,
	// 			}
	// 			temp = append(temp, &tempH)
	// 			break
	// 		} else {
	// 			temp = append(temp, h)
	// 		}
	// 	}
	//
	// 	for _, h := range temp {
	// 		heap.Push(hs, h)
	// 	}
	// 	for hs.Len() > 0 {
	// 		fmt.Printf("%v ", heap.Pop(hs))
	// 	}
	// 	fmt.Println()
	// }
}

func StartUpTrends(ts *store.TweetStore, h *handler.Handler) {
	allTweets, err := ts.GetAllTweets()
	if err != nil {
		log.Fatal(err)
	}
	for _, bm := range allTweets {
		var t *model.Tweet
		bsonBytes, _ := bson.Marshal(bm)
		_ = bson.Unmarshal(bsonBytes, &t)
		hashtags := ts.ExtractHashtags(t)
		for name, cnt := range hashtags {
			h.AddHashtag(name, t, cnt)
		}
	}
}

func populateAll(us *store.UserStore, ts *store.TweetStore, hh *store.HashtagStore) {
	for _, s := range []string{"user1", "user2", "user3", "user4"} {
		u := populateUser(s, s+"@gmail.com", s+"_pass", s+" bio", s+" profile", s+" header")
		_ = us.Create(u)
	}

	// ts.CreateTweet()
}

func populateUser(username, email, password, bio, profile, header string) *model.User {
	u := model.NewUser()
	u.Username = username
	u.Email = email
	u.Password = password
	u.Bio = bio
	u.ProfilePicture = profile
	u.HeaderPicture = header
	return u
}

package main

import (
  "net/http"
  "fmt"
  "io/ioutil"
  "time"
  "encoding/json"
  "strconv"
  "./helpers"
  "./errors"
  "log"
  "os"
  "github.com/gorilla/handlers"
)

// Number of posts to cache server side
var limit = 100


var url = "http://www.reddit.com/r/listentothis/.json?limit=" + strconv.Itoa(limit)
var uAgent = "rjpj's obscure agent v0.0.1"

// The proverbial expiration date
var refreshTime = int64(60 * 10)

var HomePage []byte


type Post struct{
    Title string
    Url string
}

type Subreddit struct{
    Name string
    LastUpdated int64
    Posts []Post
}

var Subreddits map[string]Subreddit


// takes an input sub populates that sub in the global Subreddits variable
// returns true or false if the process was succesfull or not
func populate (sub *Subreddit) bool{
    url := "https://www.reddit.com/r/" + sub.Name + "/.json?limit=100"

    // We need to create a client to add header
    client := &http.Client{}
    req , err := http.NewRequest("GET",url,nil)
    if err != nil{
        // This is a step that basically can't fail.
        // If it fails there has to be a huge underlying
        // problem so I didn't bother with fixing it.
        panic(err)
    }

    req.Header.Add("User-Agent",uAgent)
    resp , err := client.Do(req)
    defer resp.Body.Close()
    if err != nil{
      // If the request timesout or some other error occurs,
      // We repond by returning false
      log.Println(err)
      return false
    }

    // Check if the json formed is good.
    rBody, err := ioutil.ReadAll(resp.Body)
    if err != nil{
      log.Println(err)
      return false
    }

    /* 
       Here, we begin parsing the JSON
       The schema of the JSON is as follows
       {
        "kind":"Listing",
        "data":
                {
                 "after":"id of previous post",
                 "dist":"no of posts",
                 ...
                 children:
                            {
                                [
                                {
                                 "title" :  "cleopatrick - daphne did it [Rock] (2018)" // Example
                                 "url" : "https://youtu.be/ickvWCkE9Wk?t=1m32s"
                                 ...
                                },
                                {...},
                                ...]
                            }
                }
       }
       So first we must extract the json from data
       then another one from the children
       which gives us an array of Posts
    */

    var buffer interface{}
    json.Unmarshal(rBody,&buffer)
    root := buffer.(map[string]interface{})
    body := root["data"]

    // Apparently we have to assert the type of 
    // every unmarshal before we can use it.
    // Note that the type signature of children is []interface
    // and not a map, because it an array.
    children := body.(map[string]interface{})["children"].([]interface{})

    // Clear the posts of the sub
    sub.Posts = sub.Posts[:0]

    for _, rawPost := range children{
        post := rawPost.(map[string]interface{})
        data := post["data"].(map[string]interface{})
        title := data["title"].(string)
        url := data["url"].(string)
        stickied := data["stickied"].(bool)
        if stickied{
            continue
        }
        sub.Posts  = append(sub.Posts,Post{title,url})
    }

    sub.LastUpdated = int64(time.Now().Unix())

    return true
}

func getYouTubeVideos(w http.ResponseWriter , r *http.Request){
    // An end point's all

    // Check if the sub parameter is present
    paramSub := r.FormValue("sub")
    if paramSub == ""{
        w.Write(errors.EmptyParameter)
        return
    }

    // Check if the subreddit string is alphanumeric
    if !helpers.IsSubValid(paramSub){
      w.Write(errors.InvalidSub)
      return
    }

    // Check if the subreddit is already present in the global variable
    sub, ok := Subreddits[paramSub]
    if !ok{
        sub.Name = paramSub
        sub.LastUpdated = int64(time.Now().Unix())
        Subreddits[paramSub] = sub
        if !populate(&sub){
            delete(Subreddits,paramSub)
            w.Write(errors.InvalidSub)
            return
        }
        Subreddits[paramSub] = sub
    }else{
        if int64(time.Now().Unix()) > (sub.LastUpdated + refreshTime) {
            if !populate(&sub){
                w.Write(errors.CouldNotPopulate)
                return
            }
            Subreddits[paramSub] = sub
        }
    }
    // default 25 limit unless a limit is explicitly provided
    limit := 25
    paramLimit := r.FormValue("limit")
    if paramLimit != ""{
        limit, _ = strconv.Atoi(paramLimit)
    }
    limit = helpers.IntMin(limit,len(sub.Posts))

    // If another user is requesting to populate the same sub,
    // we return an error.
    if sub.Posts == nil{
        w.Write(errors.SimilarRequestProcessing)
        return
    }

    b, _ := json.Marshal(Subreddit{sub.Name,sub.LastUpdated,sub.Posts[0:limit]})
    
    w.Write(b)
}


func home(w http.ResponseWriter, r *http.Request){
  // Homepage endpoint's all
  w.Write(HomePage)
}


func server(){
  // Highest level of abstraction of the server
  LogFile , err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  defer LogFile.Close()
  if err != nil{
    panic(err)
  }
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
  http.HandleFunc("/youtubevideos",getYouTubeVideos)
  http.HandleFunc("/",home)
  LoggingHandler := handlers.LoggingHandler(LogFile, http.DefaultServeMux)
  if err := http.ListenAndServe(":8080",handlers.CompressHandler(LoggingHandler)); err != nil {
    panic(err)
  }
}

  
func main(){
    HomePage , _ = ioutil.ReadFile("./static/index.html")
    Subreddits = make(map[string]Subreddit)
    fmt.Println(url)
    server()
}

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

var limit = 100
var timeDelay time.Duration = 2

var url = "http://www.reddit.com/r/listentothis/.json?limit=" + strconv.Itoa(limit)
var uAgent = "rjpj's obscure ageint v0.0.1"
var refreshTime = int64(60 * 10)

var HomePage []byte

type Post struct{
    Title string
    Url string
}



var posts []Post

type Subreddit struct{
    Name string
    LastUpdated int64
    Posts []Post
}

var Subreddits map[string]Subreddit

func populate (sub *Subreddit) bool{
    url := "https://www.reddit.com/r/" + sub.Name + "/.json?limit=100"
    client := &http.Client{}
    req , err := http.NewRequest("GET",url,nil)
    if err != nil{
        panic(err)
    }
    req.Header.Add("User-Agent",uAgent)

    resp , err := client.Do(req)
    if err != nil{
      log.Println(err)
      return false
    }

    defer resp.Body.Close()

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
    posts = posts[:0]
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
    paramSub := r.FormValue("sub")
    if paramSub == ""{
        w.Write(errors.EmptyParameter)
        return
    }
    if !helpers.IsSubValid(paramSub){
      w.Write(errors.InvalidSub)
      return
    }

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

    limit := 25
    paramLimit := r.FormValue("limit")
    if paramLimit != ""{
        limit, _ = strconv.Atoi(paramLimit)
    }
    limit = helpers.IntMin(limit,len(sub.Posts))

    if sub.Posts == nil{
        w.Write(errors.SimilarRequestProcessing)
    }

    b, _ := json.Marshal(Subreddit{sub.Name,sub.LastUpdated,sub.Posts[0:limit]})
    
    w.Write(b)
}

func list(w http.ResponseWriter , r *http.Request){
  if len(posts) == 0 {
      w.Write(errors.NotInitializedError)
      return
  }
  limit := 25
  paramLimit := r.FormValue("limit")
  if paramLimit != ""{
      limit, _ = strconv.Atoi(paramLimit)
  }
  limit = helpers.IntMin(limit,len(posts))
  b, _ := json.Marshal(posts[0:limit])
  w.Write(b)
}

func home(w http.ResponseWriter, r *http.Request){
  w.Write(HomePage)
}


func server(){
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
  http.HandleFunc("/youtubevideos",getYouTubeVideos)
  http.HandleFunc("/",home)
  if err := http.ListenAndServe(":8080",handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)); err != nil {
    panic(err)
  }
}

  

func main(){
    HomePage , _ = ioutil.ReadFile("./static/test.html")
    Subreddits = make(map[string]Subreddit)
    fmt.Println(url)
    server()
}

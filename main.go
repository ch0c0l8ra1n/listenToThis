package main

import (
  "net/http"
  "fmt"
  "io/ioutil"
  "time"
  "encoding/json"
  "strconv"
  "./helpers"
)

var limit = 100
var timeDelay time.Duration = 2

var url = "http://www.reddit.com/r/listentothis/.json?limit=" + strconv.Itoa(limit)
var uAgent = "rjpj's listenToThis server v0.0.1"

type Post struct{
    Title string
    Url string
}

type ManyPosts struct{
    Children []Post
}

var posts []Post

type Error struct{
    ErrorCode int
    ErrorMsg string
}

var notInitializedError,_ = json.Marshal(
                        Error { 0 , 
                        "No videos currently indexed."})


func list(w http.ResponseWriter , r *http.Request){
    if len(posts) == 0 {
        w.Write(notInitializedError)
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


func server(){
  http.HandleFunc("/list",list)
  if err := http.ListenAndServe(":8080",nil); err != nil {
    panic(err)
  }
}


func cacheServer(){

  // sorta object oriented approach
  // afaik necessary for headers
  client := &http.Client{}

  req , err := http.NewRequest("GET",url,nil)
    if err != nil{
      panic(err)
  }

  req.Header.Add("User-Agent",uAgent)


  // Updates our collection of entries every 
  // few seconds
  for ;; {
    
    resp , err := client.Do(req)
    if err != nil{
      panic(err)
    }

    defer resp.Body.Close()

    rBody, err := ioutil.ReadAll(resp.Body)
    if err != nil{
      panic(err)
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
    fmt.Println(len(children))
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
        posts  = append(posts,Post{title,url})
        //fmt.Println(posts)
    }

    // rest duration
    time.Sleep(timeDelay * time.Second)

  }
  

}

func main(){
    fmt.Println(url)
    go cacheServer()
    server()
}

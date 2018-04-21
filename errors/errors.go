package errors

import (
    "encoding/json"
)

type Error struct{
    ErrorCode int64
    ErrorMsg string
}

var NotInitializedError,_ = json.Marshal(
                        Error { 0 , 
                        "No videos currently indexed."})

var EmptyParameter , _ = json.Marshal(
                        Error { 1,
                        "Required parameters were not received."})

var CouldNotPopulate , _ = json.Marshal(
                           Error { 2,
                           "Couldn't populate the list."})

var InvalidSub , _ = json.Marshal(
                           Error { 3,
                           "The Subreddit is invalid."})
var SimilarRequestProcessing , _ = json.Marshal(
                            Error{ 4,
                            "A similar request is already processing. Please try again in a few seconds."})
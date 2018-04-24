package helpers

import "regexp"

func IntMin(a int , b int) int{
    // Default math.Min function only works on floats
    if b<a {
        return b
    }
    return a
}

var IsSubValid = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
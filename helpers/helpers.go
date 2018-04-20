package helpers

func IntMin(a int , b int) int{
    // Default math.Min function only works on floats
    if b<a {
        return b
    }
    return a
}
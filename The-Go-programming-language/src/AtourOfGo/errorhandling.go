package main

import (
    "fmt"
    "math"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
    return fmt.Sprintf("cannot Sqrt negative number: %f", e)
}


func Sqrt(f float64) (float64, error) {
    if f < 0 {
     return 0, ErrNegativeSqrt(f)   
    }
    
    z := float64(1)
    newZ := float64(1)
    
    for{
        newZ = z - (z * z - f) / (2 * f)
        
        if newZ - z < math.Abs(0.000001) {
    		return newZ, nil
        }else{
          z = newZ   
        }            
	}

}

func main() {
    fmt.Println(Sqrt(2))
    fmt.Println(Sqrt(-2))

}

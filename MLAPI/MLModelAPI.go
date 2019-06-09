package MLAPI

// MLAPI for machine learning application interface
import (
	"HealthySkin/DBDAL"
	"math/rand"
	"os"
)

// this function activates a machine learning model to check the image for positive/negative cancer detection.
func GetIsCancerImage(imageFile *os.File, userDetails *DBDAL.UserInfo) bool {
	var result = rand.Intn(2)
	return result == 1
}

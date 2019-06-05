package MLAPI

// MLAPI for machine learning application interface
import (
	"HealthySkin/DBDAL"
	"math/rand"
)

func getIsCancerImage(imageBytes []byte, userDetails DBDAL.UserInfo) bool {
	//need to run python script and send him user details and image bytes in order to know if its cancer or not
	var result = rand.Intn(2)
	return result == 1
}

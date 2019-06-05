package DBDAL

//named data base data access layer

type UserInfo struct {
	id     int
	name   string
	age    int
	weight int
}

var userDetailsDict = make(map[int]UserInfo)

func saveUserInfoDetails(userInfo UserInfo) {
	userDetailsDict[userInfo.id] = userInfo
}

func getUserDetailsById(id int) *UserInfo {
	value, ok := userDetailsDict[id]
	if ok {
		return &value
	} else {
		return nil
	}
}

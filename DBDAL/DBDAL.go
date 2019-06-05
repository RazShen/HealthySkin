package DBDAL

//named data base data access layer

type UserInfo struct {
	id          int
	md5Password string
	name        string
	age         int
	weight      int
}

var userDetailsDict = make(map[int]UserInfo)

func saveUserInfoDetails(userInfo UserInfo) {
	userDetailsDict[userInfo.id] = userInfo
}

func getUserDetailsById(id int, md5Password string) *UserInfo {

	value, ok := userDetailsDict[id]
	if ok {
		if value.md5Password == md5Password {
			return &value
		}
	}
	return nil

}

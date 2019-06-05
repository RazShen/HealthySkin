package DBDAL

//named data base data access layer

type UserInfo struct {
	MD5Password string
	UserName    string
	Age         int
	Weight      int
}

var userDetailsDict = make(map[string]UserInfo)

func SaveUserInfoDetails(userInfo UserInfo) {
	userDetailsDict[userInfo.UserName] = userInfo
}

func GetUserDetailsById(userName string, md5Password string) *UserInfo {

	value, ok := userDetailsDict[userName]
	if ok {
		if value.MD5Password == md5Password {
			return &value
		}
	}
	return nil

}

package mongodb

type DB struct {
	Path    string `json:"path"`
	DBName  string `json:"db_name"`
	DBUser  string `json:"db_user"`
	DBPass  string `json:"db_pass"`
	MaxPool int    `json:"max_pool"`
}

var MongoConfig = DB{}

func CheckAndInitServiceConnection() {
	if service.baseSession == nil {
		service.URL = MongoConfig.Path
		service.DbUser = MongoConfig.DBUser
		service.DbPass = MongoConfig.DBPass
		err := service.New()
		if err != nil {
			logDB.Errorf("disconnected from %s", MongoConfig.Path)
			panic(err)
		}
	}
}

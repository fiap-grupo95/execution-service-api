package database

func ConnectDatabase() (*MongoDBConfig, error) {
	return NewMongoDBFromEnv()
}

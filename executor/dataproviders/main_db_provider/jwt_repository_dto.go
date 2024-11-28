package main_db_provider

type JWTModel struct{
	Id int32
	Hash string
	Guid string
	Expiry int32
	LastLoginIp string
}
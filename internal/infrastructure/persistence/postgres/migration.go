package postgres


func(ps *postgresDB)Migration()error{
   return  ps.DB.AutoMigrate(
	
   )
}
func TestInsert(t *testing.T) {
	// 創建模擬 db 和回應
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	// 創建 GORM DB 實例
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create gorm db: %v", err)
	}

	// 開始一個模擬交易，期望事務開始的調用
	mock.ExpectBegin()

	// 執行 Insert 函式
	user := []User{User{Name: "testuser", Age: 20}, User{Name: "wesley1", Age: 30}}
	//mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user` (`name`,`age`) values (?,?)")).WithArgs(user.Name, user.Age).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO `user`").WithArgs(user[0].Name, user[0].Age, user[1].Name, user[1].Age).WillReturnResult(sqlmock.NewResult(2, 2))
	mock.ExpectCommit()

	rowsAffected, err := Insert(gormDB, user)
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	// 確認回應和預期的一致
	if rowsAffected != 2 {
		t.Errorf("expected 1 row affected, but got %v", rowsAffected)
	}

	// 確認所有的 expectation 都已經匹配到
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

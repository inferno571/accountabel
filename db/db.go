package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
	"github.com/yincongcyincong/MuseBot/conf"
	"github.com/yincongcyincong/MuseBot/logger"
	botUtils "github.com/yincongcyincong/MuseBot/utils"
)

var (
	sqlite3TableSQLs = map[string]string{
		"users": `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id varchar(100) NOT NULL DEFAULT '0',
			update_time INTEGER NOT NULL DEFAULT '0',
			token INTEGER NOT NULL DEFAULT '0',
			avail_token INTEGER NOT NULL DEFAULT 0,
			create_time INTEGER NOT NULL DEFAULT '0',
			from_bot VARCHAR(255) NOT NULL DEFAULT '',
			llm_config TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_users_user_id ON users(user_id);
	`,
		"records": `
		CREATE TABLE IF NOT EXISTS records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id varchar(100) NOT NULL DEFAULT '0',
			question TEXT NOT NULL,
			answer TEXT NOT NULL,
			content TEXT NOT NULL,
			create_time INTEGER NOT NULL DEFAULT '0',
			update_time INTEGER NOT NULL DEFAULT '0',
			is_deleted INTEGER NOT NULL DEFAULT '0',
			token INTEGER NOT NULL DEFAULT 0,
			mode VARCHAR(100) NOT NULL DEFAULT '',
			record_type INTEGER NOT NULL DEFAULT 0, -- SQLite中用INTEGER代替tinyint
			from_bot VARCHAR(255) NOT NULL DEFAULT ''
		);
		CREATE INDEX IF NOT EXISTS idx_records_user_id ON records(user_id);
		CREATE INDEX IF NOT EXISTS idx_records_create_time ON records(create_time);
	`,
		"rag_files": `
		CREATE TABLE IF NOT EXISTS rag_files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			file_name VARCHAR(255) NOT NULL DEFAULT '',
			file_md5 VARCHAR(255) NOT NULL DEFAULT '',
			vector_id TEXT NOT NULL DEFAULT '',
			create_time INTEGER NOT NULL DEFAULT '0',
			update_time INTEGER NOT NULL DEFAULT '0',
			is_deleted INTEGER NOT NULL DEFAULT '0',
			from_bot VARCHAR(255) NOT NULL DEFAULT ''
		);
	`,
		"cron": `
		CREATE TABLE IF NOT EXISTS cron (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			cron_name VARCHAR(255) NOT NULL DEFAULT '',
			cron VARCHAR(255) NOT NULL DEFAULT '',
			target_id TEXT NOT NULL,
			group_id TEXT NOT NULL,
			command VARCHAR(255) NOT NULL DEFAULT '',
			prompt TEXT NOT NULL,
			status INTEGER NOT NULL DEFAULT 1, -- 0:disable 1:enable
			cron_job_id INTEGER NOT NULL DEFAULT '0',
			create_time INTEGER NOT NULL DEFAULT '0',
			update_time INTEGER NOT NULL DEFAULT '0',
			is_deleted INTEGER NOT NULL DEFAULT '0',
			from_bot VARCHAR(255) NOT NULL DEFAULT '',
		    type VARCHAR(255) NOT NULL DEFAULT '',
		    create_by VARCHAR(255) NOT NULL DEFAULT ''
		);
	`,
		"user_profiles": `
		CREATE TABLE IF NOT EXISTS user_profiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id VARCHAR(100) NOT NULL UNIQUE,
			platform VARCHAR(50) NOT NULL DEFAULT '',
			timezone VARCHAR(100) NOT NULL DEFAULT 'UTC',
			scheduled_checkin_time VARCHAR(10) NOT NULL DEFAULT '',
			cron_job_id INTEGER NOT NULL DEFAULT 0,
			create_time INTEGER NOT NULL DEFAULT 0,
			update_time INTEGER NOT NULL DEFAULT 0,
			addictions TEXT NOT NULL DEFAULT '',
			link_code VARCHAR(20) NOT NULL DEFAULT '',
			link_code_expires INTEGER NOT NULL DEFAULT 0
		);
		CREATE INDEX IF NOT EXISTS idx_user_profiles_user_id ON user_profiles(user_id);
	`,
		"check_ins": `
		CREATE TABLE IF NOT EXISTS check_ins (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id VARCHAR(100) NOT NULL DEFAULT '',
			timestamp INTEGER NOT NULL DEFAULT 0,
			craving_level INTEGER NOT NULL DEFAULT 0,
			relapse_status INTEGER NOT NULL DEFAULT 0,
			notes TEXT NOT NULL DEFAULT ''
		);
		CREATE INDEX IF NOT EXISTS idx_check_ins_user_id ON check_ins(user_id);
		CREATE INDEX IF NOT EXISTS idx_check_ins_timestamp ON check_ins(timestamp);
	`,
		"streaks": `
		CREATE TABLE IF NOT EXISTS streaks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id VARCHAR(100) NOT NULL UNIQUE,
			current_streak INTEGER NOT NULL DEFAULT 0,
			highest_streak INTEGER NOT NULL DEFAULT 0,
			last_check_in INTEGER NOT NULL DEFAULT 0,
			update_time INTEGER NOT NULL DEFAULT 0
		);
		CREATE INDEX IF NOT EXISTS idx_streaks_user_id ON streaks(user_id);
	`,
		"slips_log": `
		CREATE TABLE IF NOT EXISTS slips_log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id VARCHAR(100) NOT NULL,
			slip_date INTEGER NOT NULL DEFAULT 0,
			previous_streak_days INTEGER NOT NULL DEFAULT 0,
			notes TEXT NOT NULL DEFAULT ''
		);
		CREATE INDEX IF NOT EXISTS idx_slips_log_user_id ON slips_log(user_id);
	`,
		"craving_logs": `
		CREATE TABLE IF NOT EXISTS craving_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id VARCHAR(100) NOT NULL,
			intensity INTEGER NOT NULL DEFAULT 0,
			trigger_context TEXT NOT NULL DEFAULT '',
			tags TEXT NOT NULL DEFAULT '[]',
			action_taken VARCHAR(100) NOT NULL DEFAULT '',
			logged_at INTEGER NOT NULL DEFAULT 0
		);
		CREATE INDEX IF NOT EXISTS idx_craving_logs_user_id ON craving_logs(user_id);
	`,
		"coping_tools": `
		CREATE TABLE IF NOT EXISTS coping_tools (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id VARCHAR(100) NOT NULL,
			tool_type VARCHAR(100) NOT NULL DEFAULT '',
			content TEXT NOT NULL DEFAULT ''
		);
		CREATE INDEX IF NOT EXISTS idx_coping_tools_user_id ON coping_tools(user_id);
	`,
		"emergency_contacts": `
		CREATE TABLE IF NOT EXISTS emergency_contacts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id VARCHAR(100) NOT NULL,
			contact_name TEXT NOT NULL DEFAULT '',
			contact_phone TEXT NOT NULL DEFAULT '',
			relationship VARCHAR(255) NOT NULL DEFAULT ''
		);
		CREATE INDEX IF NOT EXISTS idx_emergency_contacts_user_id ON emergency_contacts(user_id);
	`,
	}

	mysqlInitializeSQLs = []string{
		// 1. users 表 (嵌入索引)
		`
       CREATE TABLE IF NOT EXISTS users (
          id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
          user_id varchar(100) NOT NULL DEFAULT '0',
          update_time INT(10) NOT NULL DEFAULT 0,
          token BIGINT NOT NULL DEFAULT 0,
          avail_token BIGINT NOT NULL DEFAULT 0,
           create_time INT(10) NOT NULL DEFAULT 0,
           from_bot VARCHAR(255) NOT NULL DEFAULT '',
           llm_config TEXT NOT NULL,
           
           -- 嵌入索引：idx_users_user_id
           INDEX idx_users_user_id (user_id)
       ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`,
		// 2. records 表 (嵌入索引)
		`
       CREATE TABLE IF NOT EXISTS records (
          id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
          user_id varchar(100) NOT NULL DEFAULT '0',
          question MEDIUMTEXT NOT NULL,
          answer MEDIUMTEXT NOT NULL,
          content MEDIUMTEXT NOT NULL,
          create_time INT(10) NOT NULL DEFAULT 0,
           update_time INT(10) NOT NULL DEFAULT 0,
          is_deleted INT(10) NOT NULL DEFAULT 0,
          token INT(10) NOT NULL DEFAULT 0,
           mode VARCHAR(100) NOT NULL DEFAULT '',
           record_type tinyint(1) NOT NULL DEFAULT 0 COMMENT '0:text, 1:image 2:video 3: web',
           from_bot VARCHAR(255) NOT NULL DEFAULT '',
           
           -- 嵌入索引：idx_records_user_id 和 idx_records_create_time
           INDEX idx_records_user_id (user_id),
           INDEX idx_records_create_time (create_time)
       ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`,
		// 3. rag_files 表 (无额外索引，仅PRIMARY KEY)
		`CREATE TABLE IF NOT EXISTS rag_files (
          id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
          file_name VARCHAR(255) NOT NULL DEFAULT '',
          file_md5 VARCHAR(255) NOT NULL DEFAULT '',
          vector_id TEXT NOT NULL,
          create_time INT(10) NOT NULL DEFAULT 0,
          update_time INT(10) NOT NULL DEFAULT 0,
          is_deleted INT(10) NOT NULL DEFAULT 0,
          from_bot VARCHAR(255) NOT NULL DEFAULT ''
       ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`,
		// 4. cron 表 (无额外索引，仅PRIMARY KEY)
		`CREATE TABLE IF NOT EXISTS cron (
          id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
          cron_name VARCHAR(255) NOT NULL DEFAULT '',
          cron VARCHAR(255) NOT NULL DEFAULT '',
          target_id TEXT NOT NULL,
          group_id TEXT NOT NULL,
          command VARCHAR(255) NOT NULL DEFAULT '',
          prompt TEXT NOT NULL,
          status tinyint(1) NOT NULL DEFAULT 1 COMMENT '0:disable 1:enable',
          cron_job_id INT(10) NOT NULL DEFAULT 0,
          create_time INT(10) NOT NULL DEFAULT 0,
          update_time INT(10) NOT NULL DEFAULT 0,
          is_deleted INT(10) NOT NULL DEFAULT 0,
          from_bot VARCHAR(255) NOT NULL DEFAULT '',
          type VARCHAR(255) NOT NULL DEFAULT '',
    	  create_by VARCHAR(255) NOT NULL DEFAULT ''
       ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`,
		// Add the new recovery tables to MySQL
		`CREATE TABLE IF NOT EXISTS slips_log (
          id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
          user_id VARCHAR(100) NOT NULL,
          slip_date INT(10) NOT NULL DEFAULT 0,
          previous_streak_days INT(10) NOT NULL DEFAULT 0,
          notes TEXT,
          INDEX idx_slips_log_user_id (user_id)
       ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS craving_logs (
          id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
          user_id VARCHAR(100) NOT NULL,
          intensity INT NOT NULL DEFAULT 0,
          trigger_context TEXT,
          tags TEXT,
          action_taken VARCHAR(100) NOT NULL DEFAULT '',
          logged_at INT(10) NOT NULL DEFAULT 0,
          INDEX idx_craving_logs_user_id (user_id)
       ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS coping_tools (
          id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
          user_id VARCHAR(100) NOT NULL,
          tool_type VARCHAR(100) NOT NULL DEFAULT '',
          content TEXT,
          INDEX idx_coping_tools_user_id (user_id)
       ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS emergency_contacts (
          id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
          user_id VARCHAR(100) NOT NULL,
          contact_name TEXT NOT NULL,
          contact_phone TEXT NOT NULL,
          relationship VARCHAR(255) NOT NULL DEFAULT '',
          INDEX idx_emergency_contacts_user_id (user_id)
       ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
	}
)

var (
	DB *sql.DB
)

type DailyStat struct {
	Date     string `json:"date"`
	NewCount int    `json:"new_count"`
}

func InitTable() {
	var err error
	if _, err = os.Stat(botUtils.GetAbsPath("data")); os.IsNotExist(err) {
		// if dir don't exist, create it.
		err := os.MkdirAll(botUtils.GetAbsPath("data"), 0755)
		if err != nil {
			logger.Fatal("create direction fail:", "err", err)
			return
		}
		logger.Info("✅ create direction success")
	}

	dbType := conf.BaseConfInfo.DBType
	if dbType == "sqlite3" {
		dbType = "sqlite"
	}
	DB, err = sql.Open(dbType, conf.BaseConfInfo.DBConf)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// init table
	switch conf.BaseConfInfo.DBType {
	case "sqlite3", "sqlite":
		err = initializeSqlite3Table(DB)
		if err != nil {
			logger.Fatal("create sqlite table fail", "err", err)
		}
	case "mysql":
		err = initializeMySQLTables(DB)
		if err != nil {
			logger.Fatal("create mysql table fail", "err", err)
		}
	}

	InsertRecord(context.Background())

	logger.Info("db initialize successfully")
}

func initializeMySQLTables(db *sql.DB) error {
	for i, sqlStr := range mysqlInitializeSQLs {
		_, err := db.Exec(sqlStr)
		if err != nil {
			logger.Error("check table fail", "err", err)
			return fmt.Errorf("execute SQL batch %d fail: %v\nSQL: %s", i+1, err, sqlStr)
		}
	}

	return nil
}

// initializeSqlite3Table check table exist or not.
func initializeSqlite3Table(db *sql.DB) error {
	for tableName, createSQL := range sqlite3TableSQLs {
		_, err := db.Exec(createSQL)
		if err != nil {
			logger.Error("check table fail", "tableName", tableName, "err", err)
			return fmt.Errorf("create table %s fail: %v", tableName, err)
		}
	}

	// Add new columns to user_profiles if they do not exist (ignore errors if columns already exist)
	db.Exec("ALTER TABLE user_profiles ADD COLUMN addictions TEXT NOT NULL DEFAULT ''")
	db.Exec("ALTER TABLE user_profiles ADD COLUMN link_code VARCHAR(20) NOT NULL DEFAULT ''")
	db.Exec("ALTER TABLE user_profiles ADD COLUMN link_code_expires INTEGER NOT NULL DEFAULT 0")
	db.Exec("ALTER TABLE user_profiles ADD COLUMN reasons_to_quit TEXT NOT NULL DEFAULT ''")
	db.Exec("ALTER TABLE user_profiles ADD COLUMN warning_signs TEXT NOT NULL DEFAULT ''")
	db.Exec("ALTER TABLE user_profiles ADD COLUMN cost_per_use DECIMAL(10,2) NOT NULL DEFAULT 0.0")
	db.Exec("ALTER TABLE user_profiles ADD COLUMN usage_frequency VARCHAR(50) NOT NULL DEFAULT ''")
	db.Exec("ALTER TABLE user_profiles ADD COLUMN google_user_id VARCHAR(255) NOT NULL DEFAULT ''")
	db.Exec("ALTER TABLE user_profiles ADD COLUMN google_name VARCHAR(255) NOT NULL DEFAULT ''")

	return nil
}

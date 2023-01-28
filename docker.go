package dst

import (
	"fmt"
	"os"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DockerPg 测试辅助函数。 利用 dockertest 生成一次性pg实例。
// version , 可选参数，用于指定pg版本， 默认为 "13"
// 返回 gorm 连接字符串， 以及用于清理此pg实例的 cleanup 函数
func DockerPg(version ...string) (string, func()) {
	pool, err := dockertest.NewPool("")
	Chk(err)

	ver := "13" // default version
	if len(version) > 0 {
		ver = version[0]
	}

	const testDbName = "test"
	const testPasswd = "test"

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        ver,
			Env: []string{
				"POSTGRES_PASSWORD=" + testPasswd,
				"POSTGRES_DB=" + testDbName,
			},
		}, func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.NeverRestart()
		})

	Chk(err)

	port := resource.GetPort("5432/tcp")
	host := "localhost"
	if s := os.Getenv("DOCKER_PG_HOST"); s != "" {
		host = s
	}
	conStr := fmt.Sprintf("host=%s port=%s user=postgres dbname=%s password=%s sslmode=disable",
		host,
		port,
		testDbName,
		testPasswd,
	)
	err = pool.Retry(func() error {
		cfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}
		gdb, err := gorm.Open(postgres.Open(conStr), cfg)
		if err != nil {
			return err
		}

		db, err := gdb.DB()
		if err != nil {
			return err
		}

		return db.Ping()
	})
	Chk(err)

	return conStr, func() {
		err := resource.Close()
		Chk(err)
	}
}

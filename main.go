package main

import (
	"context"
	"database/sql"
	"log"
	"microdata/kemendagri/bumd/controller"
	controller_bumd "microdata/kemendagri/bumd/controller/bumd"
	controller_mst "microdata/kemendagri/bumd/controller/master"
	"microdata/kemendagri/bumd/handler/http"
	"microdata/kemendagri/bumd/handler/http/bumd"
	"microdata/kemendagri/bumd/handler/http/configs"
	"microdata/kemendagri/bumd/handler/http/http_util"
	http_master "microdata/kemendagri/bumd/handler/http/master"
	_deliveryMiddleware "microdata/kemendagri/bumd/handler/http/middleware"
	"microdata/kemendagri/bumd/utils"
	"os"
	"strconv"
	"strings"
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis/v3"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	_ "microdata/kemendagri/bumd/docs" // load API Docs files (Swagger)
)

var serverName,
	serverUrl,
	serverReadTimeout,
	dbServerUrl,
	dbServerUrlMstData,
	jwtSecertKey,
	alwOrg,
	redisHost,
	redisUsername,
	redisPassword,
	minioEndpoint,
	minioAccessKey,
	minioSecretKey,
	minioBucketName string

var redisPort int

var db *sql.DB
var pgxConn,
	pgxConnMstData *pgxpool.Pool
var driver database.Driver
var migration *migrate.Migrate
var jwtMgr *utils.JWTManager
var vld *validator.Validate
var redisCl *redis.Storage
var logger *logrus.Logger
var err error
var minioConn *utils.MinioConn

func init() {
	// Server Env
	serverName = os.Getenv("SERVER_NAME")
	if serverName == "" {
		exitf("SERVER_NAME env is required")
	}
	serverUrl = os.Getenv("SERVER_URL")
	if serverUrl == "" {
		exitf("SERVER_URL env is required")
	}
	serverReadTimeout = os.Getenv("SERVER_READ_TIMEOUT")
	if serverReadTimeout == "" {
		exitf("SERVER_READ_TIMEOUT env is required")
	}

	// JWT Env
	jwtSecertKey = os.Getenv("JWT_SECRET_KEY")
	if jwtSecertKey == "" {
		exitf("JWT_SECRET_KEY env is required")
	}

	// CORS
	alwOrg = os.Getenv("SIPD_CORS_WHITELISTS")
	if alwOrg == "" {
		exitf("SIPD_CORS_WHITELISTS config is required")
	}

	// Databse Env
	dbServerUrl = os.Getenv("DB_SERVER_URL")
	if dbServerUrl == "" {
		exitf("DB_SERVER_URL config is required")
	}
	dbServerUrlMstData = os.Getenv("DB_SERVER_URL_MST_DATA")
	if dbServerUrlMstData == "" {
		exitf("DB_SERVER_URL_MST_DATA config is required")
	}

	// Redis
	redisHost = os.Getenv("REDIS_HOST")
	if redisHost == "" {
		exitf("REDIS_HOST env is required")
	}

	redisPortStr := os.Getenv("REDIS_PORT")
	if os.Getenv("REDIS_PORT") == "" {
		exitf("REDIS_PORT env is required")
	}
	redisPort, err = strconv.Atoi(redisPortStr)
	if err != nil {
		exitf("REDIS_PORT env is invalid")
	}
	/*redisUsername = os.Getenv("REDIS_USERNAME")
	if redisUsername == "" {
		exitf("REDIS_USERNAME env is required")
	}*/
	redisPassword = os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		exitf("REDIS_PASSWORD env is required")
	}

	// Minio
	minioEndpoint = os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		exitf("MINIO_ENDPOINT config is required")
	}
	minioAccessKey = os.Getenv("MINIO_ACCESS_KEY")
	if minioAccessKey == "" {
		exitf("MINIO_ACCESS_KEY config is required")
	}
	minioSecretKey = os.Getenv("MINIO_SECRET_KEY")
	if minioSecretKey == "" {
		exitf("MINIO_SECRET_KEY config is required")
	}
	minioBucketName = os.Getenv("MINIO_BUCKET_NAME")
	if minioBucketName == "" {
		exitf("MINIO_BUCKET_NAME config is required")
	}
}

func dbConnection() {
	var maxConnLifetime, maxConnIdleTime time.Duration
	var maxPoolConn int32
	maxConnLifetime = 5 * time.Minute
	maxConnIdleTime = 2 * time.Minute
	maxPoolConn = 1000

	var cfg,
		cfgMstData *pgxpool.Config

	// auth
	cfg, err = pgxpool.ParseConfig(dbServerUrl + " application_name=" + serverName)
	if err != nil {
		exitf("Unable to create db pool config auth %v\n", err)
	}
	cfg.MaxConns = maxPoolConn            // Maximum total connections in the pool
	cfg.MaxConnLifetime = maxConnLifetime // Maximum lifetime of a connection
	cfg.MaxConnIdleTime = maxConnIdleTime // Maximum time a connection can be idle
	pgxConn, err = pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		exitf("Unable to connect to database auth: %v\n", err)
	}

	// pegawai
	// cfgPegawai, err = pgxpool.ParseConfig(dbServerUrlPegawai + " application_name=" + serverName)
	// if err != nil {
	// 	exitf("Unable to create db pool config pegawai %v\n", err)
	// }
	// cfgPegawai.MaxConns = maxPoolConn            // Maximum total connections in the pool
	// cfgPegawai.MaxConnLifetime = maxConnLifetime // Maximum lifetime of a connection
	// cfgPegawai.MaxConnIdleTime = maxConnIdleTime // Maximum time a connection can be idle
	// pgxConnPegawai, err = pgxpool.NewWithConfig(context.Background(), cfgPegawai)
	// if err != nil {
	// 	exitf("Unable to connect to database pegawai: %v\n", err)
	// }

	// mst_data
	cfgMstData, err = pgxpool.ParseConfig(dbServerUrlMstData + " application_name=" + serverName)
	if err != nil {
		exitf("Unable to create db pool config mst_data %v\n", err)
	}
	cfgMstData.MaxConns = maxPoolConn            // Maximum total connections in the pool
	cfgMstData.MaxConnLifetime = maxConnLifetime // Maximum lifetime of a connection
	cfgMstData.MaxConnIdleTime = maxConnIdleTime // Maximum time a connection can be idle
	pgxConnMstData, err = pgxpool.NewWithConfig(context.Background(), cfgMstData)
	if err != nil {
		exitf("Unable to connect to database mst_data: %v\n", err)
	}

	// Minio
	minioConn, err = utils.NewMinIOConn(
		minioEndpoint,
		minioAccessKey,
		minioSecretKey,
		minioBucketName,
	)
	if err != nil {
		exitf("Unable to connect to minio: %v\n", err)
	}
}

// func dbConnectionSimple() {
// 	pgxConn, err = pgxpool.New(context.Background(), dbServerUrl+" application_name="+serverName)
// 	if err != nil {
// 		exitf("Unable to connect to database auth: %v\n", err)
// 	}

// 	// mst_data
// 	pgxConnMstData, err = pgxpool.New(context.Background(), dbServerUrlMstData+" application_name="+serverName)
// 	if err != nil {
// 		exitf("Unable to connect to database mst_data: %v\n", err)
// 	}
// }

func redisConnection() {
	redisCl = redis.New(redis.Config{
		Host:      redisHost,
		Port:      redisPort,
		Username:  redisUsername,
		Password:  redisPassword,
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		// PoolSize:  10 * runtime.GOMAXPROCS(0),
	})
}

// func setupLogger() {
// 	// Create a new instance of Logrus logger
// 	logger = logrus.New()

// 	// Configure the Graylog hook
// 	graylogHost := "127.0.0.1:12201" // Replace with your Graylog server's address and port
// 	hook := graylog.NewGraylogHook(graylogHost, map[string]interface{}{
// 		"environment": "production", // Custom fields
// 		"app":         "my-golang-app",
// 	})

// 	// Add the Graylog hook to the logger
// 	logger.Hooks.Add(hook)

// 	// Set Logrus Formatter (optional)
// 	logger.SetFormatter(&logrus.JSONFormatter{})
// }

// @title						BUMD Service
// @version					1.0
// @description				BUMD Service Rest API.
// @termsOfService				http://swagger.io/terms/
// @contact.name				API Support
// @contact.email				lifelinejar@mail.com
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @BasePath					/bumd/
func main() {
	// checking number of cpu
	// numOfCpu := runtime.NumCPU()
	// log.Printf("Number Of CPU: %d\n", numOfCpu)

	//setupLogger()

	redisConnection()

	// migrasi database dimatikan, diasumsikan data user sudah tersedia pada database.
	// database migration
	db, err = sql.Open("postgres", dbServerUrl)
	if err != nil {
		exitf("Db open error: %v\n", err)
	}
	driver, err = postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		_ = db.Close()
		exitf("Db postgres driven error: %v\n", err)
	}
	migration, err = migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		_ = db.Close()
		exitf("Unable to initiate migration: %v\n", err)
	}
	err = migration.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
	if err != nil {
		log.Printf("Migration error: %s", err.Error())
	}
	err = db.Close()
	if err != nil {
		log.Printf("Db close error: %s", err.Error())
	}
	// end database migration

	dbConnection()
	defer func() {
		pgxConn.Close()
		redisCl.Close()
		pgxConnMstData.Close()
	}()

	minioConn, err = utils.NewMinIOConn(
		os.Getenv("MINIO_ENDPOINT"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("MINIO_BUCKET_NAME"),
	)
	if err != nil {
		exitf("Unable to connect to object storage: %v\n", err)
	}

	serverReadTimeoutInt, err := strconv.Atoi(serverReadTimeout)
	if err != nil {
		exitf("Failed casting timeout context: ", err)
	}
	timeoutContext := time.Duration(serverReadTimeoutInt) * time.Second

	// Define a validator
	vld = utils.NewValidator()

	// jwt manager
	jwtMgr = utils.NewJWTManager(jwtSecertKey, serverName)

	// Define Fiber config.
	config := configs.FiberConfig()
	app := fiber.New(config)

	app.Static("/assets", "./assets")

	// Swagger handler
	app.Get("/swagger/*", swagger.HandlerDefault)

	middL := _deliveryMiddleware.InitMiddleware(app, redisCl, logger)
	//app.Use(middL.RateLimiter())
	app.Use(middL.CORS())
	app.Use(middL.LOGGER())
	//app.Use(middL.RateLimiter())
	app.Use(func(c *fiber.Ctx) error { // Middleware to check for whitelisted domains
		if alwOrg == "*" {
			// Continue to the next middleware/handler
			return c.Next()
		}

		// Use "X-Forwarded-Host" to simulate the Host header in Postman
		origin := c.Get("Origin")
		// log.Println("Origin: ", origin)

		alwOrgArr := strings.Split(alwOrg, ",")
		// log.Println("alwOrgArr: ", alwOrgArr)

		var originMatch bool
		for _, alo := range alwOrgArr {
			if origin == alo {
				originMatch = true
				break
			} else {
				/*host := c.Hostname()
				// log.Println("Host: ", host)
				if "https://"+host == alo || "http://"+host == alo {
					originMatch = true
					break
				}*/
			}
		}
		if !originMatch {
			// log.Println("not match")
			return c.Status(fiber.StatusForbidden).SendString("403 - AU: origin not allowed")
		}

		// Continue to the next middleware/handler
		return c.Next()
	})

	http.NewSiteHandler(app, controller.NewSiteController(pgxConn, timeoutContext), vld)

	// captchaStore := captcha_store.NewPostgreSQLStore(pgxConn)
	// http.NewCaptchaHandler(app, controller.NewCaptchaController(captchaStore, timeoutContext, vld))

	authController := controller.NewAuthController(
		pgxConn,
		pgxConnMstData,
		timeoutContext,
		jwtMgr,
		vld,
		redisCl,
	)
	http.NewAuthHandler(app, authController, vld)

	// private router
	rStrict := app.Group("/strict", middL.JWT()) // router for api private access

	// user_handler
	http.NewUserHandler(
		rStrict,
		vld,
		controller.NewUserController(
			pgxConn,
			timeoutContext,
			redisCl,
		),
	)

	// bumd_handler
	bumd.NewBumdHandler(
		rStrict,
		vld,
		controller_bumd.NewBumdController(pgxConn, pgxConnMstData, minioConn),
		pgxConn,
		minioConn,
	)

	// master
	// bentuk_badan_hukum
	http_master.NewBentukBadanHukumHandler(
		rStrict,
		vld,
		controller_mst.NewBentukBadanHukumController(pgxConn),
	)
	// bentuk_usaha
	http_master.NewBentukUsahaHandler(
		rStrict,
		vld,
		controller_mst.NewBentukUsahaController(pgxConn),
	)
	// jenis_dokumen
	http_master.NewJenisDokumenHandler(
		rStrict,
		vld,
		controller_mst.NewJenisDokumenController(pgxConn),
	)

	http_master.NewPendidikanHandler(
		rStrict,
		vld,
		controller_mst.NewPendidikanController(pgxConn),
	)

	http_util.StartServer(app)
}

func exitf(s string, args ...interface{}) {
	errorf(s, args...)
	os.Exit(1)
}

func errorf(s string, args ...interface{}) {
	log.Printf(s+"\n", args...)
}

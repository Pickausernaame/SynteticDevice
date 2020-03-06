package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
)

type App struct {
	Router  *gin.Engine
	Handler *Handler
}

type Agregator struct {
	Connection *pgx.ConnPool
}

type Handler struct {
	Agregator *Agregator
}

func (a *App) CreateHandler(conn *pgx.ConnPoolConfig) {
	Pool, _ := pgx.NewConnPool(*conn)
	var h = &Handler{
		Agregator: &Agregator{},
	}
	h.Agregator.Connection = Pool
	a.Handler = h
}

func CreateApp(conn *pgx.ConnPoolConfig) *App {
	var a App
	a.CreateHandler(conn)
	a.CreateRouter()
	return &a
}

func (a *App) CreateRouter() (router *gin.Engine) {
	a.Router = gin.New()
	a.Router.Use(gin.Logger())
	a.Router.Use(gin.Recovery())
	//a.Handler.Agregator.ClearTableAgr()

	a.Handler.Agregator.CreateTableAgr()

	// Создаем таблицы в бд
	api := a.Router.Group("/api")
	{
		api.POST("/nodered", a.Handler.InsertDataByNode)
		api.GET("/nodered", a.Handler.GetDataByNode)

		api.POST("/flogo", a.Handler.InsertDataByFlogo)
		api.GET("/flogo", a.Handler.GetDataByFlogo)

	}
	return
}

type Packet struct {
	AccX  float64 `json:"accX"`
	AccY  float64 `json:"accY"`
	AccZ  float64 `json:"accZ"`
	GyroX float64 `json:"gyroX"`
	GyroY float64 `json:"gyroY"`
	GyroZ float64 `json:"gyroZ"`
}

func (h *Handler) InsertDataByFlogo(c *gin.Context) {
	p := Packet{}
	err := json.NewDecoder(c.Request.Body).Decode(&p)
	if err != nil {
		fmt.Println(err)
	}
	err = h.Agregator.InsertData("FLOGO", p)
	if err != nil {
		c.Status(409)
		return
	}
	c.Status(201)
}

func (h *Handler) GetDataByFlogo(c *gin.Context) {
	p, err := h.Agregator.SelectData("FLOGO")
	if err != nil {
		c.Status(409)
		return
	}
	c.JSON(200, p)
}

func (h *Handler) InsertDataByNode(c *gin.Context) {
	p := Packet{}
	err := json.NewDecoder(c.Request.Body).Decode(&p)
	if err != nil {
		fmt.Println(err)
	}
	err = h.Agregator.InsertData("NODE", p)
	if err != nil {
		c.Status(409)
		return
	}
	c.Status(201)
}

func (h *Handler) GetDataByNode(c *gin.Context) {
	p, err := h.Agregator.SelectData("NODE")
	if err != nil {
		c.Status(409)
		return
	}
	c.JSON(200, p)
}

func (agr *Agregator) CreateTableAgr() {
	sql := `
	-- Уничтожаем существующие таблицы таблицы

	--DROP TABLE IF EXISTS flogo			CASCADE;
	--DROP TABLE IF EXISTS nodered			CASCADE;

	-- Создаем таблицы.
	-- Таблица Flogo.
	CREATE TABLE IF NOT EXISTS flogo (
		accX		double precision,
		accY		double precision,
		accZ		double precision,
		gyroX		double precision,
		gyroY		double precision,
		gyroZ		double precision			 );

	CREATE TABLE IF NOT EXISTS nodered (
		accX		double precision,
		accY		double precision,
		accZ		double precision,
		gyroX		double precision,
		gyroY		double precision,
		gyroZ		double precision			 );
`

	_, err := agr.Connection.Exec(sql)
	fmt.Println(err)
}

func (agr *Agregator) InsertData(name string, packet Packet) (err error) {
	sql := ""
	if name == "FLOGO" {
		sql = `INSERT INTO flogo (accX, accY, accZ, gyroX,  gyroY, gyroZ) VALUES ( $1, $2, $3, $4, $5, $6);`
	} else if name == "NODE" {
		sql = `INSERT INTO nodered (accX, accY, accZ, gyroX,  gyroY, gyroZ) VALUES ( $1, $2, $3, $4, $5, $6);`
	}
	_, err = agr.Connection.Exec(sql, packet.AccX, packet.AccY, packet.AccZ, packet.GyroX, packet.GyroY, packet.GyroZ)
	fmt.Println(err)
	return
}

func (agr *Agregator) SelectData(name string) (packet []Packet, err error) {
	sql := ""
	if name == "FLOGO" {
		sql = `SELECT accX, accY, accZ, gyroX,  gyroY, gyroZ FROM flogo;`
	} else if name == "NODE" {
		sql = `SELECT accX, accY, accZ, gyroX,  gyroY, gyroZ FROM nodered;`
	}

	rows, err := agr.Connection.Query(sql)
	for rows.Next() {
		p := Packet{}
		rows.Scan(&p.AccX, &p.AccY, &p.AccZ, &p.GyroX, &p.GyroY, &p.GyroZ)
		packet = append(packet, p)
	}
	defer rows.Close()
	fmt.Println(err)
	return
}

func main() {
	conf := pgx.ConnConfig{
		User:      "sayonara",
		Password:  "boy",
		Host:      "localhost",
		Port:      5432,
		Database:  "techno",
		TLSConfig: nil,
	}
	confPool := pgx.ConnPoolConfig{
		ConnConfig:     conf,
		MaxConnections: 8,
	}
	a := CreateApp(&confPool)
	a.CreateRouter()
	a.Router.Run(":5000")

}

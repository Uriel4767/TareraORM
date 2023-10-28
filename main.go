package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type Usuario struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
}

func main() {
	// Conectar a la base de datos (asegúrate de cambiar usuario y contraseña)
	var err error
	db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/finn?charset=utf8&parseTime=True&loc=Local&allowNativePasswords=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Migrar la estructura de Usuario a la base de datos
	db.AutoMigrate(&Usuario{})

	// Iniciar el enrutador de Gin
	router := gin.Default()

	// Definir rutas para el CRUD
	router.GET("/usuarios", ListarUsuarios)
	router.GET("/usuarios/:id", ObtenerUsuario)
	router.POST("/usuarios", CrearUsuario)
	router.PUT("/usuarios/:id", ActualizarUsuario)
	router.DELETE("/usuarios/:id", EliminarUsuario)

	// Servir archivos estáticos
	router.Static("/static", "./static")

	// Ruta para la página HTML
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Escuchar en el puerto 8080
	router.Run(":8080")
}

// Controladores para el CRUD
func ListarUsuarios(c *gin.Context) {
	var usuarios []Usuario
	db.Find(&usuarios)
	c.JSON(http.StatusOK, usuarios)
}

func ObtenerUsuario(c *gin.Context) {
	id := c.Param("id")
	var usuario Usuario
	if db.First(&usuario, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}
	c.JSON(http.StatusOK, usuario)
}

func CrearUsuario(c *gin.Context) {
	var usuario Usuario
	c.BindJSON(&usuario)
	db.Create(&usuario)
	c.JSON(http.StatusCreated, usuario)
}

func ActualizarUsuario(c *gin.Context) {
	id := c.Param("id")
	var usuario Usuario
	if db.First(&usuario, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}
	c.BindJSON(&usuario)
	db.Save(&usuario)
	c.JSON(http.StatusOK, usuario)
}

func EliminarUsuario(c *gin.Context) {
	id := c.Param("id")
	var usuario Usuario
	if db.First(&usuario, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}
	db.Delete(&usuario)
	c.JSON(http.StatusNoContent, gin.H{})
}

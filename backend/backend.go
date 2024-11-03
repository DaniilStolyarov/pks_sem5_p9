package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uint
	Email       string `gorm:"unique"`
	Password    string
	Name        string
	PhoneNumber string

	FavouriteItems []FavouriteItem `gorm:"foreignKey:UserID"`
	CartItems      []CartItem      `gorm:"foreignKey:UserID"`
}

type Service struct {
	gorm.Model
	ID          uint
	Name        string
	Category    string
	PriceRubles uint
	ImageHref   string
	Description string

	FavouriteItems []FavouriteItem `gorm:"foreignKey:ServiceID"`
	CartItems      []CartItem      `gorm:"foreignKey:ServiceID"`
}

type FavouriteItem struct {
	gorm.Model
	ID        uint
	UserID    uint
	ServiceID uint
	Service   Service `gorm:"foreignKey:ServiceID"`
}

type CartItem struct {
	gorm.Model
	ID        uint
	UserID    uint
	ServiceID uint
	Count     uint
	Service   Service `gorm:"foreignKey:ServiceID"`
}

var db *gorm.DB

func connectToDatabase() *gorm.DB {
	connectionString := "host=localhost user=postgres dbname=barbershop_flutter password=postgres sslmode=disable"
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func initDatabase() {
	db.Migrator().DropTable(&User{}, &Service{}, &FavouriteItem{}, &CartItem{})
	err := db.AutoMigrate(&User{}, &Service{}, &FavouriteItem{}, &CartItem{})
	if err != nil {
		panic("faled to init database")
	}
}

func fillDatabaseWithTestData() {
	daniilStolyarov := User{ID: 1, Email: "22T0318@gmail.com", Name: "Daniil Stolyarov", Password: "test_password", PhoneNumber: "89876543210"}
	db.Create(&daniilStolyarov)

	services := []Service{
		{Name: "Боб", Category: "Стрижка и укладка", PriceRubles: 799, ImageHref: "https://i.imgur.com/pLhAUHv.jpeg", Description: "Боб — это классическая и универсальная стрижка, которая подходит для любого типа волос и формы лица. Она может быть выполнена различной длины и формы, от короткого, доходящего до подбородка боба до длинного, доходящего до плеч боба. Боб идеально подходит для тех, кто хочет добавить объем и текстуру своим волосам, а также скрыть недостатки лица и подчеркнуть скулы."},
		{Name: "Цезарь", Category: "Стрижка и укладка", PriceRubles: 699, ImageHref: "https://i.imgur.com/G8eOfAE.jpeg", Description: "Цезарь — это короткая мужская стрижка с ровной челкой и короткими, одинаковой длины волосами по бокам и сзади. Она идеально подходит для мужчин с прямыми или волнистыми волосами и квадратной или круглой формой лица. Стрижка Цезарь — это стильный и универсальный вариант, который легко укладывать и поддерживать."},
		{Name: "Гарсон", Category: "Стрижка и укладка", PriceRubles: 599, ImageHref: "https://i.imgur.com/Atfpw5S.jpeg", Description: "Гарсон — это короткая, женственная стрижка, которая добавляет образу дерзости и стиля. Она характеризуется короткими, филированными волосами, часто с челкой или асимметрией. Гарсон идеально подходит для тех, кто хочет добавить объем и текстуру тонким волосам, а также скрыть недостатки лица и подчеркнуть скулы."},
		{Name: "Полубокс", Category: "Стрижка и укладка", PriceRubles: 499, ImageHref: "https://i.imgur.com/rpKjovf.jpeg", Description: "Полубокс — это короткая, практичная и универсальная мужская стрижка, которая подходит для любого типа волос и формы лица. Она характеризуется короткими волосами на висках и затылке, переходящими в более длинные волосы на макушке. Полубокс — это стильный и неприхотливый вариант, который легко укладывать и поддерживать."},
		{Name: "Шапочка", Category: "Стрижка и укладка", PriceRubles: 599, ImageHref: "https://i.imgur.com/MiLOdAO.jpeg", Description: "Шапочка — это женственная и элегантная стрижка, которая подходит для любого типа волос и формы лица. Она характеризуется короткими, округлыми волосами, которые плавно обрамляют лицо, создавая эффект шапочки. Шапочка — это универсальный вариант, который легко укладывать и поддерживать, придавая образу утонченность и шарм."},
		{Name: "Маллет", Category: "Стрижка и укладка", PriceRubles: 699, ImageHref: "https://i.imgur.com/oAlIgc2.jpeg", Description: "Маллет — это дерзкая и экстравагантная стрижка, которая характеризуется короткими волосами на макушке и длинными волосами на затылке. Маллет — это стильный и необычный вариант, который подходит для тех, кто хочет выделиться и подчеркнуть свой индивидуальный стиль. Она требует особого ухода и укладки, придавая образу неповторимость и притягательность."},
		{Name: "Милитари", Category: "Стрижка и укладка", PriceRubles: 799, ImageHref: "https://i.imgur.com/RHjZ63I.jpeg", Description: "Милитари — это мужественная и практичная стрижка, которая характеризуется короткими волосами по всей голове. Милитари — это универсальный и неприхотливый вариант, который подходит для любого типа волос и формы лица. Она не требует сложной укладки и идеально подходит для тех, кто ценит удобство и аккуратный внешний вид."},
		{Name: "Пикси", Category: "Стрижка и укладка", PriceRubles: 999, ImageHref: "https://i.imgur.com/01Z9Sga.jpeg", Description: "Пикси - это дерзкая и стильная стрижка, которая характеризуется короткими волосами на висках и затылке, переходящими в более длинные волосы на макушке. Пикси — это женственный и многогранный вариант, который подходит для любого типа волос и формы лица. Она позволяет создавать различные укладки, от гладких и элегантных до текстурных и небрежных."},
		{Name: "Короткое Каре", Category: "Стрижка и укладка", PriceRubles: 899, ImageHref: "https://i.imgur.com/HroeZo2.jpeg", Description: "Короткое Каре — это элегантная и стильная стрижка, которая характеризуется прямыми волосами, подстриженными до уровня подбородка или чуть выше. Короткое каре — это универсальный и практичный вариант, который подходит для любого типа волос и формы лица. Оно легко укладывается и позволяет создавать различные образы, от классических и сдержанных до современных и авангардных."},
		{Name: "Звезда", Category: "Стрижка и укладка", PriceRubles: 1199, ImageHref: "https://i.imgur.com/dOnl5Vz.jpeg", Description: "Звезда — это смелая и дерзкая стрижка, которая характеризуется выбритыми висками и затылком в форме звезды. Звезда — это авангардный и экстравагантный вариант, который подходит для тех, кто не боится выделяться из толпы. Она требует особого ухода и укладки, но придает образу неповторимый и притягательный вид."},
		{Name: "Теннис", Category: "Стрижка и укладка", PriceRubles: 599, ImageHref: "https://i.imgur.com/1uNHwbK.jpeg", Description: "Теннис — это спортивная и практичная стрижка, которая характеризуется короткими волосами по всей голове, подстриженными машинкой. Теннис — это универсальный и неприхотливый вариант, который подходит для любого типа волос и формы лица. Она не требует сложной укладки и идеально подходит для тех, кто ценит удобство и аккуратный внешний вид."},
	}
	for _, service := range services {
		db.Create(&service)
	}
}

func getAllServices(writer http.ResponseWriter, request *http.Request) {
	var services []Service
	db.Find(&services)
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(services)
}
func getFavourite(writer http.ResponseWriter, request *http.Request) {
	var favouriteItems []FavouriteItem
	db.Preload("Service").Find(&favouriteItems)

	services := []Service{}
	for _, item := range favouriteItems {
		services = append(services, item.Service)
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(services)
}
func getCart(writer http.ResponseWriter, request *http.Request) {
	CartItems := []CartItem{}
	db.Preload("Service").Find(&CartItems)
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(CartItems)
}
func getUserData(writer http.ResponseWriter, request *http.Request) {
	userId := request.URL.Query().Get("id")
	user := User{}
	db.First(&user, userId)
	user.Password = ""
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(user)
}
func getService(writer http.ResponseWriter, request *http.Request) {
	serviceId := request.URL.Query().Get("id")
	service := Service{}
	db.First(&service, serviceId)
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(service)
}

func addService(writer http.ResponseWriter, request *http.Request) {
	service := Service{}

	err := json.NewDecoder(request.Body).Decode(&service)
	if err != nil {
		fmt.Println(err)

	}
	db.Create(&service)
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(service)
}

func addFavourite(writer http.ResponseWriter, request *http.Request) {
	serviceId := request.URL.Query().Get("service_id")
	serviceId_parsed, _ := strconv.Atoi(serviceId)
	fmt.Println(serviceId_parsed)
	favouriteItem := FavouriteItem{ServiceID: uint(serviceId_parsed), UserID: 1, ID: 0}
	db.Create(&favouriteItem)
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(favouriteItem)
}

func addCart(writer http.ResponseWriter, request *http.Request) {
	serviceId := request.URL.Query().Get("service_id")
	serviceId_parsed, _ := strconv.Atoi(serviceId)
	CartItem := CartItem{ServiceID: uint(serviceId_parsed), UserID: 1, Count: 1}
	db.Create(&CartItem)
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(CartItem)
}

func updateCart(writer http.ResponseWriter, request *http.Request) {
	serviceId := request.URL.Query().Get("service_id")
	serviceId_parsed, _ := strconv.Atoi(serviceId)
	count := request.URL.Query().Get("count")
	count_parsed, _ := strconv.Atoi(count)

	cartItem := CartItem{}
	db.First(&cartItem, "service_id = ?", serviceId_parsed)
	cartItem.Count = uint(count_parsed)
	db.Model(&CartItem{}).Where("service_id = ?", serviceId_parsed).Updates(cartItem)

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(cartItem)
}

func updateUser(writer http.ResponseWriter, request *http.Request) {
	user := User{}
	json.NewDecoder(request.Body).Decode(&user)
	fmt.Println(user.Name + " " + strconv.Itoa(int(user.ID)) + " " + user.Password + " " + user.Email)
	result := db.Model(&User{}).Where("id = ?", int(user.ID)).Select("Name", "PhoneNumber", "Email").Updates(user)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(user)
}
func updateService(writer http.ResponseWriter, request *http.Request) {
	service := Service{}
	json.NewDecoder(request.Body).Decode(&service)
	result := db.Model(&Service{}).Where("id = ?", int(service.ID)).Updates(service)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(service)
}
func removeCartItem(writer http.ResponseWriter, request *http.Request) {
	cartItem := CartItem{}
	serviceId := request.URL.Query().Get("service_id")
	serviceId_parsed, _ := strconv.Atoi(serviceId)
	db.Where("service_id = ?", serviceId_parsed).Unscoped().Delete(&cartItem)
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(cartItem)
}
func removeFavouriteItem(writer http.ResponseWriter, request *http.Request) {
	favouriteItem := FavouriteItem{}
	serviceId := request.URL.Query().Get("service_id")
	serviceId_parsed, _ := strconv.Atoi(serviceId)
	fmt.Println(serviceId_parsed)
	db.Where("service_id = ?", serviceId_parsed).Unscoped().Delete(&favouriteItem)
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(favouriteItem)
}
func removeService(writer http.ResponseWriter, request *http.Request) {
	service := Service{}
	Id := request.URL.Query().Get("id")
	Id_parsed, _ := strconv.Atoi(Id)
	db.Where("id = ?", Id_parsed).Unscoped().Delete(&service)
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(service)
}
func main() {
	db = connectToDatabase()
	initDatabase()
	fillDatabaseWithTestData()

	router := mux.NewRouter()
	router.HandleFunc("/services", getAllServices).Methods("GET")
	router.HandleFunc("/user", getUserData).Methods("GET")
	router.HandleFunc("/favourite", getFavourite).Methods("GET")
	router.HandleFunc("/cart", getCart).Methods("GET")
	router.HandleFunc("/service", addService).Methods("POST")
	router.HandleFunc("/favourite", addFavourite).Methods("POST")
	router.HandleFunc("/cart", addCart).Methods("POST")
	router.HandleFunc("/cart", updateCart).Methods("PUT")
	router.HandleFunc("/user", updateUser).Methods("PUT")
	router.HandleFunc("/cart", removeCartItem).Methods("DELETE")
	router.HandleFunc("/favourite", removeFavouriteItem).Methods("DELETE")
	router.HandleFunc("/service", removeService).Methods("DELETE")
	router.HandleFunc("/service", updateService).Methods("PUT")
	router.HandleFunc("/service", getService).Methods("GET")

	fmt.Println("running...")
	http.ListenAndServe(":8080", router)
}

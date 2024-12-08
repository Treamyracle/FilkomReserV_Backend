package handler

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Data structures for the app
type NotificationWithRoom struct {
	ID          uint   `json:"id"`
	RoomName    string `json:"room_name"`
	BorrowDate  string `json:"borrow_date"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Status      string `json:"status"`
	File        string `json:"file"`
	Description string `json:"description"`
}

type Notification struct {
	ID          uint   `json:"id" form:"id"`           // ID notifikasi
	RoomID      uint   `json:"room_id" form:"room_id"` // ID ruangan
	BorrowDate  string `json:"borrow_date" form:"borrow_date"`
	StartTime   string `json:"start_time" form:"start_time"`
	EndTime     string `json:"end_time" form:"end_time"`
	Status      string `json:"status" form:"status"`
	File        string `json:"file" form:"file"`
	Description string `json:"description" form:"description"`
	Username    string `json:"username" form:"username"` // Menambahkan username
}

type Room struct {
	ID           uint       `json:"id"`
	Name         string     `json:"name"`
	Capacity     int        `json:"capacity"`
	Description  string     `json:"description"`
	ImageURL     string     `json:"image_url"`
	Availability []TimeSlot `json:"availability"`
	Facilities   []string   `json:"facilities"`
	UsageHistory []string   `json:"usage_history"`
}

type TimeSlot struct {
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// In-memory storage for rooms and notifications
var rooms = []Room{
	{ID: 1, Name: "GKM", Capacity: 100, Description: "Gedung Kreativitas Bersama FILKOM UB", ImageURL: "https://cdn.pixabay.com/photo/2017/03/28/12/17/chairs-2181994_1280.jpg",
		Facilities:   []string{"Videotron", "Ac Central", "Kursi", "Lampu Sorot", "Audio", "Panggung"},
		UsageHistory: []string{"Schotival 2024", "Hology 7.0"},
		Availability: []TimeSlot{
			{Date: "2024-10-12", StartTime: "14:00", EndTime: "16:00"},
			{Date: "2024-10-13", StartTime: "09:00", EndTime: "11:00"},
		}},
	{ID: 2, Name: "Algoritma G2", Capacity: 200, Description: "Gedung Algoritma G2 FILKOM UB", ImageURL: "https://cdn.pixabay.com/photo/2017/03/28/12/17/chairs-2181994_1280.jpg",
		Facilities:   []string{"Videotron", "Ac Central", "Kursi", "Lampu Sorot", "Audio"},
		UsageHistory: []string{"TechTalk 2024", "Futuristic Innovation Summit"},
		Availability: []TimeSlot{
			{Date: "2024-10-12", StartTime: "08:00", EndTime: "10:00"},
			{Date: "2024-10-15", StartTime: "10:00", EndTime: "12:00"},
		}},
}

var notifications = []Notification{
	{
		ID:          1,
		RoomID:      1,
		BorrowDate:  "2024-10-12",
		StartTime:   "14:00",
		EndTime:     "16:00",
		Status:      "diterima",
		File:        "contract1.pdf",
		Description: "Peminjaman untuk seminar teknologi",
	},
	{
		ID:          2,
		RoomID:      2,
		BorrowDate:  "2024-10-13",
		StartTime:   "09:00",
		EndTime:     "11:00",
		Status:      "proses",
		File:        "",
		Description: "Penggunaan ruangan untuk pertemuan organisasi",
	},
}

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))

var users = map[string]string{
	"admin": "password123",
	"user":  "password456",
}

// The main handler for Vercel (serverless)
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.Default()

	// Use CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500"},                   // Origin yang diizinkan
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // Metode yang diizinkan
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Header yang diizinkan
		AllowCredentials: true,                                                // Jika kamu membutuhkan cookie atau otentikasi
	}))

	// Routes
	auth := router.Group("/")
	auth.Use(JWTAuthMiddleware()) // Gunakan middleware JWT setelah CORS
	auth.GET("/rooms", getRooms)
	auth.GET("/rooms/:id", getRoomByID)
	auth.POST("/rooms", addRoom)
	auth.POST("/rooms/:id/availability", addAvailability)
	auth.PUT("/rooms/:id/availability", updateRoomAvailability)
	auth.POST("/notifications", addNotification)
	auth.GET("/notifications", getNotifications)
	auth.GET("/notifications-with-name", GetNotificationsWithRoomName)

	router.POST("/logout", LogoutHandler)

	// Run Gin router to handle the request
	router.ServeHTTP(w, r)
}

// LogoutHandler handles logout requests and removes the JWT cookie
func LogoutHandler(c *gin.Context) {
	// Remove the Authorization cookie by setting it with an expired time
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// JWTAuthMiddleware authenticates requests using the JWT in cookies
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication token required"})
			c.Abort()
			return
		}

		// Remove the "Bearer " prefix if it exists
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store username in context
		claims, _ := token.Claims.(jwt.MapClaims)
		c.Set("username", claims["username"])

		c.Next()
	}
}

// GenerateJWT generates a JWT token
func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

// LoginHandler handles login requests
// LoginHandler handles login requests and sets JWT in a cookie
// LoginHandler handles login requests and sets JWT in a cookie
func LoginHandler(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate credentials
	if password, ok := users[input.Username]; !ok || password != input.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT
	token, err := GenerateJWT(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Set JWT as a cookie
	c.SetCookie("Authorization", token, 3600, "/", "", true, false) // SameSite=None, Secure=false (untuk localhost)

	// Return response with username and user_id (assuming user_id is their username for now)
	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"username": input.Username,
		"token":    token,
	})
}

// Handler functions (such as getRooms, getRoomByID, etc.) go here...

func getRooms(c *gin.Context) {
	c.JSON(http.StatusOK, rooms)
}

func getRoomByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}
	for _, room := range rooms {
		if room.ID == uint(id) {
			c.JSON(http.StatusOK, room)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
}

func addRoom(c *gin.Context) {
	var newRoom Room
	if err := c.ShouldBindJSON(&newRoom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add the new room to the in-memory storage
	newRoom.ID = uint(len(rooms) + 1)
	rooms = append(rooms, newRoom)

	c.JSON(http.StatusOK, newRoom)
}

func addAvailability(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	var newSlot TimeSlot
	if err := c.ShouldBindJSON(&newSlot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the room by ID and add the availability slot
	for i, room := range rooms {
		if room.ID == uint(id) {
			rooms[i].Availability = append(rooms[i].Availability, newSlot)
			c.JSON(http.StatusOK, gin.H{"message": "Availability added successfully"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
}

func updateRoomAvailability(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	date := c.DefaultQuery("date", "")
	startTime := c.DefaultQuery("start_time", "")
	endTime := c.DefaultQuery("end_time", "")

	// Find the room by ID and update availability
	for i, room := range rooms {
		if room.ID == uint(id) {
			// Find and toggle the availability based on date and time
			for j, slot := range room.Availability {
				if slot.Date == date && slot.StartTime == startTime && slot.EndTime == endTime {
					rooms[i].Availability = append(rooms[i].Availability[:j], rooms[i].Availability[j+1:]...)
					c.JSON(http.StatusOK, gin.H{"message": "Room availability updated successfully"})
					return
				}
			}

			// If the time slot wasn't found, add it
			rooms[i].Availability = append(rooms[i].Availability, TimeSlot{Date: date, StartTime: startTime, EndTime: endTime})
			c.JSON(http.StatusOK, gin.H{"message": "Room availability added successfully"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
}

func addNotification(c *gin.Context) {
	var newNotification Notification
	// Bind form data ke struct Notification
	if err := c.ShouldBind(&newNotification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind form data: " + err.Error()})
		return
	}

	// Pastikan field yang diperlukan ada
	if newNotification.RoomID == 0 || newNotification.BorrowDate == "" || newNotification.StartTime == "" || newNotification.EndTime == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field room_id, borrow_date, start_time, and end_time are required"})
		return
	}

	// Ambil username dari konteks
	username, _ := c.Get("username")
	if username == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Set username pada notifikasi
	newNotification.Username = username.(string)

	// Proses file jika ada
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		newNotification.File = ""
	} else {
		defer file.Close()
		// Simpan file dan tentukan path file
		filePath := "path_to_saved_file"
		newNotification.File = filePath
	}

	// Tentukan status default jika tidak ada
	if newNotification.Status == "" {
		newNotification.Status = "proses"
	}

	// Tambahkan ID pada notifikasi baru
	newNotification.ID = uint(len(notifications) + 1)
	// Simpan notifikasi
	notifications = append(notifications, newNotification)

	// Kembalikan respons
	c.JSON(http.StatusOK, gin.H{"message": "Notification added successfully", "notification": newNotification})
}

func getNotifications(c *gin.Context) {
	c.JSON(http.StatusOK, notifications)
}

func GetNotificationsWithRoomName(c *gin.Context) {
	var result []NotificationWithRoom

	// Loop through each notification and find corresponding room details
	for _, notif := range notifications {
		for _, room := range rooms {
			if room.ID == notif.RoomID {
				// Map the notification to include room name and other details
				result = append(result, NotificationWithRoom{
					ID:          notif.ID,
					RoomName:    room.Name,
					BorrowDate:  notif.BorrowDate,
					StartTime:   notif.StartTime,
					EndTime:     notif.EndTime,
					Status:      notif.Status,
					File:        notif.File,
					Description: notif.Description,
				})
			}
		}
	}

	// Return the result with room names included
	c.JSON(http.StatusOK, result)
}

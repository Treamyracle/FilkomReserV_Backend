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
	Role     string `json:"role"`
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

type UpdateStatusRequest struct {
	ID     uint   `json:"id"`
	Status string `json:"status"`
}

// In-memory storage for rooms and notifications
var rooms = []Room{
	{
		ID:           1,
		Name:         "GKM",
		Capacity:     100,
		Description:  "Gedung Kreativitas Bersama FILKOM UB",
		ImageURL:     "https://images.unsplash.com/photo-1498050108023-c5249f4df085",
		Facilities:   []string{"Videotron", "Ac Central", "Kursi", "Lampu Sorot", "Audio", "Panggung"},
		UsageHistory: []string{"Schotival 2024", "Hology 7.0"},
		Availability: []TimeSlot{
			{Date: "2024-10-12", StartTime: "08:00", EndTime: "10:00"},
			{Date: "2024-10-12", StartTime: "10:30", EndTime: "12:30"},
			{Date: "2024-10-12", StartTime: "13:00", EndTime: "15:00"},
			{Date: "2024-10-12", StartTime: "15:30", EndTime: "18:00"},
			{Date: "2024-10-12", StartTime: "18:30", EndTime: "21:00"},
		},
	},
	{
		ID:           2,
		Name:         "Algoritma G2",
		Capacity:     200,
		Description:  "Gedung Algoritma G2 FILKOM UB",
		ImageURL:     "https://images.unsplash.com/photo-1524758631624-e2822e304c36",
		Facilities:   []string{"Videotron", "Ac Central", "Kursi", "Lampu Sorot", "Audio"},
		UsageHistory: []string{"TechTalk 2024", "Futuristic Innovation Summit"},
		Availability: []TimeSlot{
			{Date: "2024-10-12", StartTime: "08:00", EndTime: "11:00"},
			{Date: "2024-10-12", StartTime: "11:30", EndTime: "13:30"},
			{Date: "2024-10-12", StartTime: "14:00", EndTime: "16:00"},
			{Date: "2024-10-12", StartTime: "16:30", EndTime: "19:30"},
			{Date: "2024-10-12", StartTime: "20:00", EndTime: "23:00"},
		},
	},
	{
		ID:           3,
		Name:         "Lab AI",
		Capacity:     50,
		Description:  "Laboratorium Artificial Intelligence FILKOM UB",
		ImageURL:     "https://images.unsplash.com/photo-1571501679680-de32f1e7aad4",
		Facilities:   []string{"Komputer", "Proyektor", "AC", "Papan Tulis"},
		UsageHistory: []string{"Workshop AI 2024", "Hackathon UB"},
		Availability: []TimeSlot{
			{Date: "2024-10-12", StartTime: "09:00", EndTime: "11:00"},
			{Date: "2024-10-12", StartTime: "11:30", EndTime: "13:30"},
			{Date: "2024-10-12", StartTime: "14:00", EndTime: "16:00"},
			{Date: "2024-10-12", StartTime: "16:30", EndTime: "18:30"},
		},
	},
	{
		ID:           4,
		Name:         "Ruang Seminar 1",
		Capacity:     150,
		Description:  "Ruang Seminar Utama FILKOM UB",
		ImageURL:     "https://images.unsplash.com/photo-1498050108023-c5249f4df085",
		Facilities:   []string{"Mikrofon", "Proyektor", "Kursi", "Panggung"},
		UsageHistory: []string{"Seminar Nasional IT", "Konferensi Data Science"},
		Availability: []TimeSlot{
			{Date: "2024-10-12", StartTime: "08:30", EndTime: "10:30"},
			{Date: "2024-10-12", StartTime: "11:00", EndTime: "13:00"},
			{Date: "2024-10-12", StartTime: "13:30", EndTime: "15:30"},
			{Date: "2024-10-12", StartTime: "16:00", EndTime: "18:00"},
		},
	},
	{
		ID:           5,
		Name:         "Ruang Seminar 2",
		Capacity:     100,
		Description:  "Ruang Seminar Kedua FILKOM UB",
		ImageURL:     "https://images.unsplash.com/photo-1504384308090-c894fdcc538d",
		Facilities:   []string{"Proyektor", "Kursi", "Lampu Sorot", "Mikrofon"},
		UsageHistory: []string{"Workshop Desain UI", "Seminar Open Source"},
		Availability: []TimeSlot{
			{Date: "2024-10-12", StartTime: "09:00", EndTime: "11:00"},
			{Date: "2024-10-12", StartTime: "11:30", EndTime: "13:30"},
			{Date: "2024-10-12", StartTime: "14:00", EndTime: "16:00"},
			{Date: "2024-10-12", StartTime: "16:30", EndTime: "18:30"},
		},
	},
	{
		ID:           6,
		Name:         "Ruang Rapat 1",
		Capacity:     20,
		Description:  "Ruang Rapat Pertama FILKOM UB",
		ImageURL:     "https://images.unsplash.com/photo-1571501679680-de32f1e7aad4",
		Facilities:   []string{"Meja Besar", "AC", "Kursi", "Papan Tulis"},
		UsageHistory: []string{"Rapat Prodi", "Rapat Dosen"},
		Availability: []TimeSlot{
			{Date: "2024-10-12", StartTime: "08:30", EndTime: "10:30"},
			{Date: "2024-10-12", StartTime: "11:00", EndTime: "13:00"},
			{Date: "2024-10-12", StartTime: "13:30", EndTime: "15:30"},
			{Date: "2024-10-12", StartTime: "16:00", EndTime: "18:00"},
		},
	},
	{
		ID:           7,
		Name:         "Ruang Diskusi 1",
		Capacity:     30,
		Description:  "Ruang Diskusi Utama FILKOM UB",
		ImageURL:     "https://images.unsplash.com/photo-1524758631624-e2822e304c36",
		Facilities:   []string{"Meja Bundar", "Papan Tulis", "AC", "Kursi"},
		UsageHistory: []string{"Diskusi Kelompok", "Pelatihan Teamwork"},
		Availability: []TimeSlot{
			{Date: "2024-10-12", StartTime: "09:00", EndTime: "11:00"},
			{Date: "2024-10-12", StartTime: "11:30", EndTime: "13:30"},
			{Date: "2024-10-12", StartTime: "14:00", EndTime: "16:00"},
			{Date: "2024-10-12", StartTime: "16:30", EndTime: "18:30"},
		},
	},
	{
		ID:           8,
		Name:         "Ruang Diskusi 2",
		Capacity:     25,
		Description:  "Ruang Diskusi Kedua FILKOM UB",
		ImageURL:     "https://images.unsplash.com/photo-1524758631624-e2822e304c36",
		Facilities:   []string{"Meja Bundar", "Proyektor", "Kursi", "Papan Tulis"},
		UsageHistory: []string{"Diskusi Topik Riset", "Sharing Session Alumni"},
		Availability: []TimeSlot{
			{Date: "2024-10-12", StartTime: "09:00", EndTime: "11:00"},
			{Date: "2024-10-12", StartTime: "11:30", EndTime: "13:30"},
			{Date: "2024-10-12", StartTime: "14:00", EndTime: "16:00"},
			{Date: "2024-10-12", StartTime: "16:30", EndTime: "18:30"},
		},
	},
	{
		ID:           9,
		Name:         "Lab Multimedia",
		Capacity:     40,
		Description:  "Laboratorium Multimedia FILKOM UB",
		ImageURL:     "https://images.unsplash.com/photo-1524758631624-e2822e304c36",
		Facilities:   []string{"Komputer", "Proyektor", "AC", "Speaker"},
		UsageHistory: []string{"Workshop Video Editing", "Pelatihan Photoshop"},
		Availability: []TimeSlot{
			{Date: "2024-10-12", StartTime: "08:00", EndTime: "10:00"},
			{Date: "2024-10-12", StartTime: "10:30", EndTime: "12:30"},
			{Date: "2024-10-12", StartTime: "13:00", EndTime: "15:00"},
			{Date: "2024-10-12", StartTime: "15:30", EndTime: "17:30"},
		},
	},
	{
		ID:           10,
		Name:         "Auditorium",
		Capacity:     300,
		Description:  "Auditorium Utama FILKOM UB",
		ImageURL:     "https://images.unsplash.com/photo-1524758631624-e2822e304c36",
		Facilities:   []string{"Videotron", "AC Central", "Panggung", "Lampu Sorot", "Kursi"},
		UsageHistory: []string{"Konser Kampus", "Wisuda UB"},
		Availability: []TimeSlot{
			{Date: "2024-10-12", StartTime: "08:00", EndTime: "10:00"},
			{Date: "2024-10-12", StartTime: "10:30", EndTime: "12:30"},
			{Date: "2024-10-12", StartTime: "13:00", EndTime: "15:00"},
			{Date: "2024-10-12", StartTime: "15:30", EndTime: "17:30"},
		},
	},
}

var notifications = []Notification{}

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))

var users = []LoginInput{
	{Username: "admin", Password: "password123", Role: "admin"},
	{Username: "user", Password: "password123", Role: "user"},
}

// The main handler for Vercel (serverless)
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.Default()

	// Use CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://projek-akhir-aps-fix.vercel.app", "https://feaps-treamyracles-projects.vercel.app", "https://feaps.vercel.app"}, // Origin yang diizinkan
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                                                                               // Metode yang diizinkan
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},                                                                               // Header yang diizinkan
		AllowCredentials: true,                                                                                                                              // Jika kamu membutuhkan cookie atau otentikasi
	}))

	// Routes
	router.POST("/login", LoginHandler)
	router.GET("/users", GetUsersHandler)
	router.POST("/logout", LogoutHandler)

	// Group routes yang membutuhkan autentikasi
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
	auth.POST("/update-status", updateStatusHandler)

	// Run Gin router to handle the request
	router.ServeHTTP(w, r)
}

func updateStatusHandler(c *gin.Context) {
	var req UpdateStatusRequest

	// Bind input JSON ke struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Validasi status baru
	if req.Status != "Diterima" && req.Status != "Ditolak" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status harus 'Diterima' atau 'Ditolak'"})
		return
	}

	// Cari notifikasi berdasarkan ID
	for i, notification := range notifications {
		if notification.ID == req.ID {
			// Perbarui status
			notifications[i].Status = req.Status

			// Kirim respons sukses
			c.JSON(http.StatusOK, gin.H{
				"message": "Status berhasil diperbarui",
			})
			return
		}
	}

	// Jika ID tidak ditemukan
	c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
}

func GetUsersHandler(c *gin.Context) {
	// Menampilkan semua data pengguna
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
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
func LoginHandler(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari user dalam slice users
	var foundUser *LoginInput
	for _, user := range users {
		if user.Username == input.Username {
			foundUser = &user
			break
		}
	}

	// Jika user tidak ditemukan atau password tidak cocok
	if foundUser == nil || foundUser.Password != input.Password {
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

	// Return response with username and role
	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"username": input.Username,
		"role":     foundUser.Role,
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
		newNotification.File = "sementara.pdf"
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

package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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
	ID          uint   `json:"id" form:"id"`
	RoomID      uint   `json:"room_id" form:"room_id"`
	BorrowDate  string `json:"borrow_date" form:"borrow_date"`
	StartTime   string `json:"start_time" form:"start_time"`
	EndTime     string `json:"end_time" form:"end_time"`
	Status      string `json:"status" form:"status"`
	File        string `json:"file" form:"file"`
	Description string `json:"description" form:"description"`
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

// The main handler for Vercel (serverless)
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.Default()

	// Use CORS middleware
	router.Use(cors.Default())

	// Routes
	router.GET("/rooms", getRooms)
	router.GET("/rooms/:id", getRoomByID)
	router.POST("/rooms", addRoom)
	router.POST("/rooms/:id/availability", addAvailability)
	router.PUT("/rooms/:id/availability", updateRoomAvailability)
	router.POST("/notifications", addNotification)
	router.GET("/notifications", getNotifications)
	router.GET("/notifications-with-name", GetNotificationsWithRoomName)

	// Run Gin router to handle the request
	router.ServeHTTP(w, r)
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
	if err := c.ShouldBind(&newNotification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind form data: " + err.Error()})
		return
	}

	if newNotification.RoomID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field room_id, borrow_date, start_time, and end_time are required"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		newNotification.File = ""
	} else {
		defer file.Close()
		filePath := "path_to_saved_file"
		newNotification.File = filePath
	}

	if newNotification.Status == "" {
		newNotification.Status = "proses"
	}

	newNotification.ID = uint(len(notifications) + 1)
	notifications = append(notifications, newNotification)

	c.JSON(http.StatusOK, gin.H{"message": "Notification added successfully", "notification": newNotification})
}

func getNotifications(c *gin.Context) {
	c.JSON(http.StatusOK, notifications)
}

func GetNotificationsWithRoomName(c *gin.Context) {
	var result []NotificationWithRoom

	for _, notif := range notifications {
		for _, room := range rooms {
			if room.ID == notif.RoomID {
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

	c.JSON(http.StatusOK, result)
}

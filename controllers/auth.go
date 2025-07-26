package controllers

import (
	"context"
	"fmt"
	"inventory-backend/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

func SetUserCollection(db *mongo.Database) {
	userCollection = db.Collection("users")
}

// JWT Secret Key - dalam production gunakan environment variable
var jwtSecret = []byte("your-secret-key-change-this-in-production")

func init() {
	// Check if JWT_SECRET environment variable exists
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		jwtSecret = []byte(secret)
	}
}

// Register godoc
// @Summary Register user baru
// @Description Mendaftarkan user baru dengan username, email, password, dan role
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "Data user baru"
// @Success 201 {object} map[string]interface{} "User berhasil didaftarkan"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/register [post]
func Register(c *fiber.Ctx) error {
	// Get validated data from middleware
	userData := c.Locals("userData").(models.RegisterRequest)

	// Check if user already exists
	var existingUser models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": userData.Email}).Decode(&existingUser)
	if err == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Email sudah terdaftar",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengenkripsi password",
		})
	}

	// Create new user
	newUser := models.User{
		ID:        primitive.NewObjectID(),
		Username:  userData.Username,
		Email:     userData.Email,
		Password:  string(hashedPassword),
		Role:      userData.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert to database
	_, err = userCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal menyimpan user ke database",
		})
	}

	// Generate JWT token
	token, err := generateJWTToken(newUser.ID, newUser.Email, newUser.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal generate token",
		})
	}

	// Prepare response
	response := models.LoginResponse{
		Token: token,
	}
	response.User.ID = newUser.ID
	response.User.Username = newUser.Username
	response.User.Email = newUser.Email
	response.User.Role = newUser.Role

	return c.Status(201).JSON(fiber.Map{
		"message": "User berhasil didaftarkan",
		"data":    response,
	})
}

// Login godoc
// @Summary Login user
// @Description Login dengan email dan password, return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Data login"
// @Success 200 {object} map[string]interface{} "Login berhasil"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	// Get validated data from middleware
	loginData := c.Locals("loginData").(models.LoginRequest)

	// Find user by email
	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": loginData.Email}).Decode(&user)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Email atau password salah",
		})
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Email atau password salah",
		})
	}

	// Generate JWT token
	token, err := generateJWTToken(user.ID, user.Email, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal generate token",
		})
	}

	// Prepare response
	response := models.LoginResponse{
		Token: token,
	}
	response.User.ID = user.ID
	response.User.Username = user.Username
	response.User.Email = user.Email
	response.User.Role = user.Role

	return c.JSON(fiber.Map{
		"message": "Login berhasil",
		"data":    response,
	})
}


func generateJWTToken(userID primitive.ObjectID, email, role string) (string, error) {
	// Ambil lokasi WIB (Asia/Jakarta)
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// fallback ke UTC jika gagal ambil lokasi
		loc = time.UTC
	}

	now := time.Now()
	expirationTime := now.Add(24 * time.Hour)

	// Logging waktu saat token dibuat dan kadaluarsa (dalam WIB)
	fmt.Println("🔐 JWT Token Generated:")
	fmt.Println("   Sekarang :", now.In(loc))
	fmt.Println("   Expired  :", expirationTime.In(loc))

	claims := jwt.MapClaims{
		"user_id": userID.Hex(),
		"email":   email,
		"role":    role,
		"exp":     expirationTime.Unix(),
		"iat":     now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}


// GetProfile godoc
// @Summary Get user profile
// @Description Mendapatkan data profile user yang sedang login
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Profile berhasil diambil"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User tidak ditemukan"
// @Router /auth/profile [get]
// Get current user profile
func GetProfile(c *fiber.Ctx) error {
	// Get user ID from JWT middleware
	userID := c.Locals("user_id").(string)
	
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var user models.User
	err = userCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User tidak ditemukan",
		})
	}

	// Don't return password
	user.Password = ""

	return c.JSON(fiber.Map{
		"message": "Profile berhasil diambil",
		"data":    user,
	})
}

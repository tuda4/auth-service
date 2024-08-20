package util

import (
	"math/rand"
	"time"
)

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	// Define the characters that can be used in the random string
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Create a byte slice with the specified length
	result := make([]byte, length)

	// Generate a random character from the charset for each position in the byte slice
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}

	// Convert the byte slice to a string and return it
	return string(result)
}

func CreateRandomEmail() string {
	username := RandomString(10)       // Generate a random username with length 10
	domain := "example.mb-app-001.com" // Replace with your desired domain name

	return username + "@" + domain
}

func CreateRandomPhone() string {
	// Define the characters that can be used in the random phone number
	charset := "0123456789"

	// Create a byte slice with the specified length for the phone number
	phoneNumberLength := 10
	phoneNumber := make([]byte, phoneNumberLength)

	// Generate a random digit from the charset for each position in the byte slice
	for i := 0; i < phoneNumberLength; i++ {
		phoneNumber[i] = charset[rand.Intn(len(charset))]
	}

	// Convert the byte slice to a string and format it as a phone number
	formattedPhoneNumber := string(phoneNumber[:3]) + "-" + string(phoneNumber[3:6]) + "-" + string(phoneNumber[6:])

	return formattedPhoneNumber
}

func CreateRandomBirthday() time.Time {
	rand.Seed(time.Now().UnixNano())

	// Generate a random year between 1900 and 2022
	year := rand.Intn(2022-1900+1) + 1900

	// Generate a random month between 1 and 12
	month := rand.Intn(12) + 1

	// Generate a random day between 1 and 28 (assuming all months have 28 days for simplicity)
	day := rand.Intn(28) + 1

	// Create a time.Time value with the generated year, month, and day
	birthday := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)

	return birthday
}

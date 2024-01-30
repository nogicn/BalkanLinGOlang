package middleware

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/haguro/elevenlabs-go"
)

func GenerateSpeech(word string, filename string) error {
	// create context
	ctx := context.Background()
	// create client
	client := elevenlabs.NewClient(ctx, os.Getenv("ELEVENLABS_API_KEY"), 30*time.Second)

	// create request
	ttsReq := elevenlabs.TextToSpeechRequest{
		Text:    word,
		ModelID: "eleven_multilingual_v2",
	}
	audio, err := client.TextToSpeech("pNInz6obpgDQGcFmaJgB", ttsReq)
	if err != nil {
		return fiber.NewError(500, "Greška pri generisanju izgovora!")
	}

	// Write the audio file bytes to disk
	if err := os.WriteFile("./public/pronunciation/"+filename, audio, 0644); err != nil {
		return fiber.NewError(500, "Greška pri generisanju izgovora!")
	}

	log.Println("Successfully generated audio file")
	return nil
}

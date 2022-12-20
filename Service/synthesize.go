package Service

import (
	"bytes"
	"cloud.google.com/go/storage"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"context"
	"donations_alert/model"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"strconv"
	"time"
)

func Synthesize(ctx context.Context, donation model.Donations) ([]byte, string) {
	// Instantiates a client.

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.

	newAmount := strconv.Itoa(donation.Amount)

	msg := newAmount + " dari " + donation.From + ", pesan: " + donation.Message

	if donation.Amount >= 200000 {
		msg = "WAW, " + newAmount + " dari " + donation.From + ", pesan: " + donation.Message
	}
	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: msg},
		},
		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "id-ID",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_FEMALE,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}

	object, err := uploadFile(resp.AudioContent)
	if err != nil {
		log.Fatalln(err)
	}

	// The resp's AudioContent is binary.
	//filename := "output.mp3"
	//err = ioutil.WriteFile(filename, resp.AudioContent, 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Audio content written to file: %v\n", filename)

	return resp.AudioContent, object
}

func uploadFile(content []byte) (string, error) {
	bucket := "donation_alert"
	name, err := uuid.NewUUID()
	if err != nil {
		log.Fatalln(err)
	}

	object := name.String() + ".mp3"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	o = o.If(storage.Conditions{DoesNotExist: true})
	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	// attrs, err := o.Attrs(ctx)
	// if err != nil {
	//      return fmt.Errorf("object.Attrs: %v", err)
	// }
	// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// Upload an object with storage.Writer.

	audioBuffer := bytes.NewBuffer(content)

	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, audioBuffer); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Println("berhasil upload")

	acl := client.Bucket(bucket).Object(object).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("ACLHandle.Set: %v", err)
	}
	//fmt.Fprintf(w, "Blob %v uploaded.\n", object)
	return object, nil
}

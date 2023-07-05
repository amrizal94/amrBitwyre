package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

// put keys in backtick (“) to avoid errors caused by spaces or tabs
const pubkey = `-----BEGIN PGP PUBLIC KEY BLOCK-----
Comment: https://gopenpgp.org
Version: GopenPGP 2.6.1

xjMEYMx6PRYJKwYBBAHaRw8BAQdAr1NwuCF2Qax/6PmpckXgnv8e+GvbtUe32qht
uMYgMUrNJ2hpcmluZ0BiaXR3eXJlLmNvbSA8aGlyaW5nQGJpdHd5cmUuY29tPsKP
BBAWCgAgBQJgzHo9BgsJBwgDAgQVCAoCBBYCAQACGQECGwMCHgEAIQkQaIECdlD0
TWkWIQTpZNu8ZwqbxC7rjhlogQJ2UPRNaeYiAP0Uzj2/dieBZe57W2ys6//kBG/n
XMakUpzFjkMrovLXewD/QP4RYrtgcn+R2UzsmDuiBMYXs4eUaQmeZX6IjQBF9APO
OARgzHo9EgorBgEEAZdVAQUBAQdAxnBq/ovBoJxk787vQ8HrJGlAWYEJALKeZeQv
to7DnA4DAQgHwngEGBYIAAkFAmDMej0CGwwAIQkQaIECdlD0TWkWIQTpZNu8Zwqb
xC7rjhlogQJ2UPRNac7/AQCuJ9t6PmrMNmFVPVLFlaMSqvY8zLiDN1iR4ahYGEKA
8gD/RWgwCidsekCbPkkcV5xe4B9u3nLgpqcVuP2+S3w30gU=
=0Ejb
-----END PGP PUBLIC KEY BLOCK-----`

// const privkey = `-----BEGIN PGP SIGNATURE-----
// Version: ProtonMail

// wnUEARYIACcFAmSk7XQJkGiBAnZQ9E1pFiEE6WTbvGcKm8Qu644ZaIECdlD0
// TWkAAOeHAQCQH7Fs4/lyjFwOzppHHB0TOfSG50C9GjlJ9WC8hz3v2QD+Ku2B
// u10OZBeTjiEagBFKlymSx8U/kd7QBamBIeDE5gU=
// =HNir
// -----END PGP SIGNATURE-----` // encrypted private key

func main() {

	// Configure Kafka writer
	kafkaTopic := "message"
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:50234"},
		Topic:   kafkaTopic,
	})

	// Get the message from request body
	router := gin.Default()
	router.POST("/messages", func(ctx *gin.Context) {
		message, err := ctx.GetRawData()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		// encrypt binary message using public key
		armor, err := helper.EncryptBinaryMessageArmored(pubkey, message)
		if err != nil {
			log.Fatal(err)
		}

		// Create a Kafka message
		kafkaMsg := kafka.Message{
			Value: []byte(armor),
		}

		// Publish the Kafka message
		err = kafkaWriter.WriteMessages(context.Background(), kafkaMsg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// Close the Kafka writer
		err = kafkaWriter.Close()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		msg := "Message published to Redpanda"
		log.Println(msg)
		ctx.JSON(http.StatusOK, gin.H{"message": msg})

	})

	err := router.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

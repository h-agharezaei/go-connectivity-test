package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

func main() {
	log.Println("🔍 Starting connectivity tests...")

	testPostgres()
	testRedis()
	testMongo()
	testNATS()
	testKafka()

	log.Println("✅ All tests finished")
}

// ---------- Postgres ----------
func testPostgres() {
	log.Print("Postgres: ")
	conn, err := pgx.Connect(ctx, "postgres://user:password@localhost:54321/dbname")
	if err != nil {
		log.Println("❌ FAILED:", err)
		return
	}
	defer conn.Close(ctx)

	if err := conn.Ping(ctx); err != nil {
		log.Println("❌ FAILED:", err)
		return
	}
	log.Println("✅ OK")
}

// ---------- Redis ----------
func testRedis() {
	log.Print("Redis: ")
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6350",
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Println("❌ FAILED:", err)
		return
	}
	log.Println("✅ OK")
}

// ---------- Mongo ----------
func testMongo() {
	log.Print("Mongo: ")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println("❌ FAILED:", err)
		return
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := client.Ping(ctxTimeout, nil); err != nil {
		log.Println("❌ FAILED:", err)
		return
	}
	log.Println("✅ OK")
}

// ---------- NATS ----------
func testNATS() {
	log.Print("NATS: ")
	nc, err := nats.Connect("nats://localhost:42220")
	if err != nil {
		log.Println("❌ FAILED:", err)
		return
	}
	defer nc.Close()

	if err := nc.Publish("health.check", []byte("ping")); err != nil {
		log.Println("❌ FAILED:", err)
		return
	}
	log.Println("✅ OK")
}

// ---------- Kafka ----------
func testKafka() {
	log.Print("Kafka: ")

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "health-check",
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := writer.WriteMessages(ctxTimeout,
		kafka.Message{Value: []byte("ping")},
	)
	if err != nil {
		log.Println("❌ FAILED:", err)
		return
	}

	log.Println("✅ OK")
}

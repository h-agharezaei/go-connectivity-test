package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("🔍 Starting connectivity tests...")

	testPostgres()
	testRedis()
	testMongo()
	testNATS()
	testJetStream()
	testKafka()

	log.Println("✅ All tests finished")
}

// ---------- Postgres ----------
func testPostgres() {
	log.Print("Postgres: ")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/", user, password, host, port)
	conn, err := pgx.Connect(ctx, connStr)
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
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
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
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
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
	url := os.Getenv("NATS_URL")
	nc, err := nats.Connect(url)
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

func testJetStream() {
	log.Print("NATS JetStream: ")

	url := os.Getenv("NATS_JS_URL")
	nc, err := nats.Connect(url)
	if err != nil {
		log.Println("❌ FAILED (connect):", err)
		return
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Println("❌ FAILED (jetstream not enabled):", err)
		return
	}

	// 1️⃣ Ensure Stream
	streamName := "HEALTH"
	subject := "health.jetstream"

	_, err = js.StreamInfo(streamName)
	if err != nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:      streamName,
			Subjects:  []string{subject},
			Storage:   nats.MemoryStorage,
			Retention: nats.LimitsPolicy,
		})
		if err != nil {
			log.Println("❌ FAILED (create stream):", err)
			return
		}
	}

	// 2️⃣ Publish
	_, err = js.Publish(subject, []byte("ping"))
	if err != nil {
		log.Println("❌ FAILED (publish):", err)
		return
	}

	// 3️⃣ Consumer (Pull)
	sub, err := js.PullSubscribe(subject, "health-consumer",
		nats.BindStream(streamName),
	)
	if err != nil {
		log.Println("❌ FAILED (consumer):", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	msgs, err := sub.Fetch(1, nats.Context(ctx))
	if err != nil || len(msgs) == 0 {
		log.Println("❌ FAILED (fetch):", err)
		return
	}

	msgs[0].Ack()
	log.Println("✅ OK")
}

// ---------- Kafka ----------
func testKafka() {
	log.Print("Kafka: ")

	host := os.Getenv("KAFKA_HOST")
	port := os.Getenv("KAFKA_PORT")
	conn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Println("❌ FAILED:", err)
		return
	}
	defer conn.Close()
	log.Println("✅ OK")
}

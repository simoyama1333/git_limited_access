package main

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func insertBsonD(col *mongo.Collection) error {
    bsonD := bson.D{
        {"str1", "abc"},
        {"num1", 1},
        {"str2", "xyz"},
        {"num2", bson.A{2, 3, 4}},
        {"subdoc", bson.D{{"str", "subdoc"}, {"num", 987}}},
        {"date", time.Now()},
    }
    _, err := col.InsertOne(context.Background(), bsonD)
    return err
}

func insertBsonM(col *mongo.Collection) error {
    bsonM := bson.M{
        "str1": "efg",
        "num1": 11,
        "str2": "opq",
        "num2": bson.A{12, 13, 14},
        "subdoc": bson.M{"str": "subdoc", "num": 987},
        "date": time.Now(),
    }
    for i := 0; i < 10; i++ {
        _, err := col.InsertOne(context.Background(), bsonM)
        if err != nil {
            return err
        }
    }
    return nil
}

type myType struct {
    Str1 string
    Num1 int
    Str2 string
    Num2 []int
    Subdoc struct {
        Str string
        Num int
    }
    Date time.Time
}

func insertStruct(col *mongo.Collection) error {
    doc := myType{
        "hij",
        21,
        "rst",
        []int{22, 23, 24},
        struct {
            Str string
            Num int
        }{"subdoc", 987},
        time.Now(),
    }
    _, err := col.InsertOne(context.Background(), doc)
    return err
}

func mainMain() error {
    client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        return err
    }
    if err = client.Connect(context.Background()); err != nil {
        return err
    }
    defer client.Disconnect(context.Background())

    col := client.Database("test").Collection("col")
    if err = insertBsonD(col); err != nil {
        return err
    }
    if err = insertBsonM(col); err != nil {
        return err
    }
    if err = insertStruct(col); err != nil {
        return err
    }
    return nil
}

func main() {
    if err := mainMain(); err != nil {
        log.Fatal(err)
    }
    log.Println("normal end.")
}
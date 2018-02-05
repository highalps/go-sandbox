package main

// Docker container에 mongo 실행 (https://hub.docker.com/_/mongo/)
// goroutine과 channel 이용한 다중 insert 구문 실행

import (
    "fmt"
    "runtime"
    "github.com/globalsign/mgo"
    "github.com/globalsign/mgo/bson"
)

const MAX_WORKER = 1000
const MAX_JOB = 1000000

type Character struct {
  ID         bson.ObjectId    `json:"_id" bson:"_id"`
  Name       string           `json:"name" bson:"name"` 
  Power      int              `json:"power" bson:"power"`
  Health     int              `json:"health" bson:"health"`
  Mobility   int              `json:"mobility" bson:"mobility"`
  Techniques int              `json:"techniques" bson:"techniques"`
  Ranges     int              `json:"ranges" bson:"ranges"`
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU()) // 프로세스 최대로 사용
    sem := make(chan int, MAX_WORKER)

    url := "localhost:27017"
    session, err := mgo.Dial(url)
   
    if err != nil {
        fmt.Println(err)
    }

    c := session.DB("sf5").C("character")

    for i:=0 ; i<MAX_JOB ; i++ {
        sem <- 1
        go func(idx int) {
                p := Character{
                ID: bson.NewObjectId(),
                Name: "birdie",
                Power: idx,
                Health: 5,
                Mobility: 1,
                Techniques: 3,
                Ranges: 4,
            }
            err := c.Insert(p)
            
            if err != nil {
                fmt.Println(err)
            }

            <- sem
        }(i)
    }

    fmt.Println("done")
    
}

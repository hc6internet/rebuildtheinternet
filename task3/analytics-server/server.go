package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	zmq "github.com/pebbe/zmq4"
)

type AnalyticsService struct {
	sock       *zmq.Socket
	seq        uint32
	dbUser     string
	dbPasswd   string
	dbDatabase string
}

func main() {

	// variables
	serverId := os.Getenv("SERVER_ID") // number
	port := os.Getenv("PORT")

	partitionService := os.Getenv("PARTITION_SERVICE")

	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")
	dbDatabase := os.Getenv("DB_DATABASE")

	sid, err := strconv.ParseUint(serverId, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	seq := uint32(sid) << 24 // initial seq number

	log.Printf("Initial seq = %d", seq)
	log.Printf("partition service = %s", partitionService)

	// zmq socket
	sock, err := zmq.NewSocket(zmq.REQ)
	defer sock.Close()
	err = sock.Connect("tcp://" + partitionService)
	if err != nil {
		log.Fatal("sock connect failed: " + err.Error())
	}

	s := &AnalyticsService{sock: sock, seq: seq, dbUser: dbUser, dbPasswd: dbPasswd, dbDatabase: dbDatabase}

	// v1 api routes
	r := mux.NewRouter()
	r.HandleFunc("/analytics.gif", s.handler)

	log.Println("Starting up on port " + port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
}

type invalidSeqError struct {
}

func (e *invalidSeqError) Error() string {
	return "Invalid sequence number"
}

func (s *AnalyticsService) handler(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}

	remoteAddr := strings.Split(r.RemoteAddr, ":")

	m["visitor"] = remoteAddr[0]
	m["siteId"] = r.URL.Query().Get("ID")
	m["url"] = r.URL.Query().Get("URL")
	m["userAgent"] = r.URL.Query().Get("UserAgent")
	m["referrer"] = r.URL.Query().Get("Referrer")
	m["time"] = r.URL.Query().Get("Time")

	db, err := s.databaseFromId(m["siteId"])
	if err != nil {
		log.Printf("Failed to get database")
	}

	err = s.writeToDB(db, m)
	if err != nil {
		log.Printf("Failed writing to database %s: %s", db, err.Error())
	} else {
		log.Printf("Successfully writen to database %s", db)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *AnalyticsService) databaseFromId(id string) (string, error) {
	req := &PartitionRequest{SiteId: id, Seq: s.seq}
	out, err := proto.Marshal(req)
	if err != nil {
		return "", err
	}

	ret, err := s.sock.SendBytes(out, 0)
	if err != nil {
		return "", err
	}
	log.Printf("sent %d bytes", ret)
	s.seq++

	// get response back
	resp := &PartitionResponse{}
	reply, err := s.sock.RecvBytes(0)
	if err != nil {
		return "", err
	}
	if err := proto.Unmarshal(reply, resp); err != nil {
		return "", err
	}
	if resp.Seq != req.Seq {
		return "", &invalidSeqError{}
	}
	// finally, get response
	log.Printf("partition id = %s", resp.PartitionId)
	return resp.PartitionId, nil
}

func (s *AnalyticsService) writeToDB(dbHost string, m map[string]string) error {
	// db connection
	var dbStr strings.Builder
	fmt.Fprintf(&dbStr, "%s:%s@tcp(%s:3306)/%s", s.dbUser, s.dbPasswd, dbHost, s.dbDatabase)

	db, err := sql.Open("mysql", dbStr.String())
	if err != nil {
		log.Fatal("Sql open failed: " + err.Error())
	}
	defer db.Close()

	var b strings.Builder
	fmt.Fprintf(&b, `INSERT INTO visit (siteId, visitor, url, userAgent, referrer, visitTime) 
VALUES ("%s", "%s", "%s", "%s", "%s", "%s")`,
		m["siteId"], m["visitor"], m["url"], m["userAgent"], m["referrer"], m["time"])

	insStmt, err := db.Prepare(b.String())
	if err != nil {
		return err
	}
	defer insStmt.Close()

	_, err = insStmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

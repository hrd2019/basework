package datasource

import (
	"fmt"
	"testing"
)

//func TestGetPG(t *testing.T) {
//	cnn := GetPGSql(5432, "192.168.1.222", "postgres", "psql", "postgres")
//
//	if cnn != nil {
//		println(cnn.Tables)
//	}
//}

type TestPro struct {
	msgContent string
}

// 实现发送者
func (t *TestPro) MsgContent() string {
	return t.msgContent
}

// 实现接收者
func (t *TestPro) Consumer(dataByte []byte) error {
	fmt.Println("rec:", string(dataByte))
	return nil
}

func TestRabbitMQ(t *testing.T) {
	msg := fmt.Sprintf("a test only")
	testPro := &TestPro{
		msg,
	}

	queueExchange := &Queue{
		"test",
		"test",
		"192.168.1.222",
		5672,
		"test",
		"test.xx",
		"test.xx",
		"test.xx",
		"topic",
	}
	mq := New(queueExchange)
	mq.RegisterProducer(testPro)
	mq.RegisterReceiver(testPro)
	mq.Start()
}

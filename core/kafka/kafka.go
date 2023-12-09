package CoreKafka

import(
	"github.com/Shopify/sarama"
	"time"
)

//异步消息列队服务模块

var(
	//服务地址序列
	serverAddress = []string{"127.0.0.1:9092"}
	//本地测试服务地址
	debugServerAddress = []string{"127.0.0.1:9092"}
	//是否为debug
	allowDebug = true
	//配置组
	serverConfig *sarama.Config
	//超时时间
	serverTimeoutTime int64 = 5
	//连接到服务器结果集合
	serverAsyncProduct sarama.AsyncProducer
	serverSyncProduct sarama.SyncProducer
	//消费者关系
	serverClient sarama.Client
)

//初始化
func Init(tAllowDebug bool, tServerAddress []string){
	SetDebug(tAllowDebug)
	SetServerAddress(tServerAddress)
}

//设置debug
func SetDebug(tAllowDebug bool){
	allowDebug = tAllowDebug
}

//设置服务地址序列
func SetServerAddress(tServerAddress []string){
	serverAddress = tServerAddress
}

//设置debug地址序列
func SetDebugServerAddress(tDebugServerAddress []string){
	debugServerAddress = tDebugServerAddress
}

//连接到服务
func NewProducer(isAsync bool) (*sarama.Config, error){
	var err error
	serverConfig = sarama.NewConfig()
	serverConfig.Producer.Return.Successes = true
	serverConfig.Producer.Timeout = time.Second * time.Duration(serverTimeoutTime)
	if isAsync {
		if allowDebug {
			serverAsyncProduct, err = NewAsyncProducer(serverAddress)
		}else{
			serverSyncProduct, err = NewSyncProducer(serverAddress)
		}
	}else{
		if allowDebug {
			serverAsyncProduct, err = NewAsyncProducer(debugServerAddress)
		}else{
			serverSyncProduct, err = NewSyncProducer(debugServerAddress)
		}
	}
	if err != nil{
		return serverConfig, err
	}
	return serverConfig, nil
}

//注册新的异步生产者关系
// 注意，defer关闭
func NewAsyncProducer(tServerAddressList []string) (sarama.AsyncProducer, error){
	serverProducer, err := sarama.NewAsyncProducer(tServerAddressList, serverConfig)
	/**
	defer func() {
		if err := serverProducer.Close(); err != nil {
			//
		}
	}()
	*/
	return serverProducer, err
}

//注册新的同步生产者关系
// 注意，defer关闭
func NewSyncProducer(tServerAddressList []string) (sarama.SyncProducer, error){
	serverProducer, err := sarama.NewSyncProducer(tServerAddressList, serverConfig)
	/**
	defer func() {
		if err := serverProducer.Close(); err != nil {
			//
		}
	}()
	 */
	return serverProducer, err
}

//推送消息
func SendMessage(serverProducer sarama.SyncProducer, topic string, dataByte []byte) (partition int32, offset int64, err error){
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(dataByte),
	}
	part, offset, err := serverProducer.SendMessage(message)
	return part, offset, err
}

//订阅消息
// 可以作为参考，或一般的服务消费者处理方案
func NewConsumer(tServerAddressList []string) (sarama.Client, error){
	var err error
	serverClient, err = sarama.NewClient(tServerAddressList, serverConfig)
	return serverClient, err
}

func NewConsumerChildren(){

}
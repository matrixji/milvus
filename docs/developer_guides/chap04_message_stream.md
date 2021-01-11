## 8. Message Stream Service



#### 8.1 Overview



#### 8.2 Message Stream Service API

```go
type Client interface {
  CreateChannels(req CreateChannelRequest) (CreateChannelResponse, error)
  DestoryChannels(req DestoryChannelRequest) error
  DescribeChannels(req DescribeChannelRequest) (DescribeChannelResponse, error)
}
```

* *CreateChannels*

```go
type OwnerDescription struct {
  Role string
  Address string
  //Token string
  DescriptionText string
}

type CreateChannelRequest struct {
  OwnerDescription OwnerDescription
  NumChannels int
}

type CreateChannelResponse struct {
  ChannelIDs []string
}
```

* *DestoryChannels*

```go
type DestoryChannelRequest struct {
	ChannelIDs []string
}
```



* *DescribeChannels*

```go
type DescribeChannelRequest struct {
	ChannelIDs []string
}

type ChannelDescription struct {
  ChannelID string
  Owner OwnerDescription
}

type DescribeChannelResponse struct {
  Descriptions []ChannelDescription
}
```



#### A.3 Message Stream

``` go
type MsgType uint32
const {
  kInsert MsgType = 400
  kDelete MsgType = 401
  kSearch MsgType = 500
  kSearchResult MsgType = 1000
  
  kSegStatistics MsgType = 1100
  
  kTimeTick MsgType = 1200
  kTimeSync MsgType = 1201
}

type TsMsg interface {
  SetTs(ts Timestamp)
  BeginTs() Timestamp
  EndTs() Timestamp
  Type() MsgType
  Marshal(*TsMsg) []byte
  Unmarshal([]byte) *TsMsg
}

type MsgPack struct {
  BeginTs Timestamp
  EndTs Timestamp
  Msgs []TsMsg
}


type MsgStream interface {
  Produce(*MsgPack) error
  Broadcast(*MsgPack) error
  Consume() *MsgPack // message can be consumed exactly once
}

type RepackFunc(msgs []* TsMsg, hashKeys [][]int32) map[int32] *MsgPack

type PulsarMsgStream struct {
  client *pulsar.Client
  repackFunc RepackFunc
  producers []*pulsar.Producer
  consumers []*pulsar.Consumer
  unmarshal *UnmarshalDispatcher
}

func (ms *PulsarMsgStream) CreatePulsarProducers(topics []string)
func (ms *PulsarMsgStream) CreatePulsarConsumers(subname string, topics []string, unmarshal *UnmarshalDispatcher)
func (ms *PulsarMsgStream) SetRepackFunc(repackFunc RepackFunc)
func (ms *PulsarMsgStream) Produce(msgs *MsgPack) error
func (ms *PulsarMsgStream) Broadcast(msgs *MsgPack) error
func (ms *PulsarMsgStream) Consume() (*MsgPack, error)
func (ms *PulsarMsgStream) Start() error
func (ms *PulsarMsgStream) Close() error

func NewPulsarMsgStream(ctx context.Context, pulsarAddr string) *PulsarMsgStream


type PulsarTtMsgStream struct {
  client *pulsar.Client
  repackFunc RepackFunc
  producers []*pulsar.Producer
  consumers []*pulsar.Consumer
  unmarshal *UnmarshalDispatcher
  inputBuf []*TsMsg
  unsolvedBuf []*TsMsg
  msgPacks []*MsgPack
}

func (ms *PulsarTtMsgStream) CreatePulsarProducers(topics []string)
func (ms *PulsarTtMsgStream) CreatePulsarConsumers(subname string, topics []string, unmarshal *UnmarshalDispatcher)
func (ms *PulsarTtMsgStream) SetRepackFunc(repackFunc RepackFunc)
func (ms *PulsarTtMsgStream) Produce(msgs *MsgPack) error
func (ms *PulsarTtMsgStream) Broadcast(msgs *MsgPack) error
func (ms *PulsarTtMsgStream) Consume() *MsgPack //return messages in one time tick
func (ms *PulsarTtMsgStream) Start() error
func (ms *PulsarTtMsgStream) Close() error

func NewPulsarTtMsgStream(ctx context.Context, pulsarAddr string) *PulsarTtMsgStream
```



```go
type MarshalFunc func(*TsMsg) []byte
type UnmarshalFunc func([]byte) *TsMsg


type UnmarshalDispatcher struct {
	tempMap map[ReqType]UnmarshalFunc 
}

func (dispatcher *MarshalDispatcher) Unmarshal([]byte) *TsMsg
func (dispatcher *MarshalDispatcher) AddMsgTemplate(msgType MsgType, marshal MarshalFunc)
func (dispatcher *MarshalDispatcher) addDefaultMsgTemplates()

func NewUnmarshalDispatcher() *UnmarshalDispatcher
```



#### A.4 RocksMQ

RocksMQ is a RocksDB-based messaging/streaming library.

```go
type ProducerMessage struct {
  Key string
  Payload []byte
} 
```

```go
type ConsumerMessage struct {
  MsgID MessageID
  Key string
  Payload []byte
} 
```



```GO
type RocksMQ struct {
  CreateChannel(channelName string) error
  DestroyChannel(channelName string) error
  CreateConsumerGroup(groupName string) error
  DestroyConsumerGroup(groupName string) error
  
  Produce(channelName string, messages []ProducerMessage) error
  Consume(groupName string, channelName string, n int) ([]ConsumerMessage, error)
  Seek(groupName string, channelName string, msgID MessageID) error
}
```



##### A.4.1 Meta

* channel meta

```go
"$(channel_name)/start_id", MessageID
"$(channel_name)/end_id", MessageID
```

* consumer group meta

```go
"$(group_name)/$(channel_name)/current_id", MessageID
```

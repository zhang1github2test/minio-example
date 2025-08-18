package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

// Config 配置结构体
type Config struct {
	SFTPSources []SFTPConfig `yaml:"sftp_sources"`
	SFTPTarget  SFTPConfig   `yaml:"sftp_target"`
	Kafka       KafkaConfig  `yaml:"kafka"`
	Interval    int          `yaml:"interval"` // 扫描间隔，秒
}

// SFTPConfig SFTP配置
type SFTPConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	InDir      string `yaml:"in_dir"`
	OutDir     string `yaml:"out_dir"`
	PrivateKey string `yaml:"private_key,omitempty"`
}

// KafkaConfig Kafka配置
type KafkaConfig struct {
	Brokers      []string `yaml:"brokers"`
	ProduceTopic string   `yaml:"produce_topic"` // 生产者发送消息的主题
	ConsumeTopic string   `yaml:"consume_topic"` // 消费者接收消息的主题
	GroupID      string   `yaml:"group_id"`
}

// FileMessage 文件消息结构
type FileMessage struct {
	SourceHost string
	Filename   string
	Content    string
	Timestamp  time.Time
}

// SFTPClient SFTP客户端封装
type SFTPClient struct {
	client *sftp.Client
	conn   *ssh.Client
	config SFTPConfig
}

// NewSFTPClient 创建SFTP客户端
func NewSFTPClient(config SFTPConfig) (*SFTPClient, error) {
	var authMethod ssh.AuthMethod
	if config.PrivateKey != "" {
		key, err := os.ReadFile(config.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key: %v", err)
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %v", err)
		}
		authMethod = ssh.PublicKeys(signer)
	} else {
		authMethod = ssh.Password(config.Password)
	}

	sshConfig := &ssh.ClientConfig{
		User:            config.Username,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server: %v", err)
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create SFTP client: %v", err)
	}

	return &SFTPClient{
		client: client,
		conn:   conn,
		config: config,
	}, nil
}

// Close 关闭连接
func (s *SFTPClient) Close() error {
	if s.client != nil {
		s.client.Close()
	}
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// ReadFiles 读取指定目录下的.mh文件
func (s *SFTPClient) ReadFiles() ([]FileMessage, error) {
	inDir := s.config.InDir
	if inDir == "" {
		inDir = "/IN"
	}

	files, err := s.client.ReadDir(inDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %v", inDir, err)
	}

	var messages []FileMessage
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".mh") {
			continue
		}

		filePath := inDir + file.Name()
		content, err := s.readFile(filePath)
		if err != nil {
			log.Printf("Failed to read file %s: %v", filePath, err)
			continue
		}

		fmt.Printf("=== 读取文件信息 ===\n")
		fmt.Printf("来源地址: %s:%d\n", s.config.Host, s.config.Port)
		fmt.Printf("文件路径: %s\n", filePath)
		fmt.Printf("文件名: %s\n", file.Name())
		fmt.Printf("文件内容:\n%s\n", content)

		message := FileMessage{
			SourceHost: fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
			Filename:   file.Name(),
			Content:    content,
			Timestamp:  time.Now(),
		}

		messages = append(messages, message)

		// 删除已处理的文件
		if err := s.client.Remove(filePath); err != nil {
			log.Printf("Failed to delete file %s: %v", filePath, err)
		} else {
			fmt.Printf("已删除源文件: %s\n", filePath)
		}
	}

	return messages, nil
}

// readFile 读取文件内容
func (s *SFTPClient) readFile(filePath string) (string, error) {
	file, err := s.client.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// WriteFile 写入文件到目标目录
func (s *SFTPClient) WriteFile(filename, content string) error {
	outDir := s.config.OutDir
	if outDir == "" {
		outDir = "/OUT/"
	}

	// 确保目录存在
	if err := s.client.MkdirAll(outDir); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", outDir, err)
	}

	filePath := outDir + filename
	file, err := s.client.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	if _, err := file.Write([]byte(content)); err != nil {
		return fmt.Errorf("failed to write file %s: %v", filePath, err)
	}

	fmt.Printf("=== 文件上传完成 ===\n")
	fmt.Printf("目标地址: %s:%d\n", s.config.Host, s.config.Port)
	fmt.Printf("上传路径: %s\n", filePath)
	fmt.Printf("文件名: %s\n", filename)

	return nil
}

// KafkaProducer Kafka生产者
type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewKafkaProducer 创建Kafka生产者
func NewKafkaProducer(config KafkaConfig) (*KafkaProducer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}

	return &KafkaProducer{
		producer: producer,
		topic:    config.ProduceTopic,
	}, nil
}

// SendMessage 发送消息到Kafka
func (k *KafkaProducer) SendMessage(message FileMessage) error {
	fmt.Printf("=== 发送Kafka消息 ===\n")
	fmt.Printf("主题: %s\n", k.topic)
	fmt.Printf("文件名: %s\n", message.Filename)
	fmt.Printf("来源: %s\n", message.SourceHost)
	fmt.Printf("消息内容: %s\n", message.Content)

	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Key:   sarama.StringEncoder(message.Filename), // 使用文件名作为消息key
		Value: sarama.StringEncoder(message.Content),  // 直接发送文件内容
	}

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %v", err)
	}

	fmt.Printf("消息发送成功 - 分区: %d, 偏移量: %d\n", partition, offset)
	return nil
}

// Close 关闭生产者
func (k *KafkaProducer) Close() error {
	return k.producer.Close()
}

// KafkaConsumer Kafka消费者
type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
	topic    string
	groupID  string
}

// NewKafkaConsumer 创建Kafka消费者
func NewKafkaConsumer(config KafkaConfig) (*KafkaConsumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaConfig.Consumer.Group.Session.Timeout = 10 * time.Second
	saramaConfig.Consumer.Group.Heartbeat.Interval = 3 * time.Second

	consumer, err := sarama.NewConsumerGroup(config.Brokers, config.GroupID, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %v", err)
	}

	return &KafkaConsumer{
		consumer: consumer,
		topic:    config.ConsumeTopic,
		groupID:  config.GroupID,
	}, nil
}

// ConsumerGroupHandler 消费者组处理器
type ConsumerGroupHandler struct {
	sftpClient *SFTPClient
}

// Setup 设置消费者
func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup 清理消费者
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 消费消息
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			fmt.Printf("=== 接收到Kafka消息 ===\n")
			fmt.Printf("分区: %d, 偏移量: %d\n", message.Partition, message.Offset)

			// 从消息key获取文件名，如果没有key则生成默认文件名
			filename := string(message.Key)
			if filename == "" {
				filename = fmt.Sprintf("file_%d_%d.mh", message.Partition, message.Offset)
			}

			content := string(message.Value)
			fmt.Printf("文件名: %s\n", filename)
			fmt.Printf("消息内容: %s\n", content)

			// 写入文件到SFTP
			if err := h.sftpClient.WriteFile(filename, content); err != nil {
				log.Printf("Failed to write file to SFTP: %v", err)
				continue
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

// Consume 开始消费
func (k *KafkaConsumer) Consume(ctx context.Context, sftpClient *SFTPClient) error {
	handler := &ConsumerGroupHandler{sftpClient: sftpClient}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := k.consumer.Consume(ctx, []string{k.topic}, handler); err != nil {
				log.Printf("Error consuming messages: %v", err)
				time.Sleep(5 * time.Second)
			}
		}
	}
}

// Close 关闭消费者
func (k *KafkaConsumer) Close() error {
	return k.consumer.Close()
}

// loadConfig 加载配置文件
func loadConfig(configFile string) (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// 设置默认值
	if config.Interval == 0 {
		config.Interval = 10
	}

	return &config, nil
}

func startFtpKafkaTask() {
	var configFile string
	flag.StringVar(&configFile, "config", "config.yaml", "配置文件路径")
	flag.Parse()

	// 加载配置
	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("=== 启动SFTP-Kafka文件传输工具 ===\n")
	fmt.Printf("配置文件: %s\n", configFile)
	fmt.Printf("扫描间隔: %d秒\n", config.Interval)

	// 创建Kafka生产者
	producer, err := NewKafkaProducer(config.Kafka)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// 创建目标SFTP客户端
	targetSFTP, err := NewSFTPClient(config.SFTPTarget)
	if err != nil {
		log.Fatalf("Failed to create target SFTP client: %v", err)
	}
	defer targetSFTP.Close()

	// 创建Kafka消费者
	consumer, err := NewKafkaConsumer(config.Kafka)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// 启动消费者协程
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := consumer.Consume(ctx, targetSFTP); err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()

	// 定期扫描源SFTP目录
	ticker := time.NewTicker(time.Duration(config.Interval) * time.Second)
	defer ticker.Stop()
	filepath.Join()
	for {
		select {
		case <-ticker.C:
			fmt.Printf("\n=== 开始扫描源SFTP目录 ===\n")
			for i, sftpConfig := range config.SFTPSources {
				fmt.Printf("扫描源 %d: %s:%d\n", i+1, sftpConfig.Host, sftpConfig.Port)

				client, err := NewSFTPClient(sftpConfig)
				if err != nil {
					log.Printf("Failed to create SFTP client for %s:%d: %v",
						sftpConfig.Host, sftpConfig.Port, err)
					continue
				}

				messages, err := client.ReadFiles()
				client.Close()

				if err != nil {
					log.Printf("Failed to read files from %s:%d: %v",
						sftpConfig.Host, sftpConfig.Port, err)
					continue
				}

				// 发送消息到Kafka
				for _, message := range messages {
					if err := producer.SendMessage(message); err != nil {
						log.Printf("Failed to send message to Kafka: %v", err)
					}
				}
			}
		case <-ctx.Done():
			wg.Wait()
			return
		}
	}

}

func main() {
	startFtpKafkaTask()
}

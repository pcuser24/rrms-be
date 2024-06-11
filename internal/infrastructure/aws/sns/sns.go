package sns

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

type SNSClient interface {
	GetClient() *sns.Client
	ListTopics(ctx context.Context) ([]types.Topic, error)
	CreateTopic(ctx context.Context, topicName string, isFifoTopic bool, contentBasedDeduplication bool) (string, error)
	DeleteTopic(ctx context.Context, topicArn string) error
	SubscribeQueue(topicArn string, queueArn string, filterMap map[string][]string) (string, error)
	Publish(ctx context.Context, subject, message, topicArn string, attrs map[string]interface{}) error
}

type snsClient struct {
	snsClient *sns.Client
}

func NewSNSClient(conf *aws.Config) SNSClient {
	return &snsClient{
		snsClient: sns.NewFromConfig(*conf),
	}
}

func (c *snsClient) GetClient() *sns.Client {
	return c.snsClient
}

func (c *snsClient) ListTopics(ctx context.Context) ([]types.Topic, error) {
	var topics []types.Topic
	paginator := sns.NewListTopicsPaginator(c.snsClient, &sns.ListTopicsInput{})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		} else {
			topics = append(topics, output.Topics...)
		}
	}
	return topics, nil
}

func (c *snsClient) CreateTopic(ctx context.Context, topicName string, isFifoTopic bool, contentBasedDeduplication bool) (string, error) {
	topicAttributes := map[string]string{}
	if isFifoTopic {
		topicAttributes["FifoTopic"] = "true"
	}
	if contentBasedDeduplication {
		topicAttributes["ContentBasedDeduplication"] = "true"
	}
	res, err := c.snsClient.CreateTopic(ctx, &sns.CreateTopicInput{
		Name:       aws.String(topicName),
		Attributes: topicAttributes,
	})
	if err != nil {
		return "", err
	}
	return *res.TopicArn, nil
}

func (c *snsClient) DeleteTopic(ctx context.Context, topicArn string) error {
	_, err := c.snsClient.DeleteTopic(ctx, &sns.DeleteTopicInput{
		TopicArn: aws.String(topicArn)})
	return err
}

func (c *snsClient) Publish(ctx context.Context, subject, message, topicArn string, attrs map[string]interface{}) error {
	msgAttributes := make(map[string]types.MessageAttributeValue)
	for k, v := range attrs {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		msgAttributes[k] = types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(string(b)),
		}
	}

	input := sns.PublishInput{
		Subject:           aws.String(subject),
		Message:           aws.String(message),
		TopicArn:          aws.String(topicArn),
		MessageAttributes: msgAttributes,
	}
	_, err := c.snsClient.Publish(ctx, &input)
	return err
}

// SubscribeQueue subscribes an Amazon Simple Queue Service (Amazon SQS) queue to an
// Amazon SNS topic. When filterMap is not nil, it is used to specify a filter policy
// so that messages are only sent to the queue when the message has the specified attributes.
func (c *snsClient) SubscribeQueue(topicArn string, queueArn string, filterMap map[string][]string) (string, error) {
	var subscriptionArn string
	var attributes map[string]string
	if filterMap != nil {
		filterBytes, err := json.Marshal(filterMap)
		if err != nil {
			return "", err
		}
		attributes = map[string]string{"FilterPolicy": string(filterBytes)}
	}
	output, err := c.snsClient.Subscribe(context.TODO(), &sns.SubscribeInput{
		Protocol:              aws.String("sqs"),
		TopicArn:              aws.String(topicArn),
		Attributes:            attributes,
		Endpoint:              aws.String(queueArn),
		ReturnSubscriptionArn: true,
	})
	if err != nil {
		return "", err
	} else {
		subscriptionArn = *output.SubscriptionArn
	}

	return subscriptionArn, err
}

package metric

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/storage"
)

// compiletime assert that we implement the right interface
var _ Metric = (*AzureQueueMetric)(nil)

var azureAccountName string
var azureQueueClient storage.QueueServiceClient
var azureClient storage.Client
var azureInitialized = false

// AzureQueueMetric stores the current value.
type AzureQueueMetric struct {
	currentVal     int
	azureQueueName string
}

// AcsInit sets up the Azure Storage client.
func AcsInit() (err error) {
	azureAccountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	if azureAccountName == "" {
		log.Error("You need to pass in environment variable AZURE_STORAGE_ACCOUNT_NAME")
		err = fmt.Errorf("Azure storage account name not configured")
		return
	}
	azureKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	if azureAccountName == "" {
		log.Error("You need to pass in environment variable AZURE_STORAGE_ACCOUNT_KEY")
		err = fmt.Errorf("Azure storage account key not configured")
		return
	}

	azureClient, err = storage.NewBasicClient(azureAccountName, azureKey)
	if err == nil {
		azureQueueClient = azureClient.GetQueueService()
	}

	azureInitialized = true
	return
}

// NewAzureQueueMetric creates the metric.
func NewAzureQueueMetric(queueName string) *AzureQueueMetric {
	if !azureInitialized {
		AcsInit()
	}

	return &AzureQueueMetric{
		azureQueueName: queueName,
	}
}

// UpdateCurrent calls the Azure Storage API to the queue length.
func (aqm *AzureQueueMetric) UpdateCurrent() {
	metadata, err := azureQueueClient.GetMetadata(aqm.azureQueueName)
	if err != nil {
		log.Errorf("Error getting Azure queue info: %v", err)
	}
	aqm.currentVal = metadata.ApproximateMessageCount
	log.Debugf("Queue name %s length %d", aqm.azureQueueName, aqm.currentVal)
}

// Current returns the queue length.
func (aqm *AzureQueueMetric) Current() int {
	return aqm.currentVal
}

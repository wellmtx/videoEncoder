package services

import (
	"log"
	"os"
	"strconv"
	"videoEncoder/app/repositories"
	"videoEncoder/domain"
	"videoEncoder/framework/queue"

	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

type JobManager struct {
	Db               *gorm.DB
	Domain           domain.Job
	MessageChannel   chan amqp.Delivery
	JobReturnChannel chan JobWorkerResult
	RabbitMQ         *queue.RabbitMQ
}

type JobNotificationError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewJobManager(db *gorm.DB, rabbitMQ *queue.RabbitMQ, jobReturnChannel chan JobWorkerResult, messageChannel chan amqp.Delivery) *JobManager {
	return &JobManager{
		Db:               db,
		Domain:           domain.Job{},
		MessageChannel:   messageChannel,
		JobReturnChannel: jobReturnChannel,
		RabbitMQ:         rabbitMQ,
	}
}

func (j *JobManager) Start(ch *amqp.Channel) {
	videoService := NewVideoService()
	videoService.VideoRepository = repositories.VideoRepositoryDb{Db: j.Db}

	jobService := JobService{
		JobRepository: repositories.JobRepositoryDb{Db: j.Db},
		VideoService:  videoService,
	}

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKERS"))

	if err != nil {
		log.Fatalf("error loading var: CONCURRENCY_WORKERS")
	}

	for qtdProcesses := 0; qtdProcesses < concurrency; qtdProcesses++ {
		go JobWorker(j.MessageChannel, j.JobReturnChannel, jobService, j.Domain, qtdProcesses)
	}

	for jobResult := range j.JobReturnChannel {
		if jobResult.Error != nil {
			err = j.checkParseErrors(jobResult)
		} else {
			err = j.notifySuccess(jobResult, ch)
		}

		if err != nil {
			jobResult.Message.Reject(false)
		}
	}
}

func (j *JobManager) checkParseErrors(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		log.Printf("MessageID #{jobResult.Message.DeliveryTag}. Error parsing job: #{jobResult.Job.ID}")
	} else {
		log.Printf("MessageID #{jobResult.Message.DeliveryTag}. Error parsing message: #{jobResult.Error}")
	}

	// errorMsg := JobNotificationError{
	// 	Message: string(jobResult.Message.Body),
	// 	Error:   jobResult.Error.Error(),
	// }

	// jobJson, err := json.Marshal(errorMsg)

	// if err != nil {
	// 	return err
	// }

	return nil
}

package process

import (
  log "malawi-callback/logger"
  "malawi-callback/models"
  repo "malawi-callback/repository"
  "malawi-callback/repository/mssql"
  "malawi-callback/utils"
  "context"
  "github.com/aws/aws-lambda-go/events"
)

// Controller container
type Controller struct {
  secretHolder *SecretIDHolder
  sqsProducer  *SQSProducer
  repository   *repo.Repository
  config       *models.SecretModel
  requestId    *string
}

func NewController(secret string) *Controller {
  controller := Controller{
    requestId: utils.StringPtr("ROOT"),
  }
  controller.initSecret(secret)
  controller.initSqsProducer()
  return &controller
}

func (c *Controller) initSecret(secret string) {
  c.secretHolder = &SecretIDHolder{
    SecretID: secret,
    Client: CreateSMClient(),
  }
  c.config = c.secretHolder.LoadSecret()
}

func (c *Controller) initSqsProducer() {
  var err error
  c.sqsProducer, err = NewSQSProducerFromUrl(context.TODO(), CreateSQSClient(),
    &c.config.Sender.Url)
  if err != nil {
    log.Fatalf( *c.requestId, "Lambda init failed on sqs producer: %v", err)
  }
}

func (c *Controller) initRepository() {
  localRepo, err := mssql.NewRepository(c.config.Database.Config)
  if err != nil {
    log.Fatalf( *c.requestId, "Lambda init failed on repository: %v", err)
  }
  c.repository = &localRepo
}

func (c *Controller) ShutDown() {
  c.config = nil
  c.sqsProducer = nil
  c.secretHolder = nil
  if c.repository != nil {
    (*c.repository).Close()
  }
  c.repository = nil
}

func (c *Controller) PreProcess(pid *string) {
  c.requestId = pid
  // do what is needed before processing the request
}

func (c *Controller) PostProcess() {
  c.requestId = utils.StringPtr("ROOT")
}

func (c *Controller) Process(ctx context.Context, message events.SQSMessage) error {
  var err error
  // TODO Process logic here

  return err
}

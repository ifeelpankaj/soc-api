package app

import (
	"context"
	"go-server/internal/backup"
	"go-server/internal/config"
	handlers "go-server/internal/handlers/v1"
	"go-server/internal/jobs"
	"go-server/internal/models"
	repository "go-server/internal/repositories"
	service "go-server/internal/services"
	authsvc "go-server/internal/services/authSvc"
	bootstrapsvc "go-server/internal/services/bootstrapSvc"
	flatauthz "go-server/internal/services/flatAuthz"
	flatsvc "go-server/internal/services/flatSvc"
	hubsvc "go-server/internal/services/hubSvc"
	imagesvc "go-server/internal/services/imageSvc"
	maintenancesvc "go-server/internal/services/maintenanceSvc"
	"go-server/internal/services/maintenanceSvc/documentpdf"
	notificationsvc "go-server/internal/services/notificationSvc"
	plansvc "go-server/internal/services/planSvc"
	shortlinksvc "go-server/internal/services/shortLinkSvc"
	societysvc "go-server/internal/services/societySvc"
	subscriptionsvc "go-server/internal/services/subscriptionSvc"
	visitorentrysvc "go-server/internal/services/visitorEntrySvc"
	visitorsettingsvc "go-server/internal/services/visitorSettingSvc"
	"time"

	"go-server/pkg/database"
	"go-server/pkg/logger"

	"go.uber.org/zap"
)

type Dependencies struct {
	Users   repository.UserRepository
	V1      *V1Handlers
	V2      *V2Handlers
	Jobs    *jobs.Manager
	Backups *backup.Service

	Society        societysvc.SocietyService
	Flat           flatsvc.FlatService
	Plan           plansvc.PlanService
	Subscription   subscriptionsvc.SubscriptionService
	VisitorInvite  visitorentrysvc.VisitorInviteService
	VisitorEntry   visitorentrysvc.VisitorEntryService
	VisitorSetting visitorsettingsvc.VisitorSettingService
	Notification   notificationsvc.NotificationService
}

type V1Handlers struct {
	Hub                 *handlers.HubHandler
	MaintenanceDocument *handlers.MaintenanceDocumentHandler
	MaintenancePayment  *handlers.MaintenancePaymentHandler
	Maintenance         *handlers.MaintenanceHandler
	Image               *handlers.ImageHandler
	Auth                *handlers.AuthHandler
	Bootstrap           *handlers.BootstrapHandler
	Society             *handlers.SocietyHandler
	Flat                *handlers.FlatHandler
	MemberInvite        *handlers.MemberInviteHandler
	Plan                *handlers.PlanHandler
	Subscription        *handlers.SubscriptionHandler
	VisitorEntry        *handlers.VisitorEntryHandler
	VisitorSetting      *handlers.VisitorSettingHandler
	ShortLink           *handlers.ShortLinkHandler
	Notification        *handlers.NotificationHandler
	JobWebhook          *handlers.JobWebhookHandler
}

type V2Handlers struct{}

func (d *Dependencies) Shutdown() {
	logger.Info("shutting down background jobs")
	if d.Jobs != nil && !d.Jobs.Shutdown(30*time.Second) {
		logger.Warn("timed out waiting for background jobs to stop")
	}
	if d.Notification != nil {
		if err := d.Notification.Close(); err != nil {
			logger.Warn("failed to close notification service", zap.Error(err))
		}
	}
	logger.Info("background jobs stopped")
}

// StartInternalJobs starts time-sensitive jobs after the HTTP listener binds.
func (d *Dependencies) StartInternalJobs() bool {
	if d.Backups != nil {
		if err := d.Backups.StartRecovery(); err != nil {
			return false
		}
	}
	return d.Jobs != nil && d.Jobs.StartExpiry(15*time.Minute)
}

func InitializeDependencies(db *database.Database, cfg *config.Config) (*Dependencies, error) {
	logger.Info("initializing application dependencies")

	userRepo := repository.NewUserRepository(db)
	verificationRepo := repository.NewVerificationRepository(db)
	societyRepo := repository.NewSocietyRepository(db)
	societyMemberRepo := repository.NewSocietyMemberRepository(db)
	flatRepo := repository.NewFlatRepository(db)
	flatClaimRepo := repository.NewFlatClaimRepository(db)
	flatResidentRepo := repository.NewFlatResidentRepository(db)
	flatMemberInviteRepo := repository.NewFlatMemberInviteRepository(db)
	planRepo := repository.NewPlanRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	visitorSettingRepo := repository.NewVisitorSettingRepository(db)
	visitorRepo := repository.NewVisitorRepository(db)
	visitorInviteRepo := repository.NewVisitorInviteRepository(db)
	visitorEntryRepo := repository.NewVisitorEntryRepository(db)
	visitorEntryEventRepo := repository.NewVisitorEntryEventRepository(db)
	shortLinkRepo := repository.NewShortLinkRepository(db)
	deviceTokenRepo := repository.NewDeviceTokenRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	txManager := repository.NewTransactionManager(db)
	jobRepo := repository.NewJobRepository(db, txManager)

	emailSvc, err := service.NewEmailService(cfg)
	if err != nil {
		logger.Error("failed to initialize email service", zap.Error(err))
		return nil, err
	}

	subscriptionSvc := subscriptionsvc.NewSubscriptionService(subscriptionRepo, societyRepo)
	notificationSvc, err := notificationsvc.NewNotificationService(
		context.Background(),
		deviceTokenRepo,
		cfg,
		notificationRepo,
		flatResidentRepo,
		societyMemberRepo,
		subscriptionSvc,
	)
	if err != nil {
		logger.Error("failed to initialize notification service", zap.Error(err))
		return nil, err
	}

	registrationSvc := authsvc.NewRegistrationService(
		userRepo,
		verificationRepo,
		txManager,
		emailSvc,
		&cfg.Auth,
	)

	verificationSvc := authsvc.NewVerificationService(
		userRepo,
		verificationRepo,
		txManager,
		emailSvc,
		&cfg.Auth,
	)

	sessionSvc := authsvc.NewSessionService(
		userRepo,
		&cfg.Auth,
	)

	passwordSvc := authsvc.NewPasswordService(
		userRepo,
		verificationRepo,
		txManager,
		emailSvc,
		&cfg.Auth,
	)

	planSvc := plansvc.NewPlanService(planRepo)
	flatVisitorAuthz := flatauthz.New(societyMemberRepo, flatResidentRepo, flatRepo)
	visitorSettingSvc := visitorsettingsvc.NewVisitorSettingService(
		visitorSettingRepo,
		societyMemberRepo,
		flatResidentRepo,
		flatRepo,
		flatVisitorAuthz,
		txManager,
	)
	shortLinkSvc := shortlinksvc.NewShortLinkService(shortLinkRepo, cfg.PublicAppURL)
	visitorInviteSvc, visitorEntrySvc := visitorentrysvc.NewVisitorService(
		visitorRepo,
		visitorInviteRepo,
		visitorEntryRepo,
		visitorEntryEventRepo,
		visitorSettingSvc,
		societyMemberRepo,
		flatResidentRepo,
		flatRepo,
		societyRepo,
		flatVisitorAuthz,
		notificationSvc,
		shortLinkSvc,
		txManager,
	)

	societySvc := societysvc.NewSocietyService(
		societyRepo,
		societyMemberRepo,
		flatRepo,
		userRepo,
		txManager,
		flatResidentRepo,
		flatClaimRepo,
		planSvc,
		subscriptionSvc,
		visitorSettingSvc,
		visitorEntrySvc,
	)

	flatSvc := flatsvc.NewFlatService(
		flatRepo,
		flatClaimRepo,
		flatResidentRepo,
		flatMemberInviteRepo,
		societyMemberRepo,
		txManager,
		societySvc,
		subscriptionSvc,
		cfg.FlatMemberInviteTTL,
		visitorSettingSvc,
		flatVisitorAuthz,
		notificationSvc,
		shortLinkSvc,
		registrationSvc,
		sessionSvc,
	)

	bootstrapSvc := bootstrapsvc.NewBootstrapService(
		sessionSvc,
		societySvc,
		flatSvc,
	)

	jobLocation, err := time.LoadLocation(cfg.JobTimezone)
	if err != nil {
		return nil, err
	}
	expiryJob := jobs.NewExpiryJob(jobRepo)
	var cleanupImageStorage models.ImageStorage
	if cfg.ImageKit.Enabled {
		cleanupImageStorage = imagesvc.NewImageKitStorage(cfg.ImageKit)
	}

	hubConfig, err := config.LoadHubConfig()
	if err != nil {
		return nil, err
	}
	hubService := hubsvc.New(repository.NewHubRepository(db, txManager), subscriptionSvc, cleanupImageStorage, cfg.ImageKit.FolderPrefix, hubConfig)
	cleanupJob := jobs.NewCleanupJob(jobRepo, jobs.CleanupJobConfig{
		ImageStorage:              cleanupImageStorage,
		HubCleanup:                hubService.CleanupUploads,
		Location:                  jobLocation,
		Timezone:                  cfg.JobTimezone,
		VisitorEntryRetention:     time.Duration(cfg.VisitorEntryRetentionDays) * 24 * time.Hour,
		VisitorInviteRetention:    time.Duration(cfg.VisitorInviteRetentionDays) * 24 * time.Hour,
		FlatMemberInviteRetention: time.Duration(cfg.FlatInviteRetentionDays) * 24 * time.Hour,
		NotificationRetention:     time.Duration(cfg.NotificationRetentionDays) * 24 * time.Hour,
		BatchSize:                 int32(cfg.CleanupBatchSize), //nolint:gosec // Config validation bounds this value.
	})

	monthlyVisitorReportJob := jobs.NewMonthlyVisitorReportJob(jobRepo, emailSvc, jobs.MonthlyVisitorReportJobConfig{
		Location:        jobLocation,
		BatchSize:       int32(cfg.VisitorReportBatchSize), //nolint:gosec // Config validation bounds this value.
		ProcessingLease: cfg.VisitorReportProcessingLease,
		AllowResend:     cfg.VisitorReportAllowResend,
	})
	maintenanceStatuses, err := cfg.MaintenanceStatuses()
	if err != nil {
		return nil, err
	}
	maintenanceRepo := repository.NewMaintenanceRepository(db, txManager)
	paymentSvc := maintenancesvc.NewPaymentService(repository.NewMaintenancePaymentRepository(db, txManager), subscriptionSvc)
	// Maintenance billing/reminder jobs must deliver via notification_outbox only (see EnsureMaintenancePushDeliverer).
	if err := maintenancesvc.EnsureMaintenancePushDeliverer(notificationSvc); err != nil {
		return nil, err
	}
	maintenanceSvc := maintenancesvc.New(maintenanceRepo, subscriptionSvc, maintenanceStatuses, notificationRepo, notificationSvc)
	maintenanceJob := jobs.NewMaintenanceBillingJob(maintenanceSvc, jobRepo)
	reminderJob := jobs.NewMaintenanceReminderJob(maintenanceSvc, jobRepo)
	jobManager := jobs.NewManager(context.Background(), expiryJob, cleanupJob, monthlyVisitorReportJob, maintenanceJob)
	jobManager.RegisterMaintenanceReminders(reminderJob)
	backupService := backup.NewService(cfg.Backup, cfg.DBName, &backup.PGStore{Pool: db.Pool}, &backup.PGLocker{Pool: db.Pool}, backup.NewPostgres(cfg, db.Pool), backup.NewDrive(cfg.Backup), jobManager)

	v1Handlers := &V1Handlers{
		Hub:                 handlers.NewHubHandler(hubService),
		MaintenanceDocument: handlers.NewMaintenanceDocumentHandler(maintenancesvc.NewDocumentService(maintenanceSvc, paymentSvc, documentpdf.New(), subscriptionSvc)),
		MaintenancePayment:  handlers.NewMaintenancePaymentHandler(paymentSvc),
		Maintenance:         handlers.NewMaintenanceHandler(maintenanceSvc),
		Image:               handlers.NewImageHandler(newImageService(db, cfg, userRepo, societyMemberRepo, visitorEntryRepo, subscriptionSvc, visitorEntrySvc)),
		Auth: handlers.NewAuthHandler(
			registrationSvc,
			verificationSvc,
			sessionSvc,
			passwordSvc,
			&cfg.Auth,
		),
		Bootstrap:      handlers.NewBootstrapHandler(bootstrapSvc),
		Society:        handlers.NewSocietyHandler(societySvc),
		Flat:           handlers.NewFlatHandler(flatSvc),
		MemberInvite:   handlers.NewMemberInviteHandler(flatSvc),
		Plan:           handlers.NewPlanHandler(planSvc),
		Subscription:   handlers.NewSubscriptionHandler(subscriptionSvc),
		VisitorEntry:   handlers.NewVisitorEntryHandler(societySvc, visitorInviteSvc, visitorEntrySvc),
		VisitorSetting: handlers.NewVisitorSettingHandler(visitorSettingSvc),
		ShortLink:      handlers.NewShortLinkHandler(shortLinkSvc, visitorInviteSvc, flatSvc),
		Notification:   handlers.NewNotificationHandler(notificationSvc),
		JobWebhook:     handlers.NewJobWebhookHandler(jobManager, backupService),
	}

	v2Handlers := &V2Handlers{}

	logger.Info("dependencies initialized successfully")

	return &Dependencies{
		Users:          userRepo,
		V1:             v1Handlers,
		V2:             v2Handlers,
		Jobs:           jobManager,
		Backups:        backupService,
		Society:        societySvc,
		Flat:           flatSvc,
		Plan:           planSvc,
		Subscription:   subscriptionSvc,
		VisitorInvite:  visitorInviteSvc,
		VisitorEntry:   visitorEntrySvc,
		VisitorSetting: visitorSettingSvc,
		Notification:   notificationSvc,
	}, nil
}

func newImageService(db *database.Database, cfg *config.Config, users repository.UserRepository, members repository.SocietyMemberRepository, entries repository.VisitorEntryRepository, operational imagesvc.OperationalGuard, flats imagesvc.FlatEntryReader) *imagesvc.Service {
	var storage models.ImageStorage
	if cfg.ImageKit.Enabled {
		storage = imagesvc.NewImageKitStorage(cfg.ImageKit)
	}
	return imagesvc.New(storage, repository.NewImageRepository(db), imagesvc.NewAuthorization(users, members, entries, operational, flats), cfg.ImageKit.FolderPrefix, time.Now)
}

package service

import (
	"encoding/json"
	"fmt"
	"time"
	"vcs-sms/util"
)

type ReportService struct {
	esService             *ESService
	mailService           *MailService
	registeredMailService *RegisteredMailService
	serverService         *ServerService
}

func NewReportService(esService *ESService, mailService *MailService, registeredMailService *RegisteredMailService, serverService *ServerService) *ReportService {
	return &ReportService{
		esService:             esService,
		mailService:           mailService,
		registeredMailService: registeredMailService,
		serverService:         serverService,
	}
}

func (service *ReportService) SendReport(startMils int64, endMils int64, to []string) error {
	uptime_info := service.esService.CalculateUptime(startMils, endMils)
	servers := service.serverService.GetAllServers()
	onlineCount := 0
	for _, server := range servers {
		if server.Status == 1 {
			onlineCount++
		}
	}
	content, err := json.Marshal(uptime_info)
	if err != nil {
		return err
	}
	mailContent := string(content)
	mailContent += fmt.Sprintf("\nTotal servers: %d", len(servers))
	mailContent += fmt.Sprintf("\nOnline servers: %d", onlineCount)
	mailContent += fmt.Sprintf("\nOffline servers: %d", (len(servers) - onlineCount))
	err = service.mailService.SendEmail(to, string(mailContent))
	if err != nil {
		return err
	}
	return nil
}

func (service *ReportService) PeriodicReport(interval time.Duration) {
	jobTicker := util.JobTicker{
		INTERVAL_PERIOD: interval,
		HOUR_TO_TICK:    11,
		MINUTE_TO_TICK:  03,
		SECOND_TO_TICK:  0,
	}
	jobTicker.DoPeriodicTask(
		func() {
			registeredEmails := service.registeredMailService.GetAllRegisteredMails()
			mails := []string{}
			for _, email := range registeredEmails {
				mails = append(mails, email.Email)
			}
			service.SendReport(time.Now().UnixMilli()-time.Duration(24*time.Hour).Milliseconds(), time.Now().UnixMilli(), mails)
		},
	)
}

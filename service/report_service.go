package service

import (
	"encoding/json"
	"fmt"
	"time"
	"vcs-sms/model/entity"
	"vcs-sms/util"
)

type ReportService struct {
	esService             *ESService
	mailService           *MailService
	registeredMailService *RegisteredMailService
	serverService         *ServerService
	cacheService          *CacheService
}

func NewReportService(esService *ESService, mailService *MailService, registeredMailService *RegisteredMailService, serverService *ServerService, cacheService *CacheService) *ReportService {
	return &ReportService{
		esService:             esService,
		mailService:           mailService,
		registeredMailService: registeredMailService,
		serverService:         serverService,
		cacheService:          cacheService,
	}
}

type serverUptimeInfo struct {
	Name   string
	Ipv4   string
	Uptime float64
}

func (s serverUptimeInfo) ToString() string {
	return fmt.Sprintf("Name: %s\nIPv4: %s\nUptime: %.2f\n", s.Name, s.Ipv4, s.Uptime)
}

func (service *ReportService) SendReport(startMils int64, endMils int64, to []string) error {
	uptimeInfo := service.esService.CalculateUptime(startMils, endMils)

	servers := []entity.Server{}
	serversString, err := service.cacheService.Get("server:all")
	if err != nil || serversString == "" {
		servers = service.serverService.GetAllServers()
		err = service.cacheService.Set("server:all", servers)
		if err != nil {
			return err
		}
	} else {
		err = json.Unmarshal([]byte(serversString), &servers)
		if err != nil {
			return err
		}
	}

	serversUptimeInfo := []serverUptimeInfo{}
	for _, server := range uptimeInfo {
		for _, s := range servers {
			if server.ID == s.ID {
				serversUptimeInfo = append(serversUptimeInfo, serverUptimeInfo{
					Name:   s.Name,
					Ipv4:   s.IPv4,
					Uptime: server.Uptime.Value,
				})
			}
		}
	}

	onlineCount := 0
	for _, server := range servers {
		if server.Status == 1 {
			onlineCount++
		}
	}

	mailBody := ""
	for _, server := range serversUptimeInfo {
		mailBody += server.ToString()
	}

	if err != nil {
		return err
	}
	mailBody += fmt.Sprintf("\nTotal servers: %d", len(servers))
	mailBody += fmt.Sprintf("\nOnline servers: %d", onlineCount)
	mailBody += fmt.Sprintf("\nOffline servers: %d", (len(servers) - onlineCount))
	err = service.mailService.SendEmail(to, string(mailBody))
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

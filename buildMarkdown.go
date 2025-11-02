package main

import (
	"fmt"
	"strings"
	"time"
)

var mailboxTypes = map[string]string{
	"user":               "Benutzer",
	"linked":             "Verknüpft",
	"shared":             "Gemeinsam",
	"room":               "Raum",
	"equipment":          "Ausrüstung",
	"others":             "Andere",
	"unknownFutureValue": "UnbekannterZukunftswert",
}

func buildMarkdown(users []User, token string) string {
	licenseMap, _ := getLicenseMap(token)
	usageData, _ := getMailboxReport(token)

	header := fmt.Sprintf("**Letzte Aktualisierung:** %s\n\n", time.Now().Format("02 Jan 2006 15:04 MST"))
	tableHeader := "| Anzeigename | E-Mail | Typ | Lizenzen | Weiterleitung | Mail auch ins Postfach | Vollzugriff | Senden als | Senden im Auftrag | Größe | letzte Aktivität |\n"
	separator := "|-------------|--------|-----|----------|---------------|------------------------|-------------|------------|-------------------|-------|------------------|\n"

	var sb strings.Builder
	sb.WriteString(header)
	sb.WriteString(tableHeader)
	sb.WriteString(separator)

	for _, u := range users {
		if u.UserType != "Member" || len(u.Mail) == 0 {
			continue
		}
		ms, _ := getMailboxSettings(token, u.Id)
		fullAccess, sendAs, sendOnBehalf, _ := getDelegation(token, u.Id)

		licenseNames := formatLicenseNames(u.AssignedLicenses, licenseMap)
		forwardingStatus := formatForwardingStatus(ms)

		row := fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			u.DisplayName,
			u.Mail,
			mailboxTypes[ms.UserPurpose],
			licenseNames,
			ms.ForwardingAddress,
			forwardingStatus,
			fullAccess,
			sendAs,
			sendOnBehalf,
			formatBytes(usageData[u.UserPrincipalName].storageUsed),
			usageData[u.UserPrincipalName].lastActivity,
		)
		sb.WriteString(row)
	}
	return sb.String()
}

func formatLicenseNames(assigned []License, licenseMap map[string]string) string {
	if len(assigned) == 0 {
		return "-"
	}
	names := make([]string, 0, len(assigned))
	for _, l := range assigned {
		sku := strings.ToUpper(l.SkuId)
		if friendly, ok := licenseMap[sku]; ok && friendly != "" {
			if pretty := prettyNames[friendly]; pretty != "" {
				friendly = pretty
			}
			names = append(names, friendly)
		} else {
			names = append(names, l.SkuId)
		}
	}
	return strings.Join(names, ", ")
}

func formatForwardingStatus(ms MailboxSettings) string {
	if ms.ForwardingAddress == "" {
		return "-"
	}
	if ms.IsForwarding {
		return "ja"
	}
	return "nein"
}

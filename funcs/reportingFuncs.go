package funcs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reportingFuncs/entity"
	"reportingFuncs/popConstants"
	"strings"

	"strconv"

	"bufio"
	"os"
	"syscall"

	"github.com/john-k-ge/oauth2"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	resultsPerPage = 50
)

var (
	orgApiPath = "/v2/organizations?order-direction=asc&page=%d&results-per-page=" + strconv.Itoa(resultsPerPage)
	orgMemUtil = "/v2/organizations/%v/memory_usage"
)

type ReportingHelper struct {
	uaaConfig   *oauth2.Config
	sharedToken *oauth2.Token
	pop         *popConstants.PoP
}

func ReportingHelperFactory(popName string) (*ReportingHelper, error) {
	popToUse, ok := popConstants.PoPs[popName]
	if !ok {
		log.Printf("Invalid PoP name: %v", popName)
		return nil, errors.New("Invalid PoP name: " + popName)
	}

	rh := &ReportingHelper{
		uaaConfig: &oauth2.Config{
			Endpoint: oauth2.Endpoint{
				AuthURL:  popToUse.Uaa + "/oauth/authorize",
				TokenURL: popToUse.Uaa + "/oauth/token",
			},
			ClientID: "cf",
		},
		pop: popToUse,
	}

	return rh, nil
}

func (rh *ReportingHelper) getCFCreds() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Email: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Println("\nPassword typed: " + string(bytePassword))
	}
	fmt.Println()

	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password)
}

func (rh *ReportingHelper) getAuthCode() string {
	fmt.Printf("You'll need to get an authcode from: %v\n", rh.pop.Passcode)
	fmt.Print("Enter it here: ")
	byteCode, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Println("\nPasscode typed: " + string(byteCode))
	}
	fmt.Println()

	code := string(byteCode)
	return strings.TrimSpace(code)
}

func (rh *ReportingHelper) passwordLogon(uid, pass string) (*oauth2.Token, error) {
	return rh.uaaConfig.PasswordCredentialsToken(context.Background(), uid, pass)
}

func (rh *ReportingHelper) passcodeLogon(code string) (*oauth2.Token, error) {
	return rh.uaaConfig.PasscodeCredentialsToken(context.Background(), code)
}

func (rh *ReportingHelper) Authenticate() error {
	var err error
	switch len(rh.pop.Passcode) {
	case 0:
		uid, pass := rh.getCFCreds()
		rh.sharedToken, err = rh.passwordLogon(uid, pass)
	default:
		code := rh.getAuthCode()
		rh.sharedToken, err = rh.passcodeLogon(code)
	}

	if err != nil {
		log.Printf("Could not logon to `%v` CF!  %v\n", rh.pop.Name, err.Error())
		return err
	}

	return nil
}

func (rh *ReportingHelper) genCFHttpF() func(requestUrl *url.URL) []byte {
	if rh.sharedToken == nil {
		log.Printf("No saved token found.  Is the user logged in?")
		panic("No saved token found.  Is the user logged in?")
	}
	myClient := rh.uaaConfig.Client(context.Background(), rh.sharedToken)

	return func(requestUrl *url.URL) []byte {
		cfRequest, err := http.NewRequest(http.MethodGet, requestUrl.String(), strings.NewReader(""))
		if err != nil {
			log.Printf("Could not build request: %v", err)
			panic("Could not build request: " + err.Error())
		}

		cfRequest.Header.Add("Host", requestUrl.Host)
		response, err := myClient.Do(cfRequest)
		if err != nil {
			log.Printf("Oh No! %v", err)
			panic("Couldn't make the request!\n")
		}
		defer response.Body.Close()
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic("couldn't parse body")
		}

		return content
	}
}

func (rh *ReportingHelper) GetPageCount() (int, error) {
	pageCountUrl, err := url.Parse(fmt.Sprintf(rh.pop.Api+orgApiPath, 1))
	if err != nil {
		log.Printf("Failed to parse org page URL `%v`: %v", rh.pop.Api+orgApiPath, err)
		panic("Failed to parse org page URL " + err.Error())
	}

	cfHttpCall := rh.genCFHttpF()

	content := cfHttpCall(pageCountUrl)

	var batch entity.OrgResponse

	err = json.Unmarshal(content, &batch)

	if err != nil {
		log.Printf("Failed to unmarshall page count response: %v", err)
	}

	return batch.Total_Pages, err
}

// Using functions to minimize global variables
func (rh *ReportingHelper) GenOrgPageF() func(*url.URL) []*entity.OrgInfo {
	cfHttpCall := rh.genCFHttpF()
	return func(url *url.URL) []*entity.OrgInfo {
		content := cfHttpCall(url)

		var batch entity.OrgResponse

		err := json.Unmarshal(content, &batch)

		if err != nil {
			log.Printf("oh-oh: %v", err)
			panic("couldn't unmarshall response!!")
		}

		var orgs []*entity.OrgInfo

		for _, org := range batch.Resources {
			temp := &entity.OrgInfo{
				Name:         org.Entity.Name,
				Guid:         org.Metadata.Guid,
				Status:       org.Entity.Status,
				Quota_url:    org.Entity.Quota_definition_url,
				Quota_guid:   org.Entity.Quota_definition_guid,
				Spaces_url:   org.Entity.Spaces_url,
				Managers_url: org.Entity.Managers_url,
				Users_url:    org.Entity.Users_url,
				Managers:     make(map[string]string),
				Users:        make(map[string]string),
				Created:      org.Metadata.Created_at,
			}
			orgs = append(orgs, temp)
		}

		return orgs
	}
}

func (rh *ReportingHelper) GenMemUtilF() func(string) int {
	cfHttpCall := rh.genCFHttpF()
	return func(guid string) int {
		memUtilUrl, err := url.Parse(rh.pop.Api + fmt.Sprintf(orgMemUtil, guid))
		if err != nil {
			log.Printf("Failed to parse memutil url `%v`: %v", rh.pop.Api+fmt.Sprintf(orgMemUtil, guid), err)
			panic("Failed to parse memutil url: " + err.Error())
		}

		content := cfHttpCall(memUtilUrl)
		var memResponse entity.OrgMemResponse

		err = json.Unmarshal(content, &memResponse)
		if err != nil {
			log.Printf("oh-oh: %v", err)
			panic("couldn't unmarshall response!!")
		}
		return memResponse.Memory_usage_in_mb
	}
}

func (rh *ReportingHelper) GenMemQuotaF() func(*entity.OrgInfo) (string, int) {

	quotaNameCache := make(map[string]string)
	quotaMemValCache := make(map[string]int)

	cfHttpCall := rh.genCFHttpF()

	return func(myOrg *entity.OrgInfo) (string, int) {
		quotaName, foundName := quotaNameCache[myOrg.Quota_guid]
		quotaLimit, foundVal := quotaMemValCache[myOrg.Quota_guid]
		if foundName && foundVal {
			return quotaName, quotaLimit
		}

		quotaReq, err := url.Parse(rh.pop.Api + myOrg.Quota_url)
		if err != nil {
			log.Printf("Failed to parse memquota url `%v`: %v", rh.pop.Api+myOrg.Quota_url, err)
			panic("Failed to parse memquota url: " + err.Error())
		}

		content := cfHttpCall(quotaReq)

		var quotaResponse entity.OrgQuota

		err = json.Unmarshal(content, &quotaResponse)

		if err != nil {
			log.Printf("oh-oh: %v", err)
			panic("couldn't unmarshall quota response!!")
		}
		quotaNameCache[myOrg.Quota_guid] = quotaResponse.Entity.Name
		quotaMemValCache[myOrg.Quota_guid] = quotaResponse.Entity.Memory_limit

		return quotaResponse.Entity.Name, quotaResponse.Entity.Memory_limit
	}
}

func (rh *ReportingHelper) GenOrgUserF() func(info *entity.OrgInfo) map[string]string {
	cfHttpCall := rh.genCFHttpF()
	return func(org *entity.OrgInfo) map[string]string {
		managerUrl, err := url.Parse(rh.pop.Api + org.Managers_url)

		managerContent := cfHttpCall(managerUrl)
		var managersResponse entity.OrgUserResponse

		err = json.Unmarshal(managerContent, &managersResponse)
		if err != nil {
			log.Printf("oh-oh: %v", err)
			panic("couldn't unmarshall response!!")
		}

		managerMap := make(map[string]string)
		for _, manager := range managersResponse.Resources {
			managerMap[manager.Metadata.Guid] = manager.Entity.Username
		}
		return managerMap
	}
}

func (rh *ReportingHelper) GenOrgUrlListF() func(int) []*url.URL {
	var urls []*url.URL

	return func(max int) []*url.URL {
		for i := 1; i <= max; i++ {
			orgPageUrl, err := url.Parse(fmt.Sprintf(rh.pop.Api+orgApiPath, i))
			if err != nil {
				log.Printf("Failed to parse URL `%v`: %v", fmt.Sprintf(rh.pop.Api+orgApiPath, i), err)
				panic("failed to parse org page url: " + err.Error())
			}
			urls = append(urls, orgPageUrl)
		}
		return urls
	}
}

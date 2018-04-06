package funcs

import (
	"log"
	"net/url"
	"reportingFuncs/entity"
	"reportingFuncs/popConstants"
	"testing"
)

const (
	ffApiUrl = "https://api.system.aws-eu-central-1-pr.ice.predix.io"
	//oneOrg       = "/v2/organizations?order-direction=asc&page=1&results-per-page=1"
	fiftyOrgs    = "/v2/organizations?order-direction=asc&page=1&results-per-page=50"
	knownGoodPoP = "ff"
	ffPasscode   = "a2LcB7ah3r"
)

var sharedHelper *ReportingHelper

func initHelper() {
	if sharedHelper == nil {
		log.Print("Initializing sharedHelper...")
		sharedHelper, _ := ReportingHelperFactory(knownGoodPoP)
		var err error
		sharedHelper.sharedToken, err = sharedHelper.passcodeLogon(ffPasscode)
		if err != nil {
			log.Printf("Failed to get token: %v", err)
		}
		if sharedHelper.sharedToken == nil {
			log.Print("Shared token is nil!")
		}

		log.Printf(sharedHelper.sharedToken.AccessToken)
	}
}

func TestReportingHelperFactory(t *testing.T) {
	goodPops := []string{
		"us-w",
		"us-e",
		"ff",
		"jp",
		"cf3",
		"az",
	}

	_, err := ReportingHelperFactory("shoop")
	if err == nil {
		log.Printf("This should have failed, as %v is a bogus PoP", "shoop")
		t.Fail()
	}

	for _, pop := range goodPops {
		test, err := ReportingHelperFactory(pop)
		if err != nil {
			log.Printf("Failed to create ReportingHelper for %v: %v", pop, err)
			t.Fail()
		}
		if test.pop != popConstants.PoPs[pop] {
			log.Printf("PoP initialization failed for %v", pop)
			t.Fail()
		}
	}
}

func Test_genCFHttpF_NilToken(t *testing.T) {
	tester, _ := ReportingHelperFactory(knownGoodPoP)
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("genCFHttpF should have panicked without a token!")
			}
		}()
		// This function should cause a panic
		_ = tester.genCFHttpF()
	}()
}

//func TestReportingHelper_passwordLogin(t *testing.T) {
//
//}

func Test_genCFHttpF(t *testing.T) {
	initHelper()
	testCfHttp := sharedHelper.genCFHttpF()
	testUrl, err := url.Parse(sharedHelper.pop.Api + "/info")
	if err != nil {
		log.Printf("failed to parse url `%v`: %v", sharedHelper.pop.Api+"/info", err)
		t.Fail()
	}

	t.Run("makecfcall", func(t *testing.T) {
		resp := testCfHttp(testUrl)
		if len(resp) == 0 {
			log.Print("failed to get a response")
			t.Fail()
		}
		log.Printf("response: %s", resp)
	})
}

func Test_GetPageCount(t *testing.T) {
	initHelper()
	pageCount, err := sharedHelper.GetPageCount()
	if err != nil {
		log.Printf("Failed to get page count: %v", err)
		t.Fail()
	}
	if pageCount != 4 {
		log.Printf("Page count should be 4, but is `%v`", pageCount)
		t.Fail()
	}
}

func TestReportingHelper_GenOrgUrlListF(t *testing.T) {
	initHelper()
	orgUrlGen := sharedHelper.GenOrgUrlListF()
	onePage := orgUrlGen(1)
	twoPages := orgUrlGen(2)
	if len(onePage) != 1 || len(twoPages) != 2 {
		log.Printf("Not enough URLs generated")
		t.Fail()
	}
	if onePage[0].String() != ffApiUrl+fiftyOrgs {
		log.Print("URL isn't correct.")
		log.Printf("ideal : %v", ffApiUrl+fiftyOrgs)
		log.Printf("actual: %v", onePage[0])
		t.Fail()
	}
}

func TestReportingHelper_GenMemUtilF(t *testing.T) {
	initHelper()
	testFunc := sharedHelper.GenMemUtilF()
	predixSupportFf := "6b132e42-295b-4ef2-9703-c37332ac6dbc"
	mem := testFunc(predixSupportFf)
	if mem != 6688 {
		log.Printf("should be 6688, but memUtil returned is `%v`", mem)
		t.Fail()
	}
	log.Printf("memUtil returned: `%v`", mem)
}

func TestReportingHelper_GenMemQuotaF(t *testing.T) {
	initHelper()
	testFunc := sharedHelper.GenMemQuotaF()
	testOrg := &entity.OrgInfo{
		Quota_guid: "e9fd1013-8b58-490d-bd41-fd88c0266370",
		Quota_url:  "/v2/quota_definitions/e9fd1013-8b58-490d-bd41-fd88c0266370",
	}
	name, val := testFunc(testOrg)
	if name != "120GB" {
		log.Printf("Name should be `120GB`, but is: %v", name)
		t.Fail()
	}
	if val != 122880 {
		log.Printf("Val should be 122880, but is: %v", val)
		t.Fail()
	}
	log.Printf("name: %v, val: %v", name, val)
}

func TestReportingHelper_GenOrgUserF(t *testing.T) {
	initHelper()
	testFunc := sharedHelper.GenOrgUserF()
	testOrg := &entity.OrgInfo{
		Managers_url: "/v2/organizations/6b132e42-295b-4ef2-9703-c37332ac6dbc/managers",
	}

	managers := testFunc(testOrg)
	if len(managers) != 3 {
		log.Printf("expecting 3 managers, but only `%v` returned", len(managers))
		t.Fail()
	}
}

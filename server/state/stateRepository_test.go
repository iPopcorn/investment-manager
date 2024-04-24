package state

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/iPopcorn/investment-manager/types"
	"github.com/iPopcorn/investment-manager/util"
)

func TestStateRepository(t *testing.T) {
	// Arrange
	pathToTestFile, err := util.GetPathToFile("/server/state", "test-state.json")

	if err != nil {
		t.Fatalf("Failed to get path to file\n%v\n", err)
	}

	data, err := os.ReadFile(pathToTestFile)

	if err != nil {
		t.Fatalf("Failed to read file\n%v\n", err)
	}

	var expectedState types.State

	err = json.Unmarshal(data, &expectedState)

	if err != nil {
		t.Fatalf("Failed to de-serialize expected state")
	}

	testRepo := StateRepositoryFactory("test-save-state.json")

	// Act
	err = testRepo.Save(expectedState)

	if err != nil {
		t.Fatalf("Failed to save expected state")
	}

	actualState, err := testRepo.GetState()

	// Assert
	if err != nil {
		t.Fatalf("Failed to get state\n%v\n", err)
	}

	if actualState == nil {
		t.Fatalf("Repo did not get state\n")
	}

	AssertStateEqual(&expectedState, actualState, t)

	// Clean up
	pathToCreatedFile, err := util.GetPathToFile("/server/state", "test-save-state.json")

	err = os.Remove(pathToCreatedFile)
	if err != nil {
		t.Fatalf("Failed to clean up")
	}

}

func AssertStateEqual(expected, actual *types.State, t *testing.T) {
	t.Helper()

	if expected.LastUpdated != actual.LastUpdated {
		t.Errorf("Last Updated does not match\nExpected: %s\nActual: %s\n", expected.LastUpdated, actual.LastUpdated)
	}

	if len(expected.Portfolios) != len(actual.Portfolios) {
		t.Errorf("Number of portfolios in state does not match\nExpected %d\nActual %d\n", len(expected.Portfolios), len(actual.Portfolios))
	}

	if len(actual.Portfolios) != 1 {
		t.Errorf("Incorrect number of portfolios\nExpected 1\nActual %d\n", len(actual.Portfolios))
	}

	actualPortfolio := actual.Portfolios[0]
	expectedPortfolio := expected.Portfolios[0]

	if &actualPortfolio == nil {
		t.Fatalf("Actual portfolio is nil\n")
	}

	AssertStringEqual(expectedPortfolio.Name, actualPortfolio.Name, t)
	AssertStringEqual(expectedPortfolio.Uuid, actualPortfolio.Uuid, t)
	AssertStringEqual(expectedPortfolio.Type, actualPortfolio.Type, t)

	if actualPortfolio.Deleted != expectedPortfolio.Deleted {
		t.Errorf("Deleted is %t\nExpected %t\n", actualPortfolio.Deleted, expectedPortfolio.Deleted)
	}

	AssertStringEqual(expectedPortfolio.CurrentStrategy.Name, actualPortfolio.CurrentStrategy.Name, t)
	AssertStringEqual(expectedPortfolio.CurrentStrategy.Currency, actualPortfolio.CurrentStrategy.Currency, t)

	actualOpenOffersLength := len(actualPortfolio.CurrentStrategy.OpenOffers)
	expectedOpenOfferesLength := len(expectedPortfolio.CurrentStrategy.OpenOffers)

	if actualOpenOffersLength != 1 {
		t.Errorf("Incorrect number of open offers in current strategy\nExpected 1\nActual %d\n", actualOpenOffersLength)
	}

	if expectedOpenOfferesLength != actualOpenOffersLength {
		t.Errorf("Number of open offers does not match\nExpected %d\nActual %d\n", expectedOpenOfferesLength, actualOpenOffersLength)
	}

	actualOpenOffer := actualPortfolio.CurrentStrategy.OpenOffers[0]
	expectedOpenOffer := expectedPortfolio.CurrentStrategy.OpenOffers[0]

	AssertStringEqual(expectedOpenOffer.ClientOrderId, actualOpenOffer.ClientOrderId, t)
	AssertStringEqual(expectedOpenOffer.ProductId, actualOpenOffer.ProductId, t)
	AssertStringEqual(string(expectedOpenOffer.Side), string(actualOpenOffer.Side), t)
	AssertStringEqual(expectedOpenOffer.SelfTradePreventionId, actualOpenOffer.SelfTradePreventionId, t)
	AssertStringEqual(expectedOpenOffer.RetailPortfolioId, actualOpenOffer.RetailPortfolioId, t)

	if !reflect.DeepEqual(expectedOpenOffer.Config, actualOpenOffer.Config) {
		t.Errorf("Open offer config does not match\nExpected: %+v\nActual: %+v\n", expectedOpenOffer.Config, actualOpenOffer.Config)
	}
}

func AssertStringEqual(expected, actual string, t *testing.T) {
	t.Helper()

	if expected != actual {
		t.Errorf("Strings do not match\nExpected: %q\nActual: %q\n", expected, actual)
	}
}

package service

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testFunc func(context.Context, db.ConnectionPool, repositories.Repositories) error
type returnTestFunc func(context.Context, db.ConnectionPool, repositories.Repositories) interface{}
type generateRepositoriesMock func() repositories.Repositories

type verifyError func(error, *require.Assertions)
type verifyMockInteractions func(repositories.Repositories, *require.Assertions)

type repositoryInteractionTestCase struct {
	generateRepositoriesMock generateRepositoriesMock
	handler                  testFunc
	expectedError            error
	verifyError              verifyError
	verifyInteractions       verifyMockInteractions
}

type verifyContent func(interface{}, repositories.Repositories, *require.Assertions)

type returnTestCase struct {
	generateRepositoriesMock generateRepositoriesMock
	handler                  returnTestFunc
	expectedContent          interface{}
	verifyContent            verifyContent
}

type transactionTestCase struct {
	generateRepositoriesMock generateRepositoriesMock
	handler                  testFunc
}

type ServiceTestSuite struct {
	suite.Suite

	generateRepositoriesMock      generateRepositoriesMock
	generateErrorRepositoriesMock generateRepositoriesMock

	repositoryInteractionTestCases map[string]repositoryInteractionTestCase
	returnTestCases                map[string]returnTestCase
	transactionTestCases           map[string]transactionTestCase
}

func (s *ServiceTestSuite) TestWhenCallingHandler_ExpectCorrectInteraction() {
	for name, testCase := range s.repositoryInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			var repos repositories.Repositories
			if testCase.generateRepositoriesMock != nil {
				repos = testCase.generateRepositoriesMock()
			} else {
				repos = s.generateRepositoriesMock()
			}

			err := testCase.handler(context.Background(), &mockConnectionPool{}, repos)

			if testCase.verifyError != nil {
				testCase.verifyError(err, s.Require())
			} else {
				s.Require().Equal(testCase.expectedError, err)
			}
			if testCase.verifyInteractions != nil {
				testCase.verifyInteractions(repos, s.Require())
			}
		})
	}
}

func (s *ServiceTestSuite) TestWhenRepositorySucceeds_ReturnsExpectedValue() {
	for name, testCase := range s.returnTestCases {
		s.T().Run(name, func(t *testing.T) {
			var repos repositories.Repositories
			if testCase.generateRepositoriesMock != nil {
				repos = testCase.generateRepositoriesMock()
			} else {
				repos = s.generateRepositoriesMock()
			}

			actual := testCase.handler(context.Background(), &mockConnectionPool{}, repos)

			if testCase.verifyContent != nil {
				testCase.verifyContent(actual, repos, s.Require())
			} else {
				s.Require().Equal(testCase.expectedContent, actual)
			}
		})
	}
}

func (s *ServiceTestSuite) TestWhenUsingTransaction_ExpectCallsClose() {
	for name, testCase := range s.transactionTestCases {
		s.T().Run(name, func(t *testing.T) {
			var repos repositories.Repositories
			if testCase.generateRepositoriesMock != nil {
				repos = testCase.generateRepositoriesMock()
			} else {
				repos = s.generateRepositoriesMock()
			}

			m := &mockConnectionPool{}
			testCase.handler(context.Background(), m, repos)

			s.Require().Equal(1, m.tx.closeCalled)
		})
	}
}

func (s *ServiceTestSuite) TestWhenCreatingTransactionFails_ExpectErrorIsPropagated() {
	for name, testCase := range s.transactionTestCases {
		s.T().Run(name, func(t *testing.T) {
			var repos repositories.Repositories
			if testCase.generateRepositoriesMock != nil {
				repos = testCase.generateRepositoriesMock()
			} else {
				repos = s.generateRepositoriesMock()
			}

			m := &mockConnectionPool{
				err: errDefault,
			}
			err := testCase.handler(context.Background(), m, repos)

			s.Require().Equal(errDefault, err)
		})
	}
}

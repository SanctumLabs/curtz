package urlsvc

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
	"github.com/sanctumlabs/curtz/app/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUrlWriteSvc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UrlWriteSvc Suite")
}

var _ = Describe("UrlWriteSvc", func() {
	var (
		mockCtrl         *gomock.Controller
		mockUrlWriteRepo *mocks.MockUrlWriteRepository
		mockUserSvc      *mocks.MockUserService
	)

	urlWriteSvc := &UrlWriteSvc{}

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUrlWriteRepo = mocks.NewMockUrlWriteRepository(mockCtrl)
		mockUserSvc = mocks.NewMockUserService(mockCtrl)
		urlWriteSvc.repo = mockUrlWriteRepo
		urlWriteSvc.userSvc = mockUserSvc
	})

	When("Removing/Deleting a url with a given a ID", func() {
		It("should return an error when there is a failure deleting a url", func() {
			urlID := identifier.New().String()
			expectedErr := errors.New("Failed to delete url by id")

			mockUrlWriteRepo.
				EXPECT().
				Delete(urlID).
				Return(expectedErr)

			actual := urlWriteSvc.Remove(urlID)

			if assert.Error(GinkgoT(), actual) {
				assert.Equal(GinkgoT(), expectedErr, actual)
			}
		})

		It("should not return an error when there is a success deleting a url", func() {
			urlID := identifier.New().String()

			mockUrlWriteRepo.
				EXPECT().
				Delete(urlID).
				Return(nil)

			actual := urlWriteSvc.Remove(urlID)

			assert.NoError(GinkgoT(), actual)
		})
	})
})

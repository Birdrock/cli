package v6_test

import (
	"code.cloudfoundry.org/cli/api/uaa"
	"code.cloudfoundry.org/cli/command/commandfakes"
	. "code.cloudfoundry.org/cli/command/v6"
	"code.cloudfoundry.org/cli/command/v6/v6fakes"
	"code.cloudfoundry.org/cli/util/configv3"
	"code.cloudfoundry.org/cli/util/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("logout command", func() {
	var (
		cmd        LogoutCommand
		testUI     *ui.UI
		fakeActor  *v6fakes.FakeAuthActor
		fakeConfig *commandfakes.FakeConfig
		executeErr error
	)

	BeforeEach(func() {
		testUI = ui.NewTestUI(nil, NewBuffer(), NewBuffer())
		fakeConfig = new(commandfakes.FakeConfig)
		fakeActor = new(v6fakes.FakeAuthActor)
		cmd = LogoutCommand{
			UI:     testUI,
			Config: fakeConfig,
			Actor:  fakeActor,
		}
		fakeConfig.CurrentUserReturns(
			configv3.User{
				Name: "some-user",
			},
			nil)
	})

	JustBeforeEach(func() {
		executeErr = cmd.Execute(nil)
	})

	It("outputs logging out display message", func() {
		Expect(executeErr).ToNot(HaveOccurred())

		Expect(fakeConfig.UnsetUserInformationCallCount()).To(Equal(1))
		Expect(testUI.Out).To(Say("Logging out some-user..."))
		Expect(testUI.Out).To(Say("OK"))
	})

	It("calls to revoke the auth tokens", func() {
		Expect(fakeActor.RevokeCallCount()).To(Equal(1))
	})

	When("The revoking of tokens fails", func() {

		BeforeEach(func() {
			fakeActor.RevokeReturns(error(uaa.UnauthorizedError{Message: "test error"}))
		})
		It("does not impact the logout", func() {
			Expect(executeErr).ToNot(HaveOccurred())

			Expect(fakeConfig.UnsetUserInformationCallCount()).To(Equal(1))
			Expect(testUI.Out).To(Say("Logging out some-user..."))
			Expect(testUI.Out).To(Say("OK"))
		})
	})
})

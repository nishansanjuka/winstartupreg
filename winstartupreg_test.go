package winstartupreg_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nishansanjuka/winstartupreg"
)

func TestWindowsStartupRegistry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Windows Startup Registry Suite")
}

var _ = Describe("Windows Startup Registry Management", func() {
	var (
		testAppName string
		testCommand string
	)

	BeforeEach(func() {
		// Generate a unique app name for each test
		testAppName = "TestApp"

		// Create a temporary executable for testing
		tempExe, err := createTempExecutable()
		Expect(err).To(BeNil())
		testCommand = tempExe
	})

	AfterEach(func() {
		// Clean up any potential leftover registry entries
		_ = winstartupreg.RemoveStartupEntry(testAppName, winstartupreg.CurrentUserRun)
	})

	Describe("Adding Startup Entries", func() {
		Context("With valid input", func() {
			It("Should add a startup entry successfully", func() {
				entry := winstartupreg.StartupEntry{
					Name:    testAppName,
					Command: testCommand,
				}

				err := winstartupreg.AddStartupEntry(entry, winstartupreg.CurrentUserRun)
				Expect(err).To(BeNil())

				// Verify the entry was added
				entries, err := winstartupreg.ListStartupEntries(winstartupreg.CurrentUserRun)
				Expect(err).To(BeNil())
				Expect(entries).To(HaveKey(testAppName))
				Expect(entries[testAppName]).To(Equal(testCommand))
			})
		})

		Context("With invalid input", func() {
			It("Should return an error for non-existent executable", func() {
				entry := winstartupreg.StartupEntry{
					Name:    testAppName,
					Command: "/path/to/nonexistent/executable",
				}

				err := winstartupreg.AddStartupEntry(entry, winstartupreg.CurrentUserRun)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("Removing Startup Entries", func() {
		It("Should remove entry from all possible locations", func() {
			// Add test entry
			err := winstartupreg.AddStartupEntry(
				winstartupreg.StartupEntry{
					Name:    testAppName,
					Command: testCommand,
				},
				winstartupreg.CurrentUserRun,
			)
			Expect(err).To(BeNil())

			// Safe remove
			err = winstartupreg.SafeRemoveStartupEntry("TestApp")
			Expect(err).To(BeNil())

			// Verify removal from all locations
			allEntries, err := winstartupreg.ListAllStartupEntries()
			Expect(err).To(BeNil())

			for _, entries := range allEntries {
				Expect(entries).ToNot(HaveKey("TestApp"))
			}
		})
	})

	Describe("Listing Startup Entries", func() {
		Context("When entries exist", func() {
			BeforeEach(func() {
				// Add multiple entries
				entries := []winstartupreg.StartupEntry{
					{
						Name:    testAppName + "_1",
						Command: testCommand,
					},
					{
						Name:    testAppName + "_2",
						Command: testCommand,
					},
				}

				for _, entry := range entries {
					_ = winstartupreg.AddStartupEntry(entry, winstartupreg.CurrentUserRun)
				}
			})

			It("Should list all startup entries", func() {
				entries, err := winstartupreg.ListStartupEntries(winstartupreg.CurrentUserRun)
				Expect(err).To(BeNil())
				Expect(entries).To(HaveKey(testAppName + "_1"))
				Expect(entries).To(HaveKey(testAppName + "_2"))

				for key, value := range entries {
					if strings.Contains(key, "TestApp") {
						fmt.Printf("CREATED REGISTRY FOR Key: %s, Value: %s\n", key, value)
					}
				}

			})

			It(fmt.Sprintf("should remove listed entries %s and %s", testAppName+"_1", testAppName+"_2"), func() {
				// Remove the first entry
				err1 := winstartupreg.SafeRemoveStartupEntry(testAppName + "_1")
				Expect(err1).To(BeNil())

				// Remove the second entry
				err2 := winstartupreg.SafeRemoveStartupEntry(testAppName + "_2")
				Expect(err2).To(BeNil())

				// Verify both entries were removed
				entries, err := winstartupreg.ListStartupEntries(winstartupreg.CurrentUserRun)
				Expect(err).To(BeNil())
				Expect(entries).ToNot(HaveKey(testAppName + "_1"))
				Expect(entries).ToNot(HaveKey(testAppName + "_2"))

				fmt.Printf("REMOVED REGISTRY Key: %s, Value: %s\n", testAppName+"_1", testCommand)
				fmt.Printf("REMOVED REGISTRY Key: %s, Value: %s\n", testAppName+"_2", testCommand)
			})
		})
	})

	Describe("Safe Remove Startup Entry", func() {
		Context("When entry exists", func() {
			BeforeEach(func() {
				// Add an entry before each test
				entry := winstartupreg.StartupEntry{
					Name:    testAppName,
					Command: testCommand,
				}

				_ = winstartupreg.AddStartupEntry(entry, winstartupreg.CurrentUserRun)
			})

			It("Should safely remove the startup entry", func() {
				err1 := winstartupreg.SafeRemoveStartupEntry(testAppName)
				Expect(err1).To(BeNil())

				// Verify the entry was removed
				entries, err := winstartupreg.ListStartupEntries(winstartupreg.CurrentUserRun)
				Expect(err).To(BeNil())
				Expect(entries).ToNot(HaveKey(testAppName))
			})
		})

		Context("When entry does not exist", func() {
			It("Should return an error", func() {
				err := winstartupreg.SafeRemoveStartupEntry("TestApp")
				Expect(err).To(HaveOccurred())
			})
		})
	})

})

// Create a temporary executable for testing
func createTempExecutable() (string, error) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "winstartupreg-test")
	if err != nil {
		return "", err
	}

	// Create a temporary executable
	tempExe := filepath.Join(tempDir, "testapp.exe")

	// Create a dummy executable (on Windows)
	file, err := os.Create(tempExe)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write minimal executable content
	_, err = file.Write([]byte{0x4D, 0x5A}) // Minimal EXE header
	if err != nil {
		return "", err
	}

	return tempExe, nil
}

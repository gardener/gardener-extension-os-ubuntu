package validation

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1alpha1 "github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/config/v1alpha1"
)

var _ = Describe("ExtensionConfig validation", func() {
	var (
		config *configv1alpha1.ExtensionConfig
	)

	BeforeEach(func() {
		config = &configv1alpha1.ExtensionConfig{
			NTP: &configv1alpha1.NTPConfig{
				Daemon: configv1alpha1.SystemdTimesyncd,
			},
		}
	})

	It("should allow valid config", func() {
		Expect(ValidateExtensionConfig(config)).To(BeEmpty())
	})

	It("should fail with incorrect daemon name", func() {
		config.NTP.Daemon = "foo"
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeNotSupported))
		Expect(errs[0].Field).To(Equal("daemon"))
	})

	It("should succeed with valid NTPd config", func() {
		config.NTP.Daemon = configv1alpha1.NTPD
		config.NTP.NTPD = &configv1alpha1.NTPDConfig{}
		config.NTP.NTPD.Servers = []string{"ntp.ubuntu.com"}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(BeEmpty())
	})

	It("should fail with invalid NTPd config (no ntp servers provided)", func() {
		config.NTP.Daemon = configv1alpha1.NTPD
		config.NTP.NTPD = &configv1alpha1.NTPDConfig{}
		config.NTP.NTPD.Servers = []string{}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeRequired))
		Expect(errs[0].Field).To(Equal("ntpd.servers"))
	})

	It("should fail with daemon systemd-timesyncd and ntpd config set", func() {
		config.NTP.Daemon = configv1alpha1.SystemdTimesyncd
		config.NTP.NTPD = &configv1alpha1.NTPDConfig{Servers: []string{"foo.bar"}}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeForbidden))
		Expect(errs[0].Field).To(Equal("ntpd"))
	})

	It("should succeed with a valid architecture, URI and search for primary apt mirror", func() {
		config.APTConfig = &configv1alpha1.APTConfig{Primary: []configv1alpha1.APTArchive{{
			Arches:    []configv1alpha1.Architecture{configv1alpha1.ARM64},
			URI:       "http://packages.ubuntu-mirror.example.com/apt-mirror/ubuntu",
			Search:    []string{"http://archive.ubuntu.com/ubuntu/"},
			SearchDNS: false,
		}}}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(BeEmpty())
	})

	It("should fail with an invalid URI and invalid search for primary apt mirror", func() {
		config.APTConfig = &configv1alpha1.APTConfig{Primary: []configv1alpha1.APTArchive{{
			Arches:    []configv1alpha1.Architecture{configv1alpha1.ARM64},
			URI:       "packages.ubuntu-mirror.example.com/apt-mirror/ubuntu",
			Search:    []string{"archive.ubuntu.com/ubuntu/"},
			SearchDNS: false,
		}}}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(2))
	})
})

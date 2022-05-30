package aws

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/recode-sh/aws-cloud-provider/userconfig"
	"github.com/recode-sh/cli/internal/mocks"
)

func TestUserConfigLocalResolving(t *testing.T) {
	testCases := []struct {
		test                            string
		configInEnvVars                 *userconfig.Config
		errorDuringEnvVarsResolving     error
		configInFiles                   *userconfig.Config
		errorDuringConfigFilesResolving error
		profileOpts                     string
		expectedConfig                  *userconfig.Config
		expectedError                   error
	}{
		{
			test:                            "no env vars, no config files",
			errorDuringEnvVarsResolving:     userconfig.ErrMissingConfig,
			errorDuringConfigFilesResolving: userconfig.ErrMissingConfig,
			expectedConfig:                  nil,
			expectedError:                   userconfig.ErrMissingConfig,
		},

		{
			test:                            "only env vars",
			configInEnvVars:                 userconfig.NewConfig("a", "b", "c"),
			errorDuringConfigFilesResolving: userconfig.ErrMissingConfig,
			expectedConfig:                  userconfig.NewConfig("a", "b", "c"),
			expectedError:                   nil,
		},

		{
			test:                        "only config files",
			errorDuringEnvVarsResolving: userconfig.ErrMissingConfig,
			configInFiles:               userconfig.NewConfig("a", "b", "c"),
			expectedConfig:              userconfig.NewConfig("a", "b", "c"),
			expectedError:               nil,
		},

		{
			test:            "env vars and config files",
			configInEnvVars: userconfig.NewConfig("a", "b", "c"),
			configInFiles:   userconfig.NewConfig("d", "e", "f"),
			expectedConfig:  userconfig.NewConfig("a", "b", "c"),
			expectedError:   nil,
		},

		{
			test:            "env vars, config files and profile",
			configInEnvVars: userconfig.NewConfig("a", "b", "c"),
			configInFiles:   userconfig.NewConfig("d", "e", "f"),
			profileOpts:     "production",
			expectedConfig:  userconfig.NewConfig("d", "e", "f"),
			expectedError:   nil,
		},

		{
			test:                        "errored env vars and config files",
			errorDuringEnvVarsResolving: userconfig.ErrMissingAccessKeyInEnv,
			configInFiles:               userconfig.NewConfig("d", "e", "f"),
			expectedConfig:              nil,
			expectedError:               userconfig.ErrMissingAccessKeyInEnv,
		},

		{
			test:                            "env vars and errored config files",
			configInEnvVars:                 userconfig.NewConfig("a", "b", "c"),
			errorDuringConfigFilesResolving: userconfig.ErrMissingRegionInFiles,
			expectedConfig:                  userconfig.NewConfig("a", "b", "c"),
			expectedError:                   nil,
		},

		{
			test:                            "env vars, errored config files and profile",
			configInEnvVars:                 userconfig.NewConfig("a", "b", "c"),
			errorDuringConfigFilesResolving: userconfig.ErrMissingRegionInFiles,
			profileOpts:                     "production",
			expectedConfig:                  nil,
			expectedError:                   userconfig.ErrMissingRegionInFiles,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			userConfigEnvVarsResolverMock := mocks.NewAWSUserConfigEnvVarsResolver(mockCtrl)
			userConfigEnvVarsResolverMock.EXPECT().Resolve().Return(
				tc.configInEnvVars,
				tc.errorDuringEnvVarsResolving,
			).AnyTimes()

			userConfigFilesResolverMock := mocks.NewAWSUserConfigFilesResolver(mockCtrl)
			userConfigFilesResolverMock.EXPECT().Resolve().Return(
				tc.configInFiles,
				tc.errorDuringConfigFilesResolving,
			).AnyTimes()

			resolver := NewUserConfigLocalResolver(
				userConfigEnvVarsResolverMock,
				userConfigFilesResolverMock,
				UserConfigLocalResolverOpts{
					Profile: tc.profileOpts,
				},
			)

			resolvedConfig, err := resolver.Resolve()

			if tc.expectedError == nil && err != nil {
				t.Fatalf("expected no error, got '%+v'", err)
			}

			if tc.expectedError != nil && !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected error to equal '%+v', got '%+v'", tc.expectedError, err)
			}

			if tc.expectedConfig != nil && *resolvedConfig != *tc.expectedConfig {
				t.Fatalf("expected config to equal '%+v', got '%+v'", *tc.expectedConfig, *resolvedConfig)
			}

			if tc.expectedConfig == nil && resolvedConfig != nil {
				t.Fatalf("expected no config, got '%+v'", *resolvedConfig)
			}
		})
	}
}

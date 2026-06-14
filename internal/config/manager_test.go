package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConfig_Validate_MountPaths(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		wantErr     bool
		errContains string
	}{
		{
			name: "mount type fuse - ok",
			config: &Config{
				MountType: MountTypeFuse,
				MountPath: "/mnt/remotes/altmount",
				Metadata: MetadataConfig{
					RootPath: "/metadata",
				},
				WebDAV: WebDAVConfig{
					Port: 8080,
				},
				Streaming: StreamingConfig{
					MaxPrefetch: 30,
				},
				Import: ImportConfig{
					MaxProcessorWorkers:            2,
					QueueProcessingIntervalSeconds: 5,
					MaxImportConnections:           5,
					MaxDownloadPrefetch:            3,
					SegmentSamplePercentage:        1,
					ImportStrategy:                 ImportStrategyNone,
				},
				Health: HealthConfig{
					CheckIntervalSeconds:          5,
					MaxConnectionsForHealthChecks: 5,
					MaxConcurrentJobs:             1,
					SegmentSamplePercentage:       5,
				},
			},
			wantErr: false,
		},
		{
			name: "mount type rclone - ok",
			config: &Config{
				MountType: MountTypeRClone,
				MountPath: "/mnt/remotes/altmount",
				Metadata: MetadataConfig{
					RootPath: "/metadata",
				},
				WebDAV: WebDAVConfig{
					Port: 8080,
				},
				Streaming: StreamingConfig{
					MaxPrefetch: 30,
				},
				Import: ImportConfig{
					MaxProcessorWorkers:            2,
					QueueProcessingIntervalSeconds: 5,
					MaxImportConnections:           5,
					MaxDownloadPrefetch:            3,
					SegmentSamplePercentage:        1,
					ImportStrategy:                 ImportStrategyNone,
				},
				Health: HealthConfig{
					CheckIntervalSeconds:          5,
					MaxConnectionsForHealthChecks: 5,
					MaxConcurrentJobs:             1,
					SegmentSamplePercentage:       5,
				},
			},
			wantErr: false,
		},
		{
			name: "mount type none - ok",
			config: &Config{
				MountType: MountTypeNone,
				Metadata: MetadataConfig{
					RootPath: "/metadata",
				},
				WebDAV: WebDAVConfig{
					Port: 8080,
				},
				Streaming: StreamingConfig{
					MaxPrefetch: 30,
				},
				Import: ImportConfig{
					MaxProcessorWorkers:            2,
					QueueProcessingIntervalSeconds: 5,
					MaxImportConnections:           5,
					MaxDownloadPrefetch:            3,
					SegmentSamplePercentage:        1,
					ImportStrategy:                 ImportStrategyNone,
				},
				Health: HealthConfig{
					CheckIntervalSeconds:          5,
					MaxConnectionsForHealthChecks: 5,
					MaxConcurrentJobs:             1,
					SegmentSamplePercentage:       5,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_GetWebhookBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected string
	}{
		{
			name: "explicitly set",
			config: Config{
				Arrs: ArrsConfig{
					WebhookBaseURL: "http://custom:1234",
				},
				WebDAV: WebDAVConfig{
					Port: 8080,
				},
			},
			expected: "http://custom:1234",
		},
		{
			name: "default with port 8080",
			config: Config{
				Arrs: ArrsConfig{
					WebhookBaseURL: "",
				},
				WebDAV: WebDAVConfig{
					Port: 8080,
				},
			},
			expected: "http://altmount:8080",
		},
		{
			name: "default with port 8084",
			config: Config{
				Arrs: ArrsConfig{
					WebhookBaseURL: "",
				},
				WebDAV: WebDAVConfig{
					Port: 8084,
				},
			},
			expected: "http://altmount:8084",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.config.GetWebhookBaseURL())
		})
	}
}

func TestConfig_GetDownloadClientBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected string
	}{
		{
			name: "explicitly set",
			config: Config{
				SABnzbd: SABnzbdConfig{
					DownloadClientBaseURL: "http://custom:1234/sab",
				},
				WebDAV: WebDAVConfig{
					Port: 8080,
				},
			},
			expected: "http://custom:1234/sab",
		},
		{
			name: "default with port 8080",
			config: Config{
				SABnzbd: SABnzbdConfig{
					DownloadClientBaseURL: "",
				},
				WebDAV: WebDAVConfig{
					Port: 8080,
				},
			},
			expected: "http://altmount:8080/sabnzbd",
		},
		{
			name: "default with port 8084",
			config: Config{
				SABnzbd: SABnzbdConfig{
					DownloadClientBaseURL: "",
				},
				WebDAV: WebDAVConfig{
					Port: 8084,
				},
			},
			expected: "http://altmount:8084/sabnzbd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.config.GetDownloadClientBaseURL())
		})
	}
}

func TestDefaultConfig_DanishEditionSABnzbdReady(t *testing.T) {
	cfg := DefaultConfig("/config")

	assert.NotNil(t, cfg.SABnzbd.Enabled)
	assert.True(t, *cfg.SABnzbd.Enabled)
	assert.NotNil(t, cfg.Arrs.Enabled)
	assert.True(t, *cfg.Arrs.Enabled)
	assert.NotNil(t, cfg.Streaming.FailureMasking.Enabled)
	assert.True(t, *cfg.Streaming.FailureMasking.Enabled)
	assert.Equal(t, 3, cfg.Streaming.FailureMasking.Threshold)
	assert.Equal(t, 2, cfg.Import.MaxConcurrentImports)
	assert.Equal(t, 1, cfg.Import.MaxConcurrentImportsWhileStreaming)
	assert.Equal(t, []string{
		".mkv", ".mp4", ".avi", ".ts", ".m4v", ".mov", ".wmv", ".mpg", ".mpeg",
		".xvid", ".rm", ".rmvb", ".asf", ".asx", ".wtv", ".mk3d", ".dvr-ms",
	}, cfg.Import.AllowedFileExtensions)
	assert.NotNil(t, cfg.Import.AllowNestedRarExtraction)
	assert.True(t, *cfg.Import.AllowNestedRarExtraction)
	assert.NotNil(t, cfg.Import.ExpandBlurayIso)
	assert.True(t, *cfg.Import.ExpandBlurayIso)
	assert.NotNil(t, cfg.Import.RenameToNzbName)
	assert.True(t, *cfg.Import.RenameToNzbName)
	assert.NotNil(t, cfg.Import.FilterSampleFiles)
	assert.True(t, *cfg.Import.FilterSampleFiles)
	assert.NotNil(t, cfg.Import.FailedItemRetentionHours)
	assert.Equal(t, 24, *cfg.Import.FailedItemRetentionHours)
	assert.NotNil(t, cfg.Import.HistoryRetentionDays)
	assert.Equal(t, 30, *cfg.Import.HistoryRetentionDays)
	assert.NotNil(t, cfg.Import.DeleteCompletedNzb)
	assert.False(t, *cfg.Import.DeleteCompletedNzb)
	assert.NotNil(t, cfg.Health.Enabled)
	assert.True(t, *cfg.Health.Enabled)
	assert.NotNil(t, cfg.Health.ResolveRepairOnImport)
	assert.True(t, *cfg.Health.ResolveRepairOnImport)
	assert.Equal(t, 2, cfg.Health.MaxRetries)
	assert.NotNil(t, cfg.Health.Repair.Enabled)
	assert.True(t, *cfg.Health.Repair.Enabled)
	assert.Equal(t, 3, cfg.Health.Repair.MaxRepairRetries)
	assert.NotNil(t, cfg.SegmentCache.Enabled)
	assert.True(t, *cfg.SegmentCache.Enabled)
	assert.Equal(t, "/config/segment-cache", cfg.SegmentCache.CachePath)
	assert.Equal(t, 150, cfg.SegmentCache.MaxSizeGB)
	assert.Equal(t, 72, cfg.SegmentCache.ExpiryHours)
	assert.NotNil(t, cfg.Arrs.QueueCleanupEnabled)
	assert.True(t, *cfg.Arrs.QueueCleanupEnabled)
	assert.Equal(t, 300, cfg.Arrs.QueueCleanupIntervalSeconds)

	assert.Len(t, cfg.SABnzbd.Categories, 5)
	assert.Equal(t, "movies", cfg.SABnzbd.Categories[0].Name)
	assert.Equal(t, "movies", cfg.SABnzbd.Categories[0].Dir)
	assert.Equal(t, "radarr", cfg.SABnzbd.Categories[0].Type)
	assert.Equal(t, "tv", cfg.SABnzbd.Categories[1].Name)
	assert.Equal(t, "tv", cfg.SABnzbd.Categories[1].Dir)
	assert.Equal(t, "sonarr", cfg.SABnzbd.Categories[1].Type)
}

func TestApplyDockerEnvOverrides_DanishHealthControls(t *testing.T) {
	t.Setenv("ALTMOUNT_ENABLE_HEALTH", "false")
	t.Setenv("ALTMOUNT_STREAMING_FAILURE_MASKING_ENABLED", "false")
	t.Setenv("ALTMOUNT_STREAMING_FAILURE_MASKING_THRESHOLD", "7")
	t.Setenv("ALTMOUNT_IMPORT_MAX_CONCURRENT_IMPORTS", "3")
	t.Setenv("ALTMOUNT_IMPORT_MAX_CONCURRENT_IMPORTS_WHILE_STREAMING", "1")
	t.Setenv("ALTMOUNT_IMPORT_ALLOWED_FILE_EXTENSIONS", ".mkv,.mp4")
	t.Setenv("ALTMOUNT_IMPORT_ALLOW_NESTED_RAR_EXTRACTION", "false")
	t.Setenv("ALTMOUNT_IMPORT_EXPAND_BLURAY_ISO", "false")
	t.Setenv("ALTMOUNT_IMPORT_RENAME_TO_NZB_NAME", "false")
	t.Setenv("ALTMOUNT_IMPORT_FILTER_SAMPLE_FILES", "false")
	t.Setenv("ALTMOUNT_IMPORT_FAILED_ITEM_RETENTION_HOURS", "12")
	t.Setenv("ALTMOUNT_IMPORT_HISTORY_RETENTION_DAYS", "14")
	t.Setenv("ALTMOUNT_IMPORT_DELETE_COMPLETED_NZB", "true")
	t.Setenv("ALTMOUNT_HEALTH_RESOLVE_REPAIR_ON_IMPORT", "false")
	t.Setenv("ALTMOUNT_HEALTH_MAX_RETRIES", "4")
	t.Setenv("ALTMOUNT_HEALTH_SEGMENT_SAMPLE_PERCENTAGE", "2")
	t.Setenv("ALTMOUNT_HEALTH_LIBRARY_SYNC_INTERVAL_MINUTES", "720")
	t.Setenv("ALTMOUNT_HEALTH_MAX_REPAIR_RETRIES", "5")
	t.Setenv("ALTMOUNT_SEGMENT_CACHE_ENABLED", "false")
	t.Setenv("ALTMOUNT_SEGMENT_CACHE_PATH", "/tmp/segments")
	t.Setenv("ALTMOUNT_SEGMENT_CACHE_MAX_SIZE_GB", "12")
	t.Setenv("ALTMOUNT_SEGMENT_CACHE_EXPIRY_HOURS", "24")
	t.Setenv("ALTMOUNT_ARRS_QUEUE_CLEANUP_ENABLED", "false")
	t.Setenv("ALTMOUNT_ARRS_QUEUE_CLEANUP_INTERVAL_SECONDS", "120")

	cfg := DefaultConfig()
	err := applyDockerEnvOverrides(cfg)

	assert.NoError(t, err)
	assert.NotNil(t, cfg.Streaming.FailureMasking.Enabled)
	assert.False(t, *cfg.Streaming.FailureMasking.Enabled)
	assert.Equal(t, 7, cfg.Streaming.FailureMasking.Threshold)
	assert.Equal(t, 3, cfg.Import.MaxConcurrentImports)
	assert.Equal(t, 1, cfg.Import.MaxConcurrentImportsWhileStreaming)
	assert.Equal(t, []string{".mkv", ".mp4"}, cfg.Import.AllowedFileExtensions)
	assert.NotNil(t, cfg.Import.AllowNestedRarExtraction)
	assert.False(t, *cfg.Import.AllowNestedRarExtraction)
	assert.NotNil(t, cfg.Import.ExpandBlurayIso)
	assert.False(t, *cfg.Import.ExpandBlurayIso)
	assert.NotNil(t, cfg.Import.RenameToNzbName)
	assert.False(t, *cfg.Import.RenameToNzbName)
	assert.NotNil(t, cfg.Import.FilterSampleFiles)
	assert.False(t, *cfg.Import.FilterSampleFiles)
	assert.NotNil(t, cfg.Import.FailedItemRetentionHours)
	assert.Equal(t, 12, *cfg.Import.FailedItemRetentionHours)
	assert.NotNil(t, cfg.Import.HistoryRetentionDays)
	assert.Equal(t, 14, *cfg.Import.HistoryRetentionDays)
	assert.NotNil(t, cfg.Import.DeleteCompletedNzb)
	assert.True(t, *cfg.Import.DeleteCompletedNzb)
	assert.NotNil(t, cfg.Health.Enabled)
	assert.False(t, *cfg.Health.Enabled)
	assert.NotNil(t, cfg.Health.ResolveRepairOnImport)
	assert.False(t, *cfg.Health.ResolveRepairOnImport)
	assert.Equal(t, 4, cfg.Health.MaxRetries)
	assert.Equal(t, 2, cfg.Health.SegmentSamplePercentage)
	assert.Equal(t, 720, cfg.Health.LibrarySyncIntervalMinutes)
	assert.Equal(t, 5, cfg.Health.Repair.MaxRepairRetries)
	assert.NotNil(t, cfg.SegmentCache.Enabled)
	assert.False(t, *cfg.SegmentCache.Enabled)
	assert.Equal(t, "/tmp/segments", cfg.SegmentCache.CachePath)
	assert.Equal(t, 12, cfg.SegmentCache.MaxSizeGB)
	assert.Equal(t, 24, cfg.SegmentCache.ExpiryHours)
	assert.NotNil(t, cfg.Arrs.QueueCleanupEnabled)
	assert.False(t, *cfg.Arrs.QueueCleanupEnabled)
	assert.Equal(t, 120, cfg.Arrs.QueueCleanupIntervalSeconds)
}

func TestApplyDockerEnvOverrides_SABnzbdCategories(t *testing.T) {
	t.Setenv(
		"ALTMOUNT_SABNZBD_CATEGORIES",
		"movies|1|0|movies|radarr;tv|2|0|tv|sonarr;movies-2160p|3|0|movies-2160p|radarr;tv-2160p|4|0|tv-2160p|sonarr",
	)

	cfg := DefaultConfig()
	err := applyDockerEnvOverrides(cfg)

	assert.NoError(t, err)
	assert.Len(t, cfg.SABnzbd.Categories, 4)
	assert.Equal(t, SABnzbdCategory{
		Name:     "movies-2160p",
		Order:    3,
		Priority: 0,
		Dir:      "movies-2160p",
		Type:     "radarr",
	}, cfg.SABnzbd.Categories[2])
	assert.Equal(t, SABnzbdCategory{
		Name:     "tv-2160p",
		Order:    4,
		Priority: 0,
		Dir:      "tv-2160p",
		Type:     "sonarr",
	}, cfg.SABnzbd.Categories[3])
}

func TestEnvIntRejectsInvalidValue(t *testing.T) {
	t.Setenv("ALTMOUNT_HEALTH_MAX_RETRIES", "not-a-number")

	cfg := DefaultConfig()
	err := applyDockerEnvOverrides(cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ALTMOUNT_HEALTH_MAX_RETRIES")
}

func TestEnvSABnzbdCategoriesRejectsInvalidValue(t *testing.T) {
	t.Setenv("ALTMOUNT_SABNZBD_CATEGORIES", "movies|not-a-number|0|movies|radarr")

	cfg := DefaultConfig()
	err := applyDockerEnvOverrides(cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ALTMOUNT_SABNZBD_CATEGORIES")
}

func TestConfig_NetworkRoundTrip(t *testing.T) {
	in := Config{
		Network: NetworkConfig{
			HTTPProxy:  "http://proxy:3128",
			HTTPSProxy: "http://proxy:3128",
			NoProxy:    "localhost,10.0.0.0/8",
		},
	}
	b, err := yaml.Marshal(in)
	assert.NoError(t, err)

	var out Config
	err = yaml.Unmarshal(b, &out)
	assert.NoError(t, err)

	assert.Equal(t, in.Network, out.Network)
	assert.Equal(t, "http://proxy:3128", out.Network.GetHTTPProxy())
	assert.Equal(t, "http://proxy:3128", out.Network.GetHTTPSProxy())
	assert.Equal(t, "localhost,10.0.0.0/8", out.Network.GetNoProxy())
}

func TestConfig_NetworkDefaultsEmpty(t *testing.T) {
	cfg := Config{}
	assert.Empty(t, cfg.Network.HTTPProxy)
	assert.Empty(t, cfg.Network.HTTPSProxy)
	assert.Empty(t, cfg.Network.NoProxy)
}

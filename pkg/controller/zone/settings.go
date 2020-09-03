package zone

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	crdsv1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	"github.com/replicatedhq/kubeflare/pkg/logger"
	"go.uber.org/zap"
)

func ReconcileSettings(ctx context.Context, instance crdsv1alpha1.Zone, cf *cloudflare.API) error {
	logger.Debug("reconcileSettings for zone", zap.String("zoneName", instance.Name))

	if instance.Spec.Settings == nil {
		logger.Debug("instance does not contain settings to reconcile")
		return nil
	}

	zoneID, err := cf.ZoneIDByName(instance.Name)
	if err != nil {
		return errors.Wrap(err, "failed to get zone id")
	}

	zoneSettingsResponse, err := cf.ZoneSettings(zoneID)
	if err != nil {
		return errors.Wrap(err, "failed to get zone settings")
	}

	updatedZoneSettings := []cloudflare.ZoneSetting{}

	for _, zoneSetting := range zoneSettingsResponse.Result {
		switch zoneSetting.ID {
		case "advanced_ddos":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.AdvancedDDOS)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "always_use_https":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.AlwaysUseHTTPS)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "always_online":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.AlwaysOnline)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "opportunistic_onion":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.OpportunisticOnion)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "automatic_https_rewrites":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.AutomaticHTTPSRewrites)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "browser_cache_ttl":
			needsUpdate := compareAndUpdateIntZoneSetting(&zoneSetting, instance.Spec.Settings.BrowserCacheTTL)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "browser_check":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.BrowserCheck)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "cache_level":
			needsUpdate := compareAndUpdateStringZoneSetting(&zoneSetting, instance.Spec.Settings.CacheLevel)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "challenge_ttl":
			needsUpdate := compareAndUpdateIntZoneSetting(&zoneSetting, instance.Spec.Settings.ChallengeTTL)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "development_mode":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.DevelopmentMode)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "email_obfuscation":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.EmailObfuscation)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "hotlink_protection":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.HotlinkProtection)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "ip_geolocation":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.IPGeolocation)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "ipv6":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.IPV6)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "minify":
			needsUpdate := compareAndUpdateMinifyZoneSetting(&zoneSetting, instance.Spec.Settings.Minify)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "mobile_redirect":
			needsUpdate := compareAndUpdateMobileRedirectZoneSetting(&zoneSetting, instance.Spec.Settings.MobileRedirect)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "mirage":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.Mirage)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "origin_error_page_pass_thru":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.OriginErrorPagePassThru)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "opportunistic_encryption":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.OpportunisticEncryption)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "polish":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.Polish)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "webp":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.WebP)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "brotli":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.Brotli)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "prefetch_preload":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.PrefetchPreload)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "privacy_pass":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.PrivacyPass)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "response_buffering":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.ResponseBuffering)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		case "rocket_loader":
			needsUpdate := compareAndUpdateBoolZoneSetting(&zoneSetting, instance.Spec.Settings.RocketLoader)
			if needsUpdate {
				updatedZoneSettings = append(updatedZoneSettings, zoneSetting)
			}
		}

	}

	if len(updatedZoneSettings) == 0 {
		logger.Debug("no setting was changed in zone", zap.String("zoneName", instance.Name))
		return nil
	}

	logger.Debug("updating zone settings",
		zap.String("zoneName", instance.Name),
		zap.Any("updatedSettings", updatedZoneSettings))

	updateResponse, err := cf.UpdateZoneSettings(zoneID, updatedZoneSettings)
	if err != nil {
		return errors.Wrap(err, "failed to update zone settings")
	}

	if !updateResponse.Success {
		return errors.New("unsuccessful response from cloudflare api")
	}

	return nil
}

func compareAndUpdateBoolZoneSetting(zoneSetting *cloudflare.ZoneSetting, desiredValue *bool) bool {
	if desiredValue == nil {
		return false
	}

	currentValue := zoneSetting.Value == "on"
	if currentValue != *desiredValue {
		if *desiredValue {
			zoneSetting.Value = "on"
		} else {
			zoneSetting.Value = "off"
		}
		return true
	}

	return false
}

func compareAndUpdateIntZoneSetting(zoneSetting *cloudflare.ZoneSetting, desiredValue *int) bool {
	if desiredValue == nil {
		return false
	}

	currentValue := int(zoneSetting.Value.(float64))
	if currentValue != *desiredValue {
		zoneSetting.Value = *desiredValue
		return true
	}

	return false
}

func compareAndUpdateStringZoneSetting(zoneSetting *cloudflare.ZoneSetting, desiredValue *string) bool {
	if desiredValue == nil {
		return false
	}

	currentValue := zoneSetting.Value.(string)
	if currentValue != *desiredValue {
		zoneSetting.Value = *desiredValue
		return true
	}

	return false
}

func compareAndUpdateMinifyZoneSetting(zoneSetting *cloudflare.ZoneSetting, desiredValue *crdsv1alpha1.MinifySetting) bool {
	if desiredValue == nil {
		return false
	}

	isChanged := true
	currentValue := zoneSetting.Value.(map[string]interface{})

	if desiredValue.CSS != nil {
		currentCSS := currentValue["css"] == "on"
		if currentCSS != *desiredValue.CSS {
			if *desiredValue.CSS {
				zoneSetting.Value.(map[string]interface{})["css"] = "on"
			} else {
				zoneSetting.Value.(map[string]interface{})["css"] = "off"
			}
			isChanged = true
		}
	}

	if desiredValue.HTML != nil {
		currentCSS := currentValue["html"] == "on"
		if currentCSS != *desiredValue.HTML {
			if *desiredValue.HTML {
				zoneSetting.Value.(map[string]interface{})["html"] = "on"
			} else {
				zoneSetting.Value.(map[string]interface{})["html"] = "off"
			}
			isChanged = true
		}
	}

	if desiredValue.JS != nil {
		currentCSS := currentValue["js"] == "on"
		if currentCSS != *desiredValue.JS {
			if *desiredValue.JS {
				zoneSetting.Value.(map[string]interface{})["js"] = "on"
			} else {
				zoneSetting.Value.(map[string]interface{})["js"] = "off"
			}
			isChanged = true
		}
	}

	return isChanged
}

func compareAndUpdateMobileRedirectZoneSetting(zoneSetting *cloudflare.ZoneSetting, desiredValue *crdsv1alpha1.MobileRedirect) bool {
	if desiredValue == nil {
		return false
	}

	hasChanged := false

	if desiredValue.Status != nil {
		currentStatus := zoneSetting.Value.(map[string]interface{})["status"].(string) == "on"
		if *desiredValue.Status != currentStatus {
			if *desiredValue.Status {
				zoneSetting.Value.(map[string]interface{})["status"] = "on"
			} else {
				zoneSetting.Value.(map[string]interface{})["status"] = "off"
			}
			hasChanged = true
		}
	}

	if desiredValue.MobileSubdomain != nil {
		if zoneSetting.Value.(map[string]interface{})["mobile_subdomain"] == nil {
			zoneSetting.Value.(map[string]interface{})["mobile_subdomain"] = *desiredValue.MobileSubdomain
			hasChanged = true
		} else {
			currentMobileSubdomain := zoneSetting.Value.(map[string]interface{})["mobile_subdomain"].(string)
			if *desiredValue.MobileSubdomain != currentMobileSubdomain {
				zoneSetting.Value.(map[string]interface{})["mobile_subdomain"] = *desiredValue.MobileSubdomain
				hasChanged = true
			}
		}
	}

	if desiredValue.StripURI != nil {
		currentStripURI := zoneSetting.Value.(map[string]interface{})["strip_uri"].(bool)
		if *desiredValue.StripURI != currentStripURI {
			zoneSetting.Value.(map[string]interface{})["strip_uri"] = *desiredValue.StripURI
			hasChanged = true
		}
	}

	return hasChanged
}

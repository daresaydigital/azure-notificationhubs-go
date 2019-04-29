package notificationhubs

// GetContentType returns Content-Type
// associated with NotificationFormat
func (f NotificationFormat) GetContentType() string {
	switch f {
	case Template,
		AppleFormat,
		GcmFormat,
		KindleFormat,
		BaiduFormat:
		return "application/json"
	}

	return "application/xml"
}

// IsValid identifies whether notification format is valid
func (f NotificationFormat) IsValid() bool {
	return f == Template ||
		f == GcmFormat ||
		f == AppleFormat ||
		f == BaiduFormat ||
		f == KindleFormat ||
		f == WindowsFormat ||
		f == WindowsPhoneFormat
}

// IsValid identifies whether target is valid
func (f TargetPlatform) IsValid() bool {
	return f == AdmPlatform ||
		f == AdmtemplatePlatform ||
		f == ApplePlatform ||
		f == AppletemplatePlatform ||
		f == BaiduPlatform ||
		f == BaidutemplatePlatform ||
		f == GcmPlatform ||
		f == GcmtemplatePlatform ||
		f == TemplatePlatform ||
		f == WindowsphonePlatform ||
		f == WindowsphonetemplatePlatform ||
		f == WindowsPlatform ||
		f == WindowstemplatePlatform
}

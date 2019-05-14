package notificationhubs

// Public constants
const (
	Template           NotificationFormat = "template"
	AppleFormat        NotificationFormat = "apple"
	BaiduFormat        NotificationFormat = "baidu"
	GcmFormat          NotificationFormat = "gcm"
	KindleFormat       NotificationFormat = "adm"
	WindowsFormat      NotificationFormat = "windows"
	WindowsPhoneFormat NotificationFormat = "windowsphone"

	AdmPlatform                  TargetPlatform = "adm"
	AdmTemplatePlatform          TargetPlatform = "admtemplate"
	ApplePlatform                TargetPlatform = "apple"
	AppleTemplatePlatform        TargetPlatform = "appletemplate"
	BaiduPlatform                TargetPlatform = "baidu"
	BaiduTemplatePlatform        TargetPlatform = "baidutemplate"
	GcmPlatform                  TargetPlatform = "gcm"
	GcmTemplatePlatform          TargetPlatform = "gcmtemplate"
	TemplatePlatform             TargetPlatform = "template"
	WindowsphonePlatform         TargetPlatform = "windowsphone"
	WindowsphoneTemplatePlatform TargetPlatform = "windowsphonetemplate"
	WindowsPlatform              TargetPlatform = "windows"
	WindowsTemplatePlatform      TargetPlatform = "windowstemplate"

	// Abandoned: Message processing has been abandoned.
	// It will happen when the message could not be processed within the acceptable time window.
	// By default, it's 30 minutes.
	Abandoned NotificationState = "Abandoned"
	// Canceled: This scheduled message was canceled by user.
	Canceled NotificationState = "Canceled"
	// Completed: Message processing has completed.
	Completed NotificationState = "Completed"
	// Enqueued: Message has been accepted but processing has not yet begun.
	Enqueued NotificationState = "Enqueued"
	// NoTargetFound: No devices were found to send this message.
	NoTargetFound NotificationState = "NoTargetFound"
	// Processing: Message processing has started.
	Processing NotificationState = "Processing"
	// Scheduled: Message is in queue and will be sent on the scheduled time.
	Scheduled NotificationState = "Scheduled"
	// Unknown: Message processing is in an unknown state.
	Unknown NotificationState = "Unknown"

	// AbandonedNotificationMessages: Count of send requests to push service that were dropped because of a timeout.
	AbandonedNotificationMessages NotificationOutcomeName = "AbandonedNotificationMessages"
	// BadChannel: Communication to the push service failed because the channel was invalid.
	BadChannel NotificationOutcomeName = "BadChannel"
	// ChannelDisconnected: Push service disconnected.
	ChannelDisconnected NotificationOutcomeName = "ChannelDisconnected"
	// ChannelThrottled: Push service denied access due to throttling.
	ChannelThrottled NotificationOutcomeName = "ChannelThrottled"
	// Dropped: Push service indicates that the message was dropped.
	Dropped NotificationOutcomeName = "Dropped"
	// ExpiredChannel: Communication to the push service failed because the channel expired.
	ExpiredChannel NotificationOutcomeName = "ExpiredChannel"
	// InvalidCredentials: Credentials used to authenticate to the push service failed.
	InvalidCredentials NotificationOutcomeName = "InvalidCredentials"
	// InvalidNotificationSize: Push request is too large.
	InvalidNotificationSize NotificationOutcomeName = "InvalidNotificationSize"
	// NoTargets: Count of requests that found nothing to send to.
	NoTargets NotificationOutcomeName = "NoTargets"
	// PnsInterfaceError: Push service contract communication failed.
	PnsInterfaceError NotificationOutcomeName = "PnsInterfaceError"
	// PnsServerError: Push service indicated that an error happened on their side.
	PnsServerError NotificationOutcomeName = "PnsServerError"
	// PnsUnavailable: Push service is not available.
	PnsUnavailable NotificationOutcomeName = "PnsUnavailable"
	// PnsUnreachable: Push service was unresponsive.
	PnsUnreachable NotificationOutcomeName = "PnsUnreachable"
	// Skipped: Count of duplicate registrations (same PNS handle found, different registration ID).
	Skipped NotificationOutcomeName = "Skipped"
	// Success: Successfully sent the request to some number of devices.
	Success NotificationOutcomeName = "Success"
	// Throttled: Push service denied access due to throttling.
	Throttled NotificationOutcomeName = "Throttled"
	// UnknownError: An unknown error happened.
	UnknownError NotificationOutcomeName = "UnknownError"
	// WrongToken: The PNS handle was not recognized by the PNS as a valid handle.
	WrongToken NotificationOutcomeName = "WrongToken"
)

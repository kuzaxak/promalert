package main

var alertsStorage = make(map[string]Alert)

func FindAlert(alert Alert) (Alert, bool) {
	if originAlert, ok := alertsStorage[alert.Hash()]; ok {
		return originAlert, true
	} else {
		return Alert{}, false
	}
}

func AddAlert(alert Alert) {
	alertsStorage[alert.Hash()] = alert
}

func DeleteAlert(alert Alert) {
	if _, founded := FindAlert(alert); founded {
		delete(alertsStorage, alert.Hash())
	}
}

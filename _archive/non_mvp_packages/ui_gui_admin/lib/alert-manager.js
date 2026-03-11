// AlertManager handles real-time alerts
class AlertManager {
    constructor() {
        this.alerts = [];
        this.subscribers = new Set();
        this.maxAlerts = 100;
    }

    createAlert(type, severity, message, details = {}) {
        const alert = {
            id: this.generateID(),
            type: type,
            severity: severity,
            message: message,
            details: details,
            timestamp: new Date().toISOString(),
            read: false
        };
        this.alerts.unshift(alert);
        if (this.alerts.length > this.maxAlerts) {
            this.alerts = this.alerts.slice(0, this.maxAlerts);
        }
        this.notifySubscribers(alert);
        return alert;
    }

    getAlerts() {
        return this.alerts;
    }

    markAsRead(id) {
        const alert = this.alerts.find(a => a.id === id);
        if (alert) {
            alert.read = true;
            return true;
        }
        return false;
    }

    markAllAsRead() {
        this.alerts.forEach(a => a.read = true);
    }

    clearAlerts() {
        this.alerts = [];
    }

    getAlertsBySeverity(severity) {
        return this.alerts.filter(a => a.severity === severity);
    }

    subscribe(callback) {
        this.subscribers.add(callback);
        return () => this.subscribers.delete(callback);
    }

    notifySubscribers(alert) {
        this.subscribers.forEach(callback => {
            try {
                callback(alert);
            } catch (e) {
                console.error('Alert subscriber error:', e);
            }
        });
    }

    generateID() {
        return 'alert_' + Date.now() + '_' + Math.random().toString(36).substr(2, 5);
    }

    getUnreadCount() {
        return this.alerts.filter(a => !a.read).length;
    }

    clearOldAlerts(days = 7) {
        const cutoff = new Date(Date.now() - days * 24 * 60 * 60 * 1000);
        this.alerts = this.alerts.filter(a => new Date(a.timestamp) > cutoff);
    }
}

module.exports = AlertManager;
// TrafficFilter handles request/response filtering
class TrafficFilter {
    constructor() {
        this.rules = [];
        this.enabled = true;
    }

    addRule(rule) {
        this.rules.push({
            id: rule.id || this.generateID(),
            name: rule.name,
            pattern: rule.pattern,
            action: rule.action || 'block',
            type: rule.type || 'request',
            enabled: rule.enabled !== false,
            description: rule.description || ''
        });
    }

    removeRule(id) {
        this.rules = this.rules.filter(r => r.id !== id);
    }

    toggleRule(id) {
        const rule = this.rules.find(r => r.id === id);
        if (rule) {
            rule.enabled = !rule.enabled;
        }
    }

    filterRequest(request) {
        if (!this.enabled) {
            return { allowed: true, reason: '' };
        }
        for (const rule of this.rules) {
            if (rule.enabled && rule.type !== 'response') {
                if (this.matchPattern(request, rule.pattern)) {
                    if (rule.action === 'block') {
                        return { allowed: false, reason: rule.name };
                    }
                }
            }
        }
        return { allowed: true, reason: '' };
    }

    filterResponse(response) {
        if (!this.enabled) {
            return { allowed: true, reason: '' };
        }
        for (const rule of this.rules) {
            if (rule.enabled && rule.type !== 'request') {
                if (this.matchPattern(response, rule.pattern)) {
                    if (rule.action === 'block') {
                        return { allowed: false, reason: rule.name };
                    }
                }
            }
        }
        return { allowed: true, reason: '' };
    }

    matchPattern(data, pattern) {
        if (!pattern) return false;
        if (pattern.startsWith('/') && pattern.endsWith('/')) {
            try {
                const regex = new RegExp(pattern.slice(1, -1));
                return regex.test(data);
            } catch (e) {
                return false;
            }
        } else {
            return data.toLowerCase().includes(pattern.toLowerCase());
        }
    }

    generateID() {
        return 'rule_' + Math.random().toString(36).substr(2, 9);
    }

    exportRules() {
        return JSON.stringify(this.rules, null, 2);
    }

    importRules(json) {
        try {
            const rules = JSON.parse(json);
            this.rules = rules;
            return true;
        } catch (e) {
            return false;
        }
    }

    getRules() {
        return this.rules;
    }
}

module.exports = TrafficFilter;
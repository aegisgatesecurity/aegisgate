// Authentication module for AegisGate GUI
const crypto = require('crypto');

class User {
    constructor(id, username, role) {
        this.id = id;
        this.username = username;
        this.role = role;
        this.createdAt = new Date();
    }
}

class Session {
    constructor(token, userId, expiresAt) {
        this.token = token;
        this.userId = userId;
        this.expiresAt = expiresAt;
    }
}

class AuthManager {
    constructor() {
        this.users = new Map();
        this.sessions = new Map();
        this.tokenLength = 32;
    }

    generateToken(length) {
        const bytes = crypto.randomBytes(length);
        return bytes.toString('base64url');
    }

    hashPassword(password) {
        return crypto.createHash('sha256').update(password).digest('base64url');
    }

    register(username, password, role = 'user') {
        for (const user of this.users.values()) {
            if (user.username === username) {
                throw new Error('Username already exists');
            }
        }
        const userId = this.generateToken(16).substr(0, 16);
        const hashedPassword = this.hashPassword(password);
        const user = new User(userId, username, role);
        user.passwordHash = hashedPassword;
        this.users.set(username, user);
        return user;
    }

    login(username, password) {
        const user = this.users.get(username);
        if (!user) throw new Error('Invalid credentials');
        if (user.passwordHash !== this.hashPassword(password)) throw new Error('Invalid credentials');
        const token = this.generateToken(this.tokenLength);
        const expiresAt = new Date(Date.now() + 24 * 60 * 60 * 1000);
        const session = new Session(token, user.id, expiresAt);
        this.sessions.set(token, session);
        return session;
    }

    validateSession(token) {
        const session = this.sessions.get(token);
        if (!session) return null;
        if (new Date() > session.expiresAt) {
            this.sessions.delete(token);
            return null;
        }
        return this.users.get(session.userId) || null;
    }

    logout(token) { this.sessions.delete(token); }

    getUser(userId) {
        for (const user of this.users.values()) {
            if (user.id === userId) return user;
        }
        return null;
    }

    getSessionCount() { return this.sessions.size; }

    usernameExists(username) { return this.users.has(username); }
}

module.exports = AuthManager;
